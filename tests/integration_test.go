package tests

import (
	"context"
	"os"
	"testing"

	"github.com/MrSametBurgazoglu/enterprise/client"
	"github.com/MrSametBurgazoglu/enterprise/migrate"
	"github.com/MrSametBurgazoglu/enterprise/tests/db_models"
	"github.com/MrSametBurgazoglu/enterprise/tests/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		t.Skip("Skipping integration test; POSTGRES_URL environment variable is not set")
	}

	ctx := context.Background()

	// 1. Run migrations using Atlas schema auto-apply
	t.Log("Applying migrations...")
	migrate.AutoApplyMigration(
		ctx,
		postgresURL,
		"integration_test_init",
		db_models.Deneme(),
		db_models.Test(),
		db_models.Account(),
		db_models.Group(),
	)

	// 2. Initialize Database Client
	opts := &models.Options{
		Url:   postgresURL,
		Debug: true,
	}
	db, err := models.NewDB(opts)
	assert.NoError(t, err)
	defer db.Exit()

	// Clean tables before starting test to avoid residual state from previous failed runs
	_, _ = db.Exec(ctx, "TRUNCATE TABLE \"deneme\", \"group\", \"account\", \"test\" CASCADE;")

	// 3. Test Create with defaults (including UUID default uuid.New and Bool default true)
	t.Log("Testing Create with defaults...")
	deneme := models.NewDeneme(ctx, db)
	deneme.SetCount(42)
	deneme.SetDenemeType(models.DenemeTypeDeneme)
	err = deneme.Create()
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, deneme.GetID())
	assert.True(t, deneme.GetIsActive())

	// 4. Test Get
	t.Log("Testing Get...")
	fetched := models.NewDeneme(ctx, db)
	fetched.Where(fetched.IsIDEqual(deneme.GetID()))
	err = fetched.Get()
	assert.NoError(t, err)
	assert.Equal(t, 42, fetched.GetCount())
	assert.Equal(t, models.DenemeTypeDeneme, fetched.GetDenemeType())

	// 4a. Test IsTestIDNil and IsTestIDNotNil predicates
	t.Log("Testing IsTestIDNil/IsTestIDNotNil predicates...")
	fetchedNil := models.NewDeneme(ctx, db)
	fetchedNil.Where(fetchedNil.IsIDEqual(deneme.GetID()), fetchedNil.IsTestIDNil())
	err = fetchedNil.Get()
	assert.NoError(t, err)
	assert.Equal(t, 42, fetchedNil.GetCount())
	assert.Nil(t, fetchedNil.GetTestID())
	assert.Equal(t, uuid.Nil, fetchedNil.GetTestIDValue())

	fetchedNotNil := models.NewDeneme(ctx, db)
	fetchedNotNil.Where(fetchedNotNil.IsIDEqual(deneme.GetID()), fetchedNotNil.IsTestIDNotNil())
	err = fetchedNotNil.Get()
	assert.Error(t, err) // Should not find because TestID is nil

	// 5. Test Update and Get again
	t.Log("Testing Update...")
	fetched.SetCount(100)
	err = fetched.Update()
	assert.NoError(t, err)

	fetched2 := models.NewDeneme(ctx, db)
	fetched2.Where(fetched2.IsIDEqual(deneme.GetID()))
	err = fetched2.Get()
	assert.NoError(t, err)
	assert.Equal(t, 100, fetched2.GetCount())

	// 6. Test OR Predicate composition & parameter suffix-indexing
	t.Log("Testing OR Predicate and parameter suffix-indexing...")
	deneme2 := models.NewDeneme(ctx, db)
	deneme2.SetCount(200)
	deneme2.SetDenemeType(models.DenemeTypeTest)
	err = deneme2.Create()
	assert.NoError(t, err)

	list := models.NewDenemeList(ctx, db)
	list.Where(models.Or(
		list.IsCountEqual(100),
		list.IsCountEqual(200),
	))
	err = list.List()
	assert.NoError(t, err)
	assert.True(t, len(list.Items) >= 2)

	// 7. Test reserved keyword table / quoting (e.g. Group model matches SQL reserved keyword "group")
	t.Log("Testing reserved keyword table and quoting...")
	g := models.NewGroup(ctx, db)
	g.SetName("Test Group")
	g.SetSurname("Group Surname")
	g.SetData(map[string]any{"key": "value"})
	err = g.Create()
	assert.NoError(t, err)

	fetchedG := models.NewGroup(ctx, db)
	fetchedG.Where(fetchedG.IsIDEqual(g.GetID()))
	err = fetchedG.Get()
	assert.NoError(t, err)
	assert.Equal(t, "Test Group", fetchedG.GetName())
	assert.Equal(t, "value", fetchedG.GetData()["key"])

	// 8. Test the new ORM Gap Features
	t.Log("Testing new ORM Gap Features (Distinct, Partial Select, LIKE/ILIKE, Empty WHERE, Count)...")
	g1 := models.NewGroup(ctx, db)
	g1.SetName("Group A")
	g1.SetSurname("SurnameA")
	g1.SetData(map[string]any{})
	err = g1.Create()
	assert.NoError(t, err)

	g2 := models.NewGroup(ctx, db)
	g2.SetName("Group B")
	g2.SetSurname("surnameA")
	g2.SetData(map[string]any{})
	err = g2.Create()
	assert.NoError(t, err)

	g3 := models.NewGroup(ctx, db)
	g3.SetName("Group C")
	g3.SetSurname("OtherB")
	g3.SetData(map[string]any{})
	err = g3.Create()
	assert.NoError(t, err)

	// A. Empty WHERE (Unconditional List)
	allGroups := models.NewGroupList(ctx, db)
	err = allGroups.List()
	assert.NoError(t, err)
	assert.True(t, len(allGroups.Items) >= 4)

	// B. Count helper
	count, err := allGroups.Count()
	assert.NoError(t, err)
	assert.Equal(t, len(allGroups.Items), count)

	// C. SELECT DISTINCT (and ORDER BY quoting/sorting)
	distinctList := models.NewGroupList(ctx, db)
	distinctList.Order("surname")
	surnames, err := distinctList.DistinctString(models.GroupTableSurnameField)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Group Surname", "OtherB", "SurnameA", "surnameA"}, surnames)

	// D. LIKE / ILIKE predicates
	likeList := models.NewGroupList(ctx, db)
	likeList.Where(likeList.IsSurnameLike("Surname%"))
	err = likeList.List()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(likeList.Items))

	ilikeList := models.NewGroupList(ctx, db)
	ilikeList.Where(ilikeList.IsSurnameILike("surname%"))
	err = ilikeList.List()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ilikeList.Items))

	// E. Partial Field Selection (GetSelector on List)
	partialList := models.NewGroupList(ctx, db)
	partialList.GetSelector().SelectSurname()
	partialList.Where(partialList.IsSurnameILike("surname%"))
	err = partialList.List()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(partialList.Items))
	for _, item := range partialList.Items {
		assert.NotEmpty(t, item.GetSurname())
		assert.Empty(t, item.GetName()) // Not selected, should be empty!
	}

	// F. Test Aggregates on DenemeList
	t.Log("Testing Aggregates on DenemeList...")
	denemeList := models.NewDenemeList(ctx, db)
	denemeList.Where(models.Or(
		denemeList.IsCountEqual(100),
		denemeList.IsCountEqual(200),
	))

	minCount, err := denemeList.MinCount()
	assert.NoError(t, err)
	assert.Equal(t, 100, minCount)

	maxCount, err := denemeList.MaxCount()
	assert.NoError(t, err)
	assert.Equal(t, 200, maxCount)

	sumCount, err := denemeList.SumCount()
	assert.NoError(t, err)
	assert.Equal(t, 300, sumCount)

	avgCount, err := denemeList.AvgCount()
	assert.NoError(t, err)
	assert.Equal(t, 150.0, avgCount)

	// G. Test Block-Based Transactions (Commit)
	t.Log("Testing Transaction block commit...")
	var txDenemeID uuid.UUID
	err = db.Transaction(ctx, func(tx client.DatabaseTransactionClient) error {
		txDeneme := models.NewDeneme(ctx, tx)
		txDeneme.SetCount(500)
		txDeneme.SetDenemeType(models.DenemeTypeDeneme)
		if err := txDeneme.Create(); err != nil {
			return err
		}
		txDenemeID = txDeneme.GetID()
		return nil
	})
	assert.NoError(t, err)

	// Verify it exists in db
	txFetched := models.NewDeneme(ctx, db)
	txFetched.Where(txFetched.IsIDEqual(txDenemeID))
	err = txFetched.Get()
	assert.NoError(t, err)
	assert.Equal(t, 500, txFetched.GetCount())

	// G2. Test Block-Based Transactions (Rollback on Error)
	t.Log("Testing Transaction block rollback on error...")
	var txDenemeID2 uuid.UUID
	err = db.Transaction(ctx, func(tx client.DatabaseTransactionClient) error {
		txDeneme := models.NewDeneme(ctx, tx)
		txDeneme.SetCount(600)
		txDeneme.SetDenemeType(models.DenemeTypeDeneme)
		if err := txDeneme.Create(); err != nil {
			return err
		}
		txDenemeID2 = txDeneme.GetID()
		return assert.AnError // Return an error to trigger rollback
	})
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)

	// Verify it does NOT exist in db
	txFetched2 := models.NewDeneme(ctx, db)
	txFetched2.Where(txFetched2.IsIDEqual(txDenemeID2))
	err = txFetched2.Get()
	assert.Error(t, err) // Should fail to find the row

	// G3. Test Block-Based Transactions (Rollback on Panic)
	t.Log("Testing Transaction block rollback on panic...")
	var txDenemeID3 uuid.UUID
	assert.Panics(t, func() {
		_ = db.Transaction(ctx, func(tx client.DatabaseTransactionClient) error {
			txDeneme := models.NewDeneme(ctx, tx)
			txDeneme.SetCount(700)
			txDeneme.SetDenemeType(models.DenemeTypeDeneme)
			if err := txDeneme.Create(); err != nil {
				return err
			}
			txDenemeID3 = txDeneme.GetID()
			panic("something went wrong inside transaction")
		})
	})

	// Verify it does NOT exist in db
	txFetched3 := models.NewDeneme(ctx, db)
	txFetched3.Where(txFetched3.IsIDEqual(txDenemeID3))
	err = txFetched3.Get()
	assert.Error(t, err) // Should fail to find the row

	// 8b. Test new ORM features (ListWithTotal, WhereIf/WhereIn, Zero Predicates, Expression aggregates)
	t.Log("Testing new ORM Features...")

	// Test Zero Predicates (no-op WHERE)
	zeroPredList := models.NewGroupList(ctx, db)
	zeroPredList.Where()
	err = zeroPredList.List()
	assert.NoError(t, err)

	// Test WhereIf / WhereIn
	condList := models.NewGroupList(ctx, db)
	condList.WhereIf(false, condList.IsNameEqual("nonexistent"))
	condList.WhereIf(true, condList.IsNameEqual("Group A"))
	err = condList.List()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(condList.Items))
	assert.Equal(t, "Group A", condList.Items[0].GetName())

	// Test ListWithTotal
	totalList := models.NewGroupList(ctx, db)
	totalList.Where(totalList.IsSurnameILike("surname%"))
	total, err := totalList.ListWithTotal(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, 2, total) // SurnameA and surnameA
	assert.Equal(t, 1, len(totalList.Items))

	// Test Expression Aggregate with GroupBy
	exprList := models.NewGroupList(ctx, db)
	var countVal int
	scanFunc, err := exprList.Aggregate(func(a *client.Aggregate) {
		a.Count(client.Raw("COALESCE(name, 'default')"), &countVal)
		a.GroupBy(client.Raw("COALESCE(name, 'default')"))
	})
	assert.NoError(t, err)
	err = scanFunc()
	assert.NoError(t, err)
	assert.True(t, countVal > 0)

	// 9. Test Delete
	t.Log("Testing Delete...")
	err = txFetched.Delete()
	assert.NoError(t, err)

	err = deneme.Delete()
	assert.NoError(t, err)

	err = deneme2.Delete()
	assert.NoError(t, err)

	err = g.Delete()
	assert.NoError(t, err)

	err = g1.Delete()
	assert.NoError(t, err)

	err = g2.Delete()
	assert.NoError(t, err)

	err = g3.Delete()
	assert.NoError(t, err)
}
