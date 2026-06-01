package tests

import (
	"context"
	"os"
	"testing"

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

	// 9. Test Delete
	t.Log("Testing Delete...")
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
