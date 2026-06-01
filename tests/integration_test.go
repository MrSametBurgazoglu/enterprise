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

	// 8. Test Delete
	t.Log("Testing Delete...")
	err = deneme.Delete()
	assert.NoError(t, err)

	err = deneme2.Delete()
	assert.NoError(t, err)

	err = g.Delete()
	assert.NoError(t, err)
}
