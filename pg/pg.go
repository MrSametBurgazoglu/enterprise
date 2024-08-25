package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PG struct {
	Pool *pgxpool.Pool
}

func (p *PG) SetupPostgres(dbUrl string) error {
	var err error
	p.Pool, err = pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("Unable to create connection pool: %v\n", err)
	}
	return nil
}

func (p *PG) Exit() {
	p.Pool.Close()
}
