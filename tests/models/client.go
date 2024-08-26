package models

import (
	"context"
	"fmt"
	"github.com/MrSametBurgazoglu/enterprise/client"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IDatabase interface {
	NewTransaction(ctx context.Context) (client.DatabaseTransactionClient, error)
	Exit()
	AddBeginHooks(...func())
	AddEndHooks(...func())
	client.DatabaseClient
}

type Database struct {
	*pgxpool.Pool
	BeginHooks []func()
	EndHooks   []func()
}

func (d *Database) AddBeginHooks(f ...func()) {
	d.BeginHooks = f
}

func (d *Database) AddEndHooks(f ...func()) {
	d.EndHooks = f
}

func (d *Database) BeginHook() {
	for _, hook := range d.BeginHooks {
		hook()
	}
}

func (d *Database) EndHook() {
	for _, hook := range d.EndHooks {
		hook()
	}
}

func (d *Database) SetupPostgres(dbUrl string) error {
	var err error
	d.Pool, err = pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("Unable to create connection pool: %v\n", err)
	}
	return nil
}

func (d *Database) NewTransaction(ctx context.Context) (client.DatabaseTransactionClient, error) {
	tx, err := d.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Transaction{Tx: tx}, nil
}

func (d *Database) Exit() {
	d.Pool.Close()
}

type Transaction struct {
	BeginHooks []func()
	EndHooks   []func()
	pgx.Tx
}

func (t Transaction) BeginHook() {
	for _, hook := range t.BeginHooks {
		hook()
	}
}

func (t Transaction) EndHook() {
	for _, hook := range t.EndHooks {
		hook()
	}
}

func NewDB(dbUrl string) (IDatabase, error) {
	d := &Database{}
	err := d.SetupPostgres(dbUrl)
	return d, err
}
