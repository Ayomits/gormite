package gormite_databases

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"os"
	"reflect"
	"strings"
	"time"
)

// TODO: Создать утилитарную библиотеку!!!

type PgxWrappedDatabase interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (
		pgconn.CommandTag,
		error,
	)
}

type PostgresDatabase struct {
	PgX       PgxWrappedDatabase
	PgxConfig *pgxpool.Config
	// pgxConn - fallback, cannot be replaced, used for close connection ability
	pgxConn *PgXWrapper
	logger  *log.Logger
}

// PostgresDatabaseInterface - TODO: Remove after move to utility library
type PostgresDatabaseInterface interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Select(sql string, args ...interface{}) QueryInterface
	Get(sql string, args ...interface{}) QueryInterface
}

// todo: сделать интерфейс-обертку чтобы можно было использовать транзакцию в качестве основы

func NewPostgresDatabase(ctx context.Context, dsn string) *PostgresDatabase {
	// TODO: Move to util func
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

	return &PostgresDatabase{PgX: conn, PgxConfig: config, pgxConn: conn, logger: logger}
}

type PgXWrapper struct {
	*pgxpool.Pool
	logger *log.Logger
}

func trimSQL(sql string) string {
	lines := strings.Split(sql, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	sql = strings.Join(lines, " ")

	return sql
}

type SqlError struct {
	err  error
	sql  string
	args []any
}

func (e *SqlError) Error() string {
	return e.err.Error()
}

func (e *SqlError) GetExtraData() map[string]any {
	return map[string]any{
		"sql":  e.sql,
		"args": e.args,
	}
}

func (w *PgXWrapper) Exec(
	ctx context.Context,
	sqlQuery string,
	args ...interface{},
) (pgconn.CommandTag, error) {
	v, err := w.Pool.Exec(ctx, sqlQuery, args...)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		w.logger.Warn(err.Error(), "sql", trimSQL(sqlQuery), "args", args)
	}

	return v, err
}

func (d *PostgresDatabase) Select(sql string, args ...interface{}) QueryInterface {
	return &Query{db: d.PgX, sql: sql, args: args, logger: d.logger}
}

func (d *PostgresDatabase) Get(sql string, args ...interface{}) QueryInterface {
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
	args ...interface{},
) (pgconn.CommandTag, error) {
	tag, err := d.PgX.Exec(ctx, sql, args...)

	return tag, errors.WithStack(err)
}

func (d *PostgresDatabase) Destruct() {
	d.pgxConn.Close()
}

type Query struct {
	db     PgxWrappedDatabase
	logger *log.Logger

	sql       string
	args      []interface{}
	scan      []interface{}
	scanFirst bool
	scanCol   bool
}
type QueryInterface interface {
	Scan(dest ...interface{}) QueryInterface
	ScanCol(dest ...interface{}) QueryInterface
	Exec(ctx context.Context) error
}

func (q *Query) Scan(dest ...interface{}) QueryInterface {
	q.scan = append(q.scan, dest...)

	return q
}

func (q *Query) ScanCol(dest ...interface{}) QueryInterface {
	q.scan = append(q.scan, dest...)
	q.scanCol = true

	return q
}

type queryError struct {
	err   error
	query *Query
}

func (e *queryError) Error() string {
	return e.err.Error()
}

