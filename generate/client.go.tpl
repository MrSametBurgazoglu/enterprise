package models

import (
	"context"
	"fmt"
	"github.com/MrSametBurgazoglu/enterprise/client"
	"github.com/MrSametBurgazoglu/enterprise/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
)

type Options struct {
	Url   string
	Debug bool
}

type IDatabase interface {
	NewTransaction(ctx context.Context, options ...pgx.TxOptions) (client.DatabaseTransactionClient, error)
	Exit()
	client.DatabaseClient
}

type Database struct {
	pool    *pgxpool.Pool
	Options *Options
	Logger  logger.Logger
}

func (d *Database) Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
	commandTag, err = d.pool.Exec(ctx, sql, arguments)
	if err != nil {
		d.Logger.LogError(ctx, sql, arguments[0].(pgx.NamedArgs), err)
	} else if d.Options.Debug {
		d.Logger.Log(ctx, sql, arguments[0].(pgx.NamedArgs))
	}
	return commandTag, err
}

func (d *Database) Query(ctx context.Context, sql string, args ...any) (rows pgx.Rows, err error) {
	rows, err = d.pool.Query(ctx, sql, args)
	if err != nil {
		d.Logger.LogError(ctx, sql, args[0].(pgx.NamedArgs), err)
	} else if d.Options.Debug {
		d.Logger.Log(ctx, sql, args[0].(pgx.NamedArgs))
	}
	return rows, err
}

func (d *Database) QueryRow(ctx context.Context, sql string, args ...any) (row pgx.Row) {
	row = d.pool.QueryRow(ctx, sql, args)
	if d.Options.Debug {
		d.Logger.Log(ctx, sql, args[0].(pgx.NamedArgs))
	}
	return row
}

func (d *Database) SetupPostgres(dbUrl string) error {
	var err error
	d.pool, err = pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("Unable to create connection pool: %v\n", err)
	}
	return nil
}

func (d *Database) NewTransaction(ctx context.Context, options ...pgx.TxOptions) (client.DatabaseTransactionClient, error) {
	var (
		err error
		tx  pgx.Tx
	)

	if len(options) == 0 {
		tx, err = d.pool.Begin(ctx)
	} else {
		tx, err = d.pool.BeginTx(ctx, options[0])
	}

	if err != nil {
		return nil, err
	}

	return &Transaction{Tx: tx}, nil
}

func (d *Database) Exit() {
	d.pool.Close()
}

type Transaction struct {
	Logger logger.Logger
	pgx.Tx
}

func (t Transaction) GetLogger() logger.Logger {
	return t.Logger
}

func (t Transaction) SavePoint(ctx context.Context) (client.DatabaseTransactionClient, error) {
	tx, err := t.Tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Transaction{Tx: tx, Logger: t.Logger}, nil
}

func NewDB(options *Options) (IDatabase, error) {
	d := &Database{Options: options}
	err := d.SetupPostgres(options.Url)
	slogger := initializeSlog(slog.LevelDebug)
	d.Logger = logger.NewLogger(slogger)
	return d, err
}

func initializeSlog(logLevel slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}
