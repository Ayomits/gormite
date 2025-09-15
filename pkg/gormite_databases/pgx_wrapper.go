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

type PgXWrapper struct {
	*pgxpool.Pool
	logger *log.Logger
}

func (w *PgXWrapper) Exec(
	ctx context.Context,
	sqlQuery string,
	args ...any,
) (pgconn.CommandTag, error) {
	v, err := w.Pool.Exec(ctx, sqlQuery, args...)

	if err != nil && !errors.Is(err, databaseSql.ErrNoRows) {
		w.logger.Warn(err.Error(), "sql", trimSQL(sqlQuery), "args", args)
	}

	return v, err
}

func (w *PgXWrapper) Query(
	ctx context.Context,
	sqlQuery string,
	args ...any,
) (pgx.Rows, error) {
	v, err := w.Pool.Query(ctx, sqlQuery, args...)

	if err != nil && !errors.Is(err, databaseSql.ErrNoRows) {
		w.logger.Warn(err.Error(), "sql", trimSQL(sqlQuery), "args", args)
	}

	return v, err
}
