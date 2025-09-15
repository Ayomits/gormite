package gormite_databases

import (
	"context"
	databaseSql "database/sql"
	gdh "github.com/KoNekoD/gormite/pkg/gormite_databases_helpers"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

type Query struct {
	db     PgxWrappedDatabase
	logger *log.Logger

	sql       string
	args      []any
	scan      []any
	scanFirst bool
	scanCol   bool
}

func (q *Query) Scan(dest ...any) gdh.QueryInterface {
	q.scan = append(q.scan, dest...)

	return q
}

func (q *Query) ScanCol(dest ...any) gdh.QueryInterface {
	q.scan = append(q.scan, dest...)
	q.scanCol = true

	return q
}

func (q *Query) Exec(ctx context.Context) error {
	err := q.ExecWrapped(ctx)

	if err != nil && !errors.Is(err, databaseSql.ErrNoRows) {
		err = &queryError{err: err, query: q}

		q.logger.Warn(err.Error(), "sql", trimSQL(q.sql), "args", q.args)
	}

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (q *Query) ExecWrapped(ctx context.Context) error {
	rows, err := q.db.Query(ctx, q.sql, q.args...)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := rows.Err(); err != nil {
		return errors.WithStack(err)
	}

	defer rows.Close()

	if q.scanCol {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return errors.WithStack(err)
			}
			return databaseSql.ErrNoRows
		}
		err = rows.Scan(q.scan...)
		if err != nil {
			return err
		}
		return nil
	}

	columns := rows.FieldDescriptions()

	positionsList, err := getPositionsList(columns, q.scan)
	if err != nil {
		return err
	}

	if q.scanFirst {
		return scanFirst(rows, len(columns), positionsList, q.scan)
	}

	return scanAll(rows, len(columns), positionsList, q.scan)
}
