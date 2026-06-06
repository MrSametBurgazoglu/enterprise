package tests

import (
	"context"
	"testing"

	"github.com/MrSametBurgazoglu/enterprise/client"
	"github.com/MrSametBurgazoglu/enterprise/tests/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

type mockRows struct{}

func (m *mockRows) Close()                                       {}
func (m *mockRows) Err() error                                   { return nil }
func (m *mockRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (m *mockRows) Next() bool                                   { return false }
func (m *mockRows) Scan(dest ...any) error                       { return nil }
func (m *mockRows) Values() ([]any, error)                       { return nil, nil }
func (m *mockRows) RawValues() [][]byte                          { return nil }
func (m *mockRows) Conn() *pgx.Conn                              { return nil }

type mockRow struct {
	val int
}

func (m *mockRow) Scan(dest ...any) error {
	if len(dest) == 1 {
		if ptr, ok := dest[0].(*int); ok {
			*ptr = m.val
			return nil
		}
	}
	return nil
}

type mockDB struct {
	lastSQL  string
	lastArgs []any
	queries  []string
	countVal int
}

func (m *mockDB) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	m.lastSQL = sql
	m.lastArgs = arguments
	m.queries = append(m.queries, sql)
	return pgconn.NewCommandTag("UPDATE 1"), nil
}

func (m *mockDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	m.lastSQL = sql
	m.lastArgs = args
	m.queries = append(m.queries, sql)
	return &mockRows{}, nil
}

func (m *mockDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	m.lastSQL = sql
	m.lastArgs = args
	m.queries = append(m.queries, sql)
	return &mockRow{val: m.countVal}
}

func TestWhere_ZeroPredicates(t *testing.T) {
	ctx := context.Background()
	db := &mockDB{}
	list := models.NewTestList(ctx, db)

	// Call Where with zero predicates (Bug 14)
	list.Where()
	err := list.List()
	assert.NoError(t, err)

	// Verify that the generated SQL does NOT have "WHERE" at all (or "WHERE (())")
	assert.Contains(t, db.lastSQL, "SELECT")
	assert.NotContains(t, db.lastSQL, "WHERE")
}

func TestWhereIf_WhereIn(t *testing.T) {
	ctx := context.Background()

	// 1. Condition is false
	db := &mockDB{}
	list := models.NewTestList(ctx, db)
	list.WhereIf(false, list.IsNameEqual("test_name"))
	list.WhereIn(false, list.IsNameIN("a", "b"))
	err := list.List()
	assert.NoError(t, err)
	assert.NotContains(t, db.lastSQL, "WHERE")

	// 2. Condition is true
	db = &mockDB{}
	list = models.NewTestList(ctx, db)
	list.WhereIf(true, list.IsNameEqual("test_name"))
	err = list.List()
	assert.NoError(t, err)
	assert.Contains(t, db.lastSQL, "WHERE")
	assert.Contains(t, db.lastSQL, "= @test__name_1")
}

func TestListWithTotal(t *testing.T) {
	ctx := context.Background()
	db := &mockDB{countVal: 15}
	list := models.NewTestList(ctx, db)

	list.Where(list.IsNameEqual("hello"))
	total, err := list.ListWithTotal(10, 5)
	assert.NoError(t, err)
	assert.Equal(t, 15, total)

	// We expect two queries: Count and paged List
	assert.Len(t, db.queries, 2)
	assert.Contains(t, db.queries[0], "SELECT COUNT(*)")
	assert.Contains(t, db.queries[1], "LIMIT 5 OFFSET 10")
}

func TestAggregate_ExpressionAutoQuoting(t *testing.T) {
	ctx := context.Background()
	db := &mockDB{}
	list := models.NewTestList(ctx, db)

	// Call aggregate with an expression (Bug 13)
	var avgVal float64
	_, err := list.Aggregate(func(a *client.Aggregate) {
		a.Avg("DATE(created_at)", &avgVal)
		a.GroupBy("DATE(created_at)")
	})
	assert.NoError(t, err)

	// The SQL should NOT contain double-quoted `"DATE(created_at)"`
	assert.Contains(t, db.lastSQL, "AVG(DATE(created_at))")
	assert.Contains(t, db.lastSQL, "GROUP BY DATE(created_at)")
	assert.NotContains(t, db.lastSQL, "\"DATE(created_at)\"")
}