func (q *Query) Exec(ctx context.Context) error {
	err := q.ExecWrapped(ctx)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
			return sql.ErrNoRows
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

func scanAll(
	rows pgx.Rows,
	columnsNum int,
	positionsList [][][]int,
	dest []interface{},
) error {
	anyFound := false
	for rows.Next() {
		// Possible problem when only one struct used but results rows more than 1
		// All will be mashed into 1 entity, the last one will remain, respectively
		anyFound = true
		err := scanStructs(rows, columnsNum, positionsList, dest)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if anyFound {
		return nil
	}

	rowsErr := rows.Err()
	if rowsErr == nil {
		rowsErr = sql.ErrNoRows
	}

	return errors.WithStack(rowsErr)
}

func scanFirst(
	rows pgx.Rows,
	columnsNum int,
	positionsList [][][]int,
	dest []interface{},
) error {
	if rows.Next() {
		err := scanStructs(rows, columnsNum, positionsList, dest)
		if err != nil {
			return err
		}

		return nil
	}

	if err := rows.Err(); err != nil {
		return errors.WithStack(err)
	}

	return sql.ErrNoRows
}

func scanStructs(
	rows pgx.Rows,
	lenColumns int,
	positionsList [][][]int,
	dest []interface{},
) error {
	allValues := make([]interface{}, 0, lenColumns)
	items := make([]reflect.Value, 0, len(dest))

	for i, dest := range dest {
		values, item := structValues(positionsList[i], dest)
		allValues = append(allValues, values...)
		items = append(items, item)
	}

	err := rows.Scan(allValues...)
	if err != nil {
		return errors.WithStack(err)
	}

	for i, item := range items {
		if reflect.ValueOf(dest[i]).Elem().Kind() != reflect.Slice { // struct or slice
			reflect.ValueOf(dest[i]).Elem().Set(item.Elem())

			continue
		}

		slice := reflect.ValueOf(dest[i]).Elem()

		if slice.Type().Elem().Kind() == reflect.Pointer {
			slice.Set(reflect.Append(slice, item))
		} else {
			slice.Set(reflect.Append(slice, item.Elem()))
		}
	}

	return nil
}

func structValues(
	positions [][]int,
	dest interface{},
) ([]interface{}, reflect.Value) {
	var values []interface{}

	structType := reflect.TypeOf(dest).Elem()
	if structType.Kind() == reflect.Slice { // struct or slice
		structType = structType.Elem()
	}

	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}

	item := reflect.New(structType)

	for _, p := range positions {
		values = append(values, getField(item.Elem(), p).Addr().Interface())
	}

	return values, item
}

func getField(v reflect.Value, p []int) reflect.Value {
	if len(p) == 1 {
		return v.Field(p[0])
	}

	return getField(v.Field(p[0]), p[1:])
}

func getPositionsList(
	columns []pgconn.FieldDescription,
	dest []interface{},
) ([][][]int, error) {
	positionsList := make([][][]int, 0, len(dest))

	var missedColumns []pgconn.FieldDescription

	for _, d := range dest {
		destType := reflect.TypeOf(d)
		if destType.Kind() != reflect.Pointer {
			return nil, fmt.Errorf("destination must be a pointer")
		}

		var structType reflect.Type

		sliceType := destType.Elem()
		switch { // struct or slice
		case sliceType.Kind() == reflect.Slice:
			structType = sliceType.Elem()
			if structType.Kind() == reflect.Pointer {
				structType = structType.Elem()
			}
		case sliceType.Kind() == reflect.Struct:
			structType = sliceType
		default:
			return nil, fmt.Errorf("destination must be a pointer of slice")
		}

		if structType.Kind() != reflect.Struct {
			return nil, fmt.Errorf("destination must be a pointer of slice of struct or pointer struct")
		}

		positions := getPositions(columns, structType)
		positionsList = append(positionsList, positions)
		for i, position := range positions {
			if position == nil {
				missedColumns = append(missedColumns, columns[i])
			}
		}
	}

	if len(missedColumns) != 0 {
		columnsStrings := make([]string, 0, len(missedColumns))
		for _, c := range missedColumns {
			columnsStrings = append(columnsStrings, c.Name)
		}
		return nil, fmt.Errorf("%d columns not found in destination structs(%s)", len(missedColumns), strings.Join(columnsStrings, ", "))
	}

	return positionsList, nil
}

func getPositions(columns []pgconn.FieldDescription, t reflect.Type) [][]int {
	positions := make([][]int, 0, t.NumField())
	discoverStruct(columns, &positions, t, nil)

	return positions
}

func discoverStruct(
	allColumns []pgconn.FieldDescription,
	positions *[][]int,
	t reflect.Type,
	prefix []int,
) int {
	dbFieldsMap := make(map[string][]int, t.NumField())
	columns := allColumns

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		name := f.Tag.Get("db")

		if !f.IsExported() || name == "-" {
			continue
		}

		if name == "" {
			name = f.Name
		}

		if f.Type.Kind() == reflect.Struct && (f.Anonymous || name == "*") {
			columns = columns[len(dbFieldsMap):]
			used := discoverStruct(
				columns,
				positions,
				f.Type,
				append(prefix, i),
			)
			columns = columns[used:]

			continue
		}

		dbFieldsMap[name] = append(prefix, i)
	}

	// no columns(invalid sql)
	if columns == nil {
		return 0
	}

	for i := 0; i < len(dbFieldsMap); i++ {
		*positions = append(*positions, dbFieldsMap[columns[i].Name])
	}

	return len(dbFieldsMap)
}
