package gormite_databases

import (
	"context"
	gdh "github.com/KoNekoD/gormite/pkg/gormite_databases_helpers"
	"github.com/KoNekoD/pgx-colon-query-rewriter/pkg/pgxcqr"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"os"
	"time"
)

type PgxWrappedDatabase interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (
		pgconn.CommandTag,
		error,
	)
}

type PostgresDatabase struct {
	PgX       PgxWrappedDatabase
	PgxConfig *pgxpool.Config

	pgxConn *PgXWrapper
	logger  *log.Logger
}

func NewPostgresDatabase(ctx context.Context, dsn string) *PostgresDatabase {
	module := "storage"
	opts := log.Options{
		ReportTimestamp: true,
		Prefix:          module,
		TimeFormat:      time.DateTime,
		Level:           log.DebugLevel,
	}
	logger := log.NewWithOptions(os.Stdout, opts)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Fatal("Cannot parse config", "err", err)
	}

	pgXPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	conn := &PgXWrapper{Pool: pgXPool, logger: logger}

	return &PostgresDatabase{
		PgX:       conn,
		PgxConfig: config,
		pgxConn:   conn,
		logger:    logger,
	}
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

func (d *PostgresDatabase) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (gdh.CommandTag, error) {
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
