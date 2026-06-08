# Enterprise ORM — AI Agent Code Generation & Usage Guide

This guide is designed for AI coding assistants and developers. It outlines the directory patterns, model definition conventions, code generation workflows, and query APIs of the Enterprise ORM.

---

## 1. Directory Pattern

Projects utilizing the Enterprise ORM must adhere to the following standard layout:

```text
your-project/
├── db_models/        # 1. Human-defined database table definitions (schema definitions)
│   ├── account.go
│   ├── group.go
│   └── constants.go  # Helper constants for table & column names
├── generate/         # 2. Code-generation entrypoint
│   └── generate.go   # Program containing main() that executes model generation
├── models/           # 3. Auto-generated ORM packages (DO NOT edit manually)
│   ├── client.go     # Database initialization, transactions, and SQL loggers
│   ├── account.go    # Model struct, getters, setters, hooks, and list/bulk operations
│   └── account_predicates.go  # Code-generated WHERE filters for the model
├── migrate/          # 4. Schema migrations runner
│   └── migrate.go    # Program to compute diffs and run migrations
└── main.go           # Application entrypoint
```

---

## 2. Defining Tables (Schemas)

Table schemas are defined as Go functions under the `db_models` package. Each function returns a `*models.Table`.

### Schema Rules & Guidelines
1. **Quoting & Keyword Safety**: Enterprise automatically double-quotes all table and column names in SQL. Always specify identifiers exactly as they should appear in DB.
2. **Primary Keys**: Every table must have a primary key defined via `SetIDField()`.
3. **Builder Chaining**: Field builders return instances supporting method chaining (`SetNillable()`, `AddSerial()`, `Default()`, `DefaultFunc()`).

### Supported Fields
- **Boolean**: `models.BoolField(name)`
- **Byte**: `models.ByteField(name)`
- **UUID**: `models.UUIDField(name)`
- **Integers**: `models.IntField(name)`, `models.SmallIntField(name)`, `models.BigIntField(name)`
- **Unsigned Integer**: `models.UintField(name)`
- **Floats**: `models.Float32Field(name)`, `models.Float64Field(name)`
- **String**: `models.StringField(name)`
- **Enum**: `models.EnumField(name, stringSlice)`
- **Time**: `models.TimeField(name)`
- **JSON**: `models.JSONField(name)`
- **Custom Type**: `models.CustomField(name, postgresType, valueScannerInstance)` (requires implementing `sql.Scanner` and `driver.Valuer`)

### Declaring Relationships
- **Many-to-One**: `models.ManyToOne(targetTable, targetPK, sourceFK)`
- **One-to-Many**: `models.OneToMany(targetTable, sourcePK, targetFK)`
- **Many-to-Many**: `models.ManyToMany(targetTable, sourceFK, targetFK, targetPK, joinTableName)`

### Example: `db_models/account.go`
```go
package db_models

import (
	"github.com/MrSametBurgazoglu/enterprise/models"
	"github.com/google/uuid"
)

func Account() *models.Table {
	idField := models.UUIDField("ID").DefaultFunc(uuid.New)

	tb := &models.Table{
		Fields: []models.FieldI{
			idField,
			models.StringField("Name"),
			models.StringField("Surname"),
			models.UUIDField("DenemeID").SetNillable(),
			models.UintField("Serial").AddSerial(),
		},
		Relations: []*models.Relation{
			models.ManyToOne("Deneme", "id", "deneme_id"),
			models.ManyToMany("Group", "account_id", "group_id", "id", "account_group"),
		},
	}

	tb.SetTableName("account")
	tb.SetIDField(idField)
	
	// Add composite index
	tb.AddIndex("account_name_surname_idx", "Name", "Surname")

	return tb
}
```

---

## 3. Code Generation Workflow

The generator script executes model generation. Create `generate/generate.go`:

```go
package main

import (
	"github.com/MrSametBurgazoglu/enterprise/generate"
	"your-project/db_models"
)

func main() {
	generate.Models(
		db_models.Account(),
		db_models.Group(),
	)
}
```

### Execution
Run the script from your terminal:
```bash
go run generate/generate.go
```
This automatically compiles the templates and writes out the complete `models/` Go files.

