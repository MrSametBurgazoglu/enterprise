package enterprisetest

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os/exec"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// NewDB spins up an ephemeral PostgreSQL container using docker, returning the DB URL.
// The container is automatically destroyed when the test finishes using t.Cleanup.
func NewDB(t *testing.T) string {
	t.Helper()

	// Spin up a Postgres container with a random port to avoid collisions
	port := 10000 + rand.Intn(10000)
	containerName := fmt.Sprintf("enterprise-test-postgres-%d", port)

	// Run container in detached mode
	cmd := exec.Command("docker", "run",
		"--name", containerName,
		"-d",
		"-p", fmt.Sprintf("%d:5432", port),
		"-e", "POSTGRES_PASSWORD=password",
		"-e", "POSTGRES_DB=enterprise",
		"postgres:15-alpine",
	)

	if err := cmd.Run(); err != nil {
		t.Fatalf("enterprisetest: failed to start postgres container: %v", err)
	}

	// Register cleanup to kill and remove the container
	t.Cleanup(func() {
		_ = exec.Command("docker", "rm", "-f", containerName).Run()
	})

	dbURL := fmt.Sprintf("postgres://postgres:password@localhost:%d/enterprise?sslmode=disable&search_path=public", port)

	// Poll until database is accepting connections
	var db *sql.DB
	var err error
	for i := 0; i < 40; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				db.Close()
				return dbURL
			}
		}
		if db != nil {
			db.Close()
		}
		time.Sleep(250 * time.Millisecond)
	}

	t.Fatalf("enterprisetest: postgres container failed to become ready: %v", err)
	return ""
}
