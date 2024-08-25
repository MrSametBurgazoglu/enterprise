package migrate

import (
	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlclient"
	"context"
	"fmt"
	"github.com/MrSametBurgazoglu/enterprise/models"
	_ "github.com/lib/pq"
	"log"
	"net/url"
	"os"
)

func Migrate(ctx context.Context, postgresUrl, migrationPath, planName string, tables []*models.Table) {
	// Define the migration directory
	dir := getMigrationDirectory(migrationPath)

	parsedURL := getURL(postgresUrl)

	drv, closeClient := getDriver(ctx, parsedURL)
	defer closeClient()

	// Inspect the current state of the database
	currentState, err := drv.InspectSchema(ctx, parsedURL.Schema, &schema.InspectOptions{Exclude: []string{"atlas_schema_revisions"}})
	if err != nil {
		log.Fatal(err)
	}

	desired := TransformSchemaToAtlasSchema(parsedURL.Schema, tables)

	// Compute the diff between the current and desired state
	diff, err := drv.SchemaDiff(currentState, desired)
	if err != nil {
		fmt.Printf("failed computing schema diff: %v\n", err)
		os.Exit(1)
	}

	// Plan the migration
	plan, err := drv.PlanChanges(ctx, planName, diff)
	if err != nil {
		fmt.Printf("failed planning changes: %v\n", err)
		os.Exit(1)
	}

	// Create the planner
	pl := migrate.NewPlanner(drv, dir)

	// Write the migration plan
	if err := pl.WritePlan(plan); err != nil {
		fmt.Printf("failed writing migration plan: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("migration plan created successfully")
}

func AutoApplyMigration(ctx context.Context, postgresUrl, planName string, tables ...*models.Table) {
	parsedURL := getURL(postgresUrl)

	drv, closeClient := getDriver(ctx, parsedURL)
	defer closeClient()

	// Inspect the current state of the database
	currentState, err := drv.InspectSchema(ctx, parsedURL.Schema, &schema.InspectOptions{Exclude: []string{"atlas_schema_revisions"}})
	if err != nil {
		log.Fatal(err)
	}

	desired := TransformSchemaToAtlasSchema(parsedURL.Schema, tables)

	// Compute the diff between the current and desired state
	diff, err := drv.SchemaDiff(currentState, desired)
	if err != nil {
		fmt.Printf("failed computing schema diff: %v\n", err)
		os.Exit(1)
	}

	// Plan the migration
	plan, err := drv.PlanChanges(ctx, planName, diff)
	if err != nil {
		fmt.Printf("failed planning changes: %v\n", err)
		os.Exit(1)
	}

	for _, c := range plan.Changes {
		if _, err := drv.ExecContext(ctx, c.Cmd, c.Args...); err != nil {
			panic(err)
		}
	}
}

func getURL(postgresURL string) *sqlclient.URL {
	u, err := url.Parse(postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	return &sqlclient.URL{URL: u, DSN: u.String(), Schema: u.Query().Get("search_path")}
}

func getMigrationDirectory(migrationPath string) *migrate.LocalDir {
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		if err = os.MkdirAll(migrationPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	dir, err := migrate.NewLocalDir(migrationPath)
	if err != nil {
		log.Fatal(err)
	}

	err = migrate.Validate(dir)
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getDriver(ctx context.Context, postgresUrl *sqlclient.URL) (*postgres.Driver, func()) {
	client, err := sqlclient.OpenURL(ctx, postgresUrl.URL)
	if err != nil {
		log.Fatal(err)
	}
	drv := client.Driver.(*postgres.Driver)
	return drv, func() { client.Close() }
}
