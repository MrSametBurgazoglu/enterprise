package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/MrSametBurgazoglu/enterprise/client"
	"github.com/MrSametBurgazoglu/enterprise/enterprisetest"
	"github.com/MrSametBurgazoglu/enterprise/migrate"
	"github.com/MrSametBurgazoglu/enterprise/tests/db_models"
	"github.com/MrSametBurgazoglu/enterprise/tests/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

type testObserver struct {
	mu      sync.Mutex
	queries []string
}

func (o *testObserver) Before(ctx context.Context, sql string, args pgx.NamedArgs) context.Context {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.queries = append(o.queries, sql)
	return ctx
}

func (o *testObserver) After(ctx context.Context, sql string, args pgx.NamedArgs, err error) {}

func TestNewDBHarnessAndFeatures(t *testing.T) {
	// 1. Spin up ephemeral PG via our first-party test harness!
	postgresURL := enterprisetest.NewDB(t)

	ctx := context.Background()

	// 2. Apply migrations
	migrate.AutoApplyMigration(
		ctx,
		postgresURL,
		"harness_test_init",
		db_models.Deneme(),
		db_models.Test(),
		db_models.Account(),
		db_models.Group(),
	)

	// 3. Initialize DB Client with Observability Hook
	observer := &testObserver{}
	opts := &models.Options{
		Url:      postgresURL,
		Debug:    true,
		Observer: observer,
	}
	db, err := models.NewDB(opts)
	assert.NoError(t, err)
	defer db.Exit()

	// 4. Test Column Projection Select(fields...)
	t.Log("Testing Column Projection...")
	d := models.NewDeneme(ctx, db)
	d.SetCount(10)
	d.SetDenemeType(models.DenemeTypeDeneme)
	err = d.Create()
	assert.NoError(t, err)

	fetched := models.NewDeneme(ctx, db)
	fetched.Select(models.DenemeTableCountField) // project only Count field
	fetched.Where(fetched.IsIDEqual(d.GetID()))
	err = fetched.Get()
	assert.NoError(t, err)
	assert.Equal(t, 10, fetched.GetCount())
	assert.Equal(t, models.DenemeType(""), fetched.GetDenemeType()) // not selected!

	// 5. Test Lazy WhereIfFn
	t.Log("Testing WhereIfFn...")
	var nilPtr *int
	list := models.NewDenemeList(ctx, db)
	// This would panic with WhereIf since nilPtr derefs, but WhereIfFn guards it safely
	list.WhereIfFn(nilPtr != nil, func() client.PredicateI {
		return list.IsCountEqual(*nilPtr)
	})
	err = list.List()
	assert.NoError(t, err)
	assert.True(t, len(list.Items) >= 1)

	// 6. Test GetByIDs and GetByIDsMap
	t.Log("Testing GetByIDs and GetByIDsMap...")
	d2 := models.NewDeneme(ctx, db)
	d2.SetCount(20)
	d2.SetDenemeType(models.DenemeTypeTest)
	err = d2.Create()
	assert.NoError(t, err)

	list2 := models.NewDenemeList(ctx, db)
	items, err := list2.GetByIDs(d.GetID(), d2.GetID())
	assert.NoError(t, err)
	assert.Len(t, items, 2)

	list3 := models.NewDenemeList(ctx, db)
	itemsMap, err := list3.GetByIDsMap(d.GetID(), d2.GetID())
	assert.NoError(t, err)
	assert.Len(t, itemsMap, 2)
	assert.Equal(t, 10, itemsMap[d.GetID()].GetCount())
	assert.Equal(t, 20, itemsMap[d2.GetID()].GetCount())

	// 7. Test SumExpr and AvgExpr
	t.Log("Testing SumExpr and AvgExpr...")
	list4 := models.NewDenemeList(ctx, db)
	var sumResult float64
	err = list4.SumExpr("count * 2", "numeric", &sumResult)
	assert.NoError(t, err)
	assert.Equal(t, float64(60), sumResult) // (10*2) + (20*2) = 60

	// 8. Test Range-over-func AggregateSeq
	t.Log("Testing AggregateSeq...")
	list5 := models.NewDenemeList(ctx, db)
	var totalSum int
	for _, err := range list5.AggregateSeq(func(a *client.Aggregate) {
		a.Sum("count", &totalSum)
	}) {
		assert.NoError(t, err)
	}
	assert.Equal(t, 30, totalSum)

	// 9. Test First-class Transaction db.Tx
	t.Log("Testing db.Tx transaction block...")
	var txDenemeID uuid.UUID
	err = db.Tx(ctx, func(tx *models.Transaction) error {
		txDeneme := models.NewDeneme(ctx, tx)
		txDeneme.SetCount(99)
		txDeneme.SetDenemeType(models.DenemeTypeDeneme)
		err := txDeneme.Create()
		assert.NoError(t, err)
		txDenemeID = txDeneme.GetID()
		return nil
	})
	assert.NoError(t, err)

	// Verify transaction committed successfully
	txFetched := models.NewDeneme(ctx, db)
	txFetched.Where(txFetched.IsIDEqual(txDenemeID))
	err = txFetched.Get()
	assert.NoError(t, err)
	assert.Equal(t, 99, txFetched.GetCount())

	// 10. Verify structured QueryObserver caught queries
	observer.mu.Lock()
	defer observer.mu.Unlock()
	assert.NotEmpty(t, observer.queries)
	assert.Contains(t, observer.queries[len(observer.queries)-1], "SELECT")
}