---

## 4. Database Setup & Migrations

Initialize the database using `models.Options` and run auto-applied or file-based migrations.

```go
package main

import (
	"context"
	"log"

	"github.com/MrSametBurgazoglu/enterprise/migrate"
	"your-project/db_models"
	"your-project/models"
)

func main() {
	ctx := context.Background()
	postgresURL := "postgresql://user:pass@localhost:5432/db?sslmode=disable"

	// 1. Run migrations
	tables := []*models.Table{
		db_models.Account(),
		db_models.Group(),
	}
	migrate.AutoApplyMigration(ctx, postgresURL, "init_schema", tables...)

	// 2. Setup Client
	opts := &models.Options{
		Url:   postgresURL,
		Debug: true, // Output SQL query logs via slog to stdout
	}
	db, err := models.NewDB(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Exit()
}
```

---

## 5. Eager Relation Loading & Eager Filtering

To load child relations and filter them dynamically in a **single query** (preventing N+1 query patterns), use the `With[RelationName]` methods.

```go
account := models.NewAccount(ctx, db)
account.Where(account.IsIDEqual(accountID))

// Join and filter relations inline
account.WithDeneme(func(d *models.Deneme) {
    d.Where(d.IsActiveEqual(true))
})
account.WithGroupList(func(gl *models.GroupList) {
    gl.Where(gl.IsNameEqual("Admin Group"))
})

err = account.Get()
if err == nil {
    // Nested child entities are automatically populated
    if account.Deneme != nil {
        log.Println("Active Deneme:", account.Deneme.GetID())
    }
}
```

---

## 6. Query Predicates & logical Composition

Where clauses accept list of `client.PredicateI` interfaces.

- **Logical Composition**: Nest queries using `models.And(...)` and `models.Or(...)`.
- **Dynamic Filters**: Use `WhereIf(cond, predicate)` and `WhereIn(cond, predicate)` to append predicates only if `cond == true`.
- **Zero Predicates**: Calling `Where()` with no arguments behaves as a no-op (no SQL WHERE clause).
- **Wildcards**: Use case-sensitive `IsXLike` and case-insensitive `IsXILike` pattern matching.

```go
list := models.NewAccountList(ctx, db)

// Complex filters
list.Where(
    models.And(
        models.Or(
            list.IsNameEqual("Samet"),
            list.IsNameILike("admin%"),
        ),
        list.IsSerialGreater(10),
    ),
)

// Conditional filters
list.WhereIf(filterByStatus, list.IsStatusEqual("active"))

err = list.List()
```

---

## 7. Bulk and In-Database Operations

Avoid loading entities into memory for bulk inserts, updates, or deletes.

```go
list := models.NewAccountList(ctx, db)

// 1. Bulk Create
acc1 := models.NewAccount(ctx, db)
acc1.SetName("Alice")
acc2 := models.NewAccount(ctx, db)
acc2.SetName("Bob")
err := list.Create(acc1, acc2) // Single INSERT statement

// 2. In-Database Update (UpdateWhere)
list.Where(list.IsStatusEqual("pending"))
rowsUpdated, err := list.UpdateWhere(map[string]any{
	"status": "active",
})

// 3. In-Database Delete (DeleteWhere)
list.Where(list.IsStatusEqual("banned"))
rowsDeleted, err := list.DeleteWhere()
```

---

## 8. Aggregates and Iterators (Go 1.23)

Enterprise supports basic aggregates (`Min`, `Max`, `Sum`, `Avg`) on model fields as well as standard library iterators.

### Model Convenience Methods
```go
denemeList := models.NewDenemeList(ctx, db)
maxVal, err := denemeList.MaxCount()
```

### Group By and Go 1.23 iter.Seq2 Iterators
```go
groupList := models.NewGroupList(ctx, db)

var (
	groupName string
	groupSize int
)

// Compile query and get iterator
iterator := groupList.AggregateRowsSeq(func(a *client.Aggregate) {
	a.Field("name", &groupName)
	a.Count("id", &groupSize)
	a.GroupBy("name")
})

// Iterate over results safely
for index, err := range iterator {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Group %d: Name: %s, Size: %d\n", index, groupName, groupSize)
}
```
