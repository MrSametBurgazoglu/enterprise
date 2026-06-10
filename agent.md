# Enterprise ORM — Agent Guide

A type-safe, code-generated PostgreSQL ORM for Go built on top of [pgx](https://github.com/jackc/pgx) and [Atlas](https://atlasgo.io).

---

## Project Layout

```
enterprise/
├── client/          # Core query builder: Where, WhereList, PredicateI, Client
├── generate/        # Go text/template files (.tpl) for code generation
├── logger/          # Structured JSON logger (slog-based)
├── migrate/         # Atlas-based schema migration helpers
├── mock/            # pgxmock wrapper for unit tests
├── models/          # Field type definitions (FieldI interface + concrete field types)
├── pg/              # Low-level pgx connection helpers
├── tests/           # All tests (unit + integration)
│   ├── db_models/   # Schema definitions used by the test suite
│   ├── models/      # Code-generated ORM models (do not edit by hand)
│   ├── generate/    # Code-generation driver (go run generate/generate.go)
│   ├── makefile     # make targets: generate-models, run-tests, integration-test
│   └── *_test.go    # Unit and integration test files
├── docker-compose.yml  # Local PostgreSQL for integration tests
├── go.mod / go.sum
└── agent.md         # This file
```

---

## Key Packages and Their Roles

### `models/` — Field Definitions

Defines `FieldI` (interface) and concrete field types (`UUIDDBField`, `StringDBField`, `IntDBField`, etc.).

Every field type embeds `*Field` (the base struct) which holds:
- `defaultFunc reflect.Value` — an arbitrary function called at `SetDefaults()` time via reflection.
- `HaveDefault bool`, `RequiredPackages []string`.

**Adding a new field type:**
1. Create `models/mytype.go` embedding `*Field`.
2. Implement `GetDefault() string` — delegate to `i.Field.GetDefault()` when `i.defaultFunc.IsValid()`, otherwise return a literal.
3. Implement `DefaultFunc(v func() MyType) *MyTypeDBField` — call `i.Field.DefaultFunc(v)` and return self.
4. Add a `FieldTypeMyType` constant in `models/field_type.go`.
5. Add a `case models.FieldTypeMyType` in `migrate/transformer.go → TransformFieldToAtlasColumn`.

### `client/` — Query Builder

| File | Responsibility |
|---|---|
| `where.go` | `PredicateI` interface, `Where` struct, `LogicalPredicate`, `Or()`, `And()` |
| `where_list.go` | `WhereList` — a slice of `PredicateI`, `Parse(tableName, args)` |
| `relation.go` | `Relation`, `RelationCondition`, join-string building, relation-where parsing |
| `client.go` | `Client` struct — `Get`, `List`, `Create`, `Update`, `Delete`, `Refresh`, `Aggregate`, `BulkCreate`, `BulkUpdate`, `BulkDelete`, and their SQL-builder helpers |
| `res.go` | `Res` — result holder for SQL strings |
| `order.go` | `Order` struct |
| `paging.go` | `Paging` struct |
| `aggregate.go` | `Aggregate` struct |

#### `PredicateI` interface

```go
type PredicateI interface {
    Parse(tableName string, args pgx.NamedArgs) string
    GetName() string
}
```

Both `*Where` and `*LogicalPredicate` implement this.  
`Or(preds ...PredicateI)` and `And(preds ...PredicateI)` return a `*LogicalPredicate`.

#### Parameter naming convention

Parameters are named `<table>__<column>_<index>` (e.g. `account__id_1`).  
The index is computed per-`Parse` call by scanning `args` for existing keys — guaranteeing uniqueness even when the same predicate appears multiple times (fixes the duplicate-parameter overwrite bug).

#### SQL quoting

All table names and column names in generated SQL are double-quoted (`"table"."column"`) to avoid collisions with PostgreSQL reserved keywords (e.g. `group`, `order`, `user`).

### `generate/` — Code Generator Templates

Templates use Go's `text/template` package. Run via `go run generate/generate.go` from inside `tests/`.

| Template | Generates |
|---|---|
| `schema_struct.go.tpl` | Model struct, getters/setters, `SetDefaults()`, `Create/Update/Delete/Get/List` methods, `*Result` type |
| `predicates.go.tpl` | `*Predicate` struct with `Where(PredicateI...)`, `ORWhere(PredicateI...)`, `IsXEqual`, `IsXGreater`, etc. |
| `client.go.tpl` | Package-level `Database`, `IDatabase`, `Transaction`, `NewDB()`, `Or()`, `And()` helpers |

**Important template rules:**
- `{{if .IsNillable}}t.{{.GetNameLower}}{{else}}&t.{{.GetNameLower}}{{end}}` — nillable fields are already pointers; do not take their address again.
- `t.serialFields = nil` is emitted at the top of `SetDefaults()` to prevent duplicate RETURNING columns on repeated calls.
- `Where` and `ORWhere` accept `...client.PredicateI`, not `...*client.Where`.

### `migrate/` — Schema Migrations

`AutoApplyMigration(ctx, postgresURL, planName, tables...)` — diffs the current DB schema against the desired Atlas schema and applies changes in-process (no migration files).

`TransformFieldToAtlasColumn` maps `models.FieldType*` → Atlas `schema.Type`:

| Field type | Atlas type |
|---|---|
| `FieldTypeSmallInt` (serial) | `postgres.SerialType{T: TypeSmallSerial}` |
| `FieldTypeInt` (serial) | `postgres.SerialType{T: TypeSerial}` |
| `FieldTypeBigInt` (serial) | `postgres.SerialType{T: TypeBigSerial}` |
| `FieldTypeUint` (serial) | `postgres.SerialType{T: TypeSerial}` |
| `FieldTypeCustom` (`"text"`) | `schema.StringType{T: "text"}` |
| `FieldTypeCustom` (other) | `schema.UnsupportedType{T: customType}` |
| `FieldTypeDecimal` | `schema.DecimalType{T: postgres.TypeNumeric}` |
| `FieldTypeStringArray` | `postgres.ArrayType{T: "text[]"}` |

**Index columns** are resolved through a `fieldNameMap` (Go name → DB column name) before being passed to Atlas, ensuring column references use the actual snake_case DB names.

---

## Testing

### Unit tests

```bash
cd tests
make run-tests       # go test -v
# or from project root:
go test ./...
```

All unit tests use `mock.NewMockClient()` (pgxmock). SQL expectations include:
- Double-quoted identifiers: `"account"."id"`
- Suffix-indexed parameters: `@account__id_1`

### Integration tests

```bash
cd tests
make integration-test
```

This target:
1. Runs `make generate-models` (regenerates `tests/models/`).
2. Starts a `postgres:16-alpine` Docker container.
3. Runs `TestIntegration` with `POSTGRES_URL` pointing at it.
4. Tears down the container.

`TestIntegration` covers:
- `AutoApplyMigration` (Atlas schema auto-apply).
- **Create** with `SetDefaults()` (UUID auto-generation, bool default `true`).
- **Get** after create.
- **Update** + **Get** to verify the change.
- **OR predicate** composition and suffix-indexed parameter uniqueness.
- **Reserved keyword table** (`"group"`) create + get.
- **Delete** for all created rows.

---

## Development Workflow

### 1. Change a schema model

Edit files in `tests/db_models/`, then regenerate:

```bash
cd tests
make generate-models
```

### 2. Add a new ORM feature

1. Edit the relevant `client/*.go` file(s).
2. If the feature affects generated code, edit the corresponding `.tpl` file in `generate/`.
3. Regenerate models: `make generate-models`.
4. Update mock-based unit tests in `tests/*_test.go` to match new SQL expectations.
5. Add integration test coverage in `tests/integration_test.go`.
6. Verify: `go test ./...` and `make integration-test`.

### 3. Add a new field type to the ORM

1. `models/mynewtype.go` — implement `FieldI`.
2. `models/field_type.go` — add `FieldTypeMyNewType` constant.
3. `migrate/transformer.go` — add a `case` in `TransformFieldToAtlasColumn`.
4. `generate/schema_struct.go.tpl` — add any special-case template logic if needed.
5. `generate/predicates.go.tpl` — add predicate helpers if the field supports filtering.
6. Regenerate and test.

---

## Architecture Decisions

| Decision | Rationale |
|---|---|
| `pgx.NamedArgs` threaded through all query builders | Enables clean parameter deduplication and suffix-indexing without a global counter |
| `PredicateI` interface over concrete `*Where` | Allows arbitrary logical nesting (`Or`, `And`) without changing call sites |
| Double-quoting all identifiers | Prevents reserved-keyword SQL errors for fields/tables named `group`, `order`, `user`, etc. |
| `postgres.SerialType` (not `schema.IntegerType`) | Atlas requires its own `SerialType` struct for `serial`/`bigserial` columns; using `IntegerType` produces a runtime pq error |
| `SetDefaults()` resets `serialFields` | Prevents duplicate `RETURNING` columns when `SetDefaults()` is called more than once |
| `defaultFunc reflect.Value` in base `Field` | Generalizes dynamic defaults to all field types without duplicating `FuncStruct` logic |

---

## Dependencies

| Package | Version | Purpose |
|---|---|---|
| `github.com/jackc/pgx/v5` | v5.x | PostgreSQL driver + `pgx.NamedArgs` |
| `ariga.io/atlas` | v0.25.0 | Schema diffing and migration planning |
| `github.com/google/uuid` | latest | UUID generation for default values |
| `github.com/lib/pq` | latest | Atlas sqlclient dialect driver |
| `github.com/pashagolub/pgxmock/v4` | v4.x | Mock database for unit tests |
| `github.com/stretchr/testify` | latest | Test assertions |
