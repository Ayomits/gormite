package gormite_databases

import (
	"context"
	gdh "github.com/KoNekoD/gormite/pkg/gormite_databases_helpers"
	"github.com/KoNekoD/gormite/pkg/utils"
	"github.com/KoNekoD/pgx-colon-query-rewriter/pkg/pgxcqr"
	"github.com/charmbracelet/log"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type newPostgresDatabaseOption struct {
	onError func(method string, err error, sql string, args ...any)
}

type NewPostgresDatabaseOptionFn func(o *newPostgresDatabaseOption)

func WithOnError(onError func(method string, err error, sql string, args ...any)) NewPostgresDatabaseOptionFn {
	return func(o *newPostgresDatabaseOption) {
		o.onError = onError
	}
}

type PgxWrappedDatabase interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type PostgresDatabase struct {
	PgX       PgxWrappedDatabase
	PgxConfig *pgxpool.Config

	pgxConn *PgxWrapper
	logger  *log.Logger
}

func NewPostgresDatabase(ctx context.Context, dsn string, opts ...NewPostgresDatabaseOptionFn) *PostgresDatabase {
	logger := utils.NewLogger("storage")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Fatal("Cannot parse config", "err", err)
	}

	pgxPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	conn := &PgxWrapper{Pool: pgxPool, logger: logger}

	onError := func(method string, err error, sql string, args ...any) {
		logger.Warn(err.Error(), "sql", sql, "args", args)
	}

	o := &newPostgresDatabaseOption{onError: onError}
	for _, opt := range opts {
		opt(o)
	}

	conn.onErr = o.onError

	return &PostgresDatabase{PgX: conn, PgxConfig: config, pgxConn: conn, logger: logger}
}

func (d *PostgresDatabase) WrapInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	opts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

	tx, err := d.PgX.(*PgxWrapper).BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	orig := d.PgX
	d.PgX = tx
	defer func() { d.PgX = orig }()

	if err = fn(ctx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return multierror.Append(err, errors.Wrap(rollbackErr, "failed to rollback transaction"))
		}

		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (d *PostgresDatabase) Select(sql string, args ...any) gdh.QueryInterface {
	return &Query{db: d.PgX, sql: sql, args: args, logger: d.logger}
}

func (d *PostgresDatabase) Get(sql string, args ...any) gdh.QueryInterface {
	return &Query{
		db:        d.PgX,
		sql:       sql,
		args:      args,
		scanFirst: true,
		logger:    d.logger,
	}
}

func (d *PostgresDatabase) Exec(ctx context.Context, sql string, args ...any) (gdh.CommandTag, error) {
	tag, err := d.PgX.Exec(ctx, sql, args...)

	return tag, errors.WithStack(err)
}

func (d *PostgresDatabase) Query(ctx context.Context, sql string, args ...any) (gdh.Rows, error) {
	rows, err := d.PgX.Query(ctx, sql, args...)

	return rows, errors.WithStack(err)
}

func (d *PostgresDatabase) GetNamedArgs(args any) any {
	return pgxcqr.NamedArgs(args.(map[string]any))
}

func (d *PostgresDatabase) Destruct() {
	d.pgxConn.Close()
}
