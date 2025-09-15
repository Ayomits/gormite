package gormite_databases

import (
	"context"
	databaseSql "database/sql"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PgxWrapper struct {
	*pgxpool.Pool
	logger *log.Logger
	onErr  func(method string, err error, sql string, args ...any)
}

func (w *PgxWrapper) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	v, err := w.Pool.Exec(ctx, sql, args...)

	if err != nil && !errors.Is(err, databaseSql.ErrNoRows) {
		w.onErr("Exec", err, trimSQL(sql), args...)
	}

	return v, err
}

func (w *PgxWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	v, err := w.Pool.Query(ctx, sql, args...)

	if err != nil && !errors.Is(err, databaseSql.ErrNoRows) {
		w.onErr("Query", err, trimSQL(sql), args...)
	}

	return v, err
}
