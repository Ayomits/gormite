package postgres_platform

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/assets"
	"github.com/KoNekoD/gormite/pkg/diff_dtos"
	"github.com/KoNekoD/gormite/pkg/enums"
	"github.com/KoNekoD/gormite/pkg/keywords"
	"github.com/KoNekoD/gormite/pkg/platforms"
	"github.com/KoNekoD/gormite/pkg/schema_managers"
	"github.com/KoNekoD/gormite/pkg/schema_managers/postgres_schema_manager"
	"github.com/elliotchance/pie/v2"
	"slices"
	"strconv"
	"strings"
)

type PostgreSQLPlatform struct {
	*platforms.AbstractPlatform
	UseBooleanTrueFalseStrings bool
	BooleanLiterals            map[string][]string
}

func NewPostgreSQLPlatform() *PostgreSQLPlatform {
	v := &PostgreSQLPlatform{}

	v.AbstractPlatform = platforms.NewAbstractPlatform(v)

	// PostgreSQL booleans literals
	v.BooleanLiterals = map[string][]string{
		"true":  {"t", "true", "y", "yes", "on", "1"},
		"false": {"f", "false", "n", "no", "off", "0"},
	}

	return v
}

func (p *PostgreSQLPlatform) SetUseBooleanTrueFalseStrings(flag bool) {
	p.UseBooleanTrueFalseStrings = flag
}
func (p *PostgreSQLPlatform) GetRegexpExpression() string {
	return `SIMILAR TO`
}
func (p *PostgreSQLPlatform) GetLocateExpression(
	string string,
	substring string,
	start *string,
) string {
	if start != nil {
		startStr := *start
		string = p.GetSubstringExpression(string, *start, nil)

		return `CASE WHEN (POSITION(` + substring + ` IN ` + string + `) = 0) THEN 0 ELSE (POSITION(` + substring + ` IN ` + string + `) + ` + startStr + ` - 1) END`
	}

	return fmt.Sprintf(`POSITION(%s IN %s)`, substring, string)
}
func (p *PostgreSQLPlatform) GetDateArithmeticIntervalExpression(
	date string,
	operator string,
	interval string,
	unit enums.DateIntervalUnit,
) string {
	if unit == enums.DateIntervalUnitQuarter {
		interval = p.MultiplyInterval(interval, 3)
		unit = enums.DateIntervalUnitMonth
	}

	return `(` + date + ` ` + operator + ` (` + interval + " || ` " + string(unit) + "`)::interval)"
}
func (p *PostgreSQLPlatform) GetDateDiffExpression(
	date1 string,
	date2 string,
) string {
	return `(DATE(` + date1 + `)-DATE(` + date2 + `))`
}
func (p *PostgreSQLPlatform) GetCurrentDatabaseExpression() string {
	return `CURRENT_DATABASE()`
}
func (p *PostgreSQLPlatform) SupportsSequences() bool {
	return true
}
func (p *PostgreSQLPlatform) SupportsSchemas() bool {
	return true
}
func (p *PostgreSQLPlatform) SupportsIdentityColumns() bool {
	return true
}

func (p *PostgreSQLPlatform) SupportsPartialIndexes() bool {
	return true
}

func (p *PostgreSQLPlatform) SupportsCommentOnStatement() bool {
	return true
}

func (p *PostgreSQLPlatform) GetListDatabasesSQL() string {
	return `SELECT datname FROM pg_database`
}

func (p *PostgreSQLPlatform) GetListSequencesSQL(database string) string {
	return `SELECT sequence_name AS relname,
                       sequence_schema AS schemaname,
                       minimum_value AS min_value,
                       increment AS increment_by
                FROM   information_schema.sequences
                WHERE  sequence_catalog = ` + p.QuoteStringLiteral(database) + " AND    sequence_schema NOT LIKE 'pg\\_%' AND    sequence_schema != 'information_schema'"
}

func (p *PostgreSQLPlatform) GetListViewsSQL(database string) string {
	return `SELECT quote_ident(table_name) AS viewname,
                       table_schema AS schemaname,
                       view_definition AS definition
                FROM   information_schema.views
                WHERE  view_definition IS NOT NULL`
}

func (p *PostgreSQLPlatform) GetAdvancedForeignKeyOptionsSQL(foreignKey *assets.ForeignKeyConstraint) string {
	query := ``

	if foreignKey.HasOption(`match`) {
		query += ` MATCH ` + foreignKey.GetOption(`match`).(string)
	}

	query += p.AbstractPlatform.GetAdvancedForeignKeyOptionsSQL(foreignKey)

	if foreignKey.HasOption(`deferrable`) && foreignKey.GetOption(`deferrable`) != false {
		query += ` DEFERRABLE`
	} else {
		query += ` NOT DEFERRABLE`
	}

	if foreignKey.HasOption(`deferred`) && foreignKey.GetOption(`deferred`) != false {
		query += ` INITIALLY DEFERRED`
	} else {
		query += ` INITIALLY IMMEDIATE`
	}

	return query
}
func (p *PostgreSQLPlatform) GetAlterTableSQL(diff *diff_dtos.TableDiff) []string {
	sql := make([]string, 0)
	commentsSQL := make([]string, 0)

	table := diff.GetOldTable()

	tableNameSQL := table.GetQuotedName(p)

	for _, addedColumn := range diff.GetAddedColumns() {
		query := `ADD ` + p.GetColumnDeclarationSQL(
			addedColumn.GetQuotedName(p),
			addedColumn.ToArray(),
		)

		sql = append(sql, `ALTER TABLE `+tableNameSQL+` `+query)

		comment := addedColumn.GetComment()

		if comment == `` {
			continue
		}

		commentsSQL = append(
			commentsSQL, p.GetCommentOnColumnSQL(
				tableNameSQL,
				addedColumn.GetQuotedName(p),
				comment,
			),
		)
	}

	for _, droppedColumn := range diff.GetDroppedColumns() {
		query := `DROP ` + droppedColumn.GetQuotedName(p)
		sql = append(sql, `ALTER TABLE `+tableNameSQL+` `+query)
	}

	for _, columnDiff := range diff.GetChangedColumns() {
		oldColumn := columnDiff.GetOldColumn()
		newColumn := columnDiff.GetNewColumn()

		oldColumnName := oldColumn.GetQuotedName(p)
		newColumnName := newColumn.GetQuotedName(p)

		if columnDiff.HasNameChanged() {
			sql = append(
				sql,
				p.GetRenameColumnSQL(
					tableNameSQL,
					oldColumnName,
					newColumnName,
				)...,
			)
		}

		if columnDiff.HasTypeChanged() ||
			columnDiff.HasPrecisionChanged() ||
			columnDiff.HasScaleChanged() ||
			columnDiff.HasFixedChanged() ||
			columnDiff.HasLengthChanged() {
			typeVar := newColumn.GetColumnType()

			// SERIAL/BIGSERIAL are not "real" types and we can`t alter a column to that type
			columnDefinition := newColumn.ToArray()
			columnDefinition[`autoincrement`] = false

			// here was a server version check before, but DBAL API does not support this anymore.
			query := `ALTER ` + newColumnName + ` TYPE ` + typeVar.GetSQLDeclaration(
				columnDefinition,
				p,
			)
			sql = append(sql, `ALTER TABLE `+tableNameSQL+` `+query)
		}

		if columnDiff.HasDefaultChanged() {
			defaultClause := ""
			if newColumn.GetColumnDefault() == nil {
				defaultClause = ` DROP DEFAULT`
			} else {
				defaultClause = ` SET` + p.GetDefaultValueDeclarationSQL(newColumn.ToArray())
			}

			query := `ALTER ` + newColumnName + defaultClause
			sql = append(sql, `ALTER TABLE `+tableNameSQL+` `+query)
		}

		if columnDiff.HasNotNullChanged() {
			action := ""
			if newColumn.GetNotNull() {
				action = `SET`
			} else {
				action = `DROP`
			}
			query := `ALTER ` + newColumnName + ` ` + action + ` NOT NULL`
			sql = append(sql, `ALTER TABLE `+tableNameSQL+` `+query)
		}

		if columnDiff.HasAutoIncrementChanged() {
			query := ""
			if newColumn.GetAutoincrement() {
				query = `ADD GENERATED BY DEFAULT AS IDENTITY`
			} else {
				query = `DROP IDENTITY`
			}

			sql = append(
				sql,
				`ALTER TABLE `+tableNameSQL+` ALTER `+newColumnName+` `+query,
			)
		}

		if !columnDiff.HasCommentChanged() {
			continue
		}

		commentsSQL = append(
			commentsSQL, p.GetCommentOnColumnSQL(
				tableNameSQL,
				newColumn.GetQuotedName(p),
				newColumn.GetComment(),
			),
		)
	}

	sql = append(p.GetPreAlterTableIndexForeignKeySQL(diff), sql...)
	sql = append(sql, commentsSQL...)
	sql = append(sql, p.GetPostAlterTableIndexForeignKeySQL(diff)...)

	return sql
}
func (p *PostgreSQLPlatform) GetRenameIndexSQL(
	oldIndexName string,
	index *assets.Index,
	tableName string,
) []string {
	if strings.Contains(tableName, `.`) {
		schema := strings.Split(tableName, ".")[0]
		oldIndexName = schema + `.` + oldIndexName
	}

	return []string{`ALTER INDEX ` + oldIndexName + ` RENAME TO ` + index.GetQuotedName(p)}
}
func (p *PostgreSQLPlatform) GetCreateSequenceSQL(sequence *assets.Sequence) string {
	return `CREATE SEQUENCE ` + sequence.GetQuotedName(p) +
		` INCREMENT BY ` + strconv.Itoa(sequence.GetAllocationSize()) +
		` MINVALUE ` + strconv.Itoa(sequence.GetInitialValue()) +
		` START ` + strconv.Itoa(sequence.GetInitialValue()) +
		p.GetSequenceCacheSQL(sequence)
}
func (p *PostgreSQLPlatform) GetAlterSequenceSQL(sequence *assets.Sequence) string {
	return `ALTER SEQUENCE ` + sequence.GetQuotedName(p) +
		` INCREMENT BY ` + strconv.Itoa(sequence.GetAllocationSize()) +
		p.GetSequenceCacheSQL(sequence)
}
func (p *PostgreSQLPlatform) GetSequenceCacheSQL(sequence *assets.Sequence) string {
	cache := sequence.GetCache()
	if cache != nil && *cache > 1 {
		cacheInt := *cache
		return ` CACHE ` + strconv.Itoa(cacheInt)
	}

	return ``
}
func (p *PostgreSQLPlatform) GetDropSequenceSQL(name string) string {
	return p.AbstractPlatform.GetDropSequenceSQL(name) + ` CASCADE`
}
func (p *PostgreSQLPlatform) GetDropForeignKeySQL(
	foreignKey string,
	table string,
) string {
	return p.GetDropConstraintSQL(foreignKey, table)
}
func (p *PostgreSQLPlatform) GetDropIndexSQL(name string, table string) string {
	if name == `"primary"` {
		constraintName := table + `_pkey`

		return p.GetDropConstraintSQL(constraintName, table)
	}

	return p.AbstractPlatform.GetDropIndexSQL(name, table)
}
func (p *PostgreSQLPlatform) _getCreateTableSQL(
	name string,
	columns []map[string]any,
	options map[string]any,
) []string {
	queryFields := p.GetColumnDeclarationListSQL(columns)

	if v, ok := options[`primary`]; ok {
		keyColumns := pie.Unique(v.([]string))
		queryFields += `, PRIMARY KEY(` + strings.Join(keyColumns, `, `) + `)`
	}

	unloggedV, unloggedOk := options[`unlogged`]
	unlogged := ""
	if unloggedOk && unloggedV == true {
		unlogged = ` UNLOGGED`
	} else {
		unlogged = ``
	}
	query := `CREATE` + unlogged + ` TABLE ` + name + ` (` + queryFields + `)`

	sql := []string{query}

	if v, ok := options[`indexes`]; ok {
		for _, index := range v.(map[string]*assets.Index) {
			sql = append(sql, p.GetCreateIndexSQL(index, name))
		}
	}

	if v, ok := options[`uniqueConstraints`]; ok {
		for _, uniqueConstraint := range v.(map[string]*assets.UniqueConstraint) {
			sql = append(
				sql,
				p.GetCreateUniqueConstraintSQL(uniqueConstraint, name),
			)
		}
	}

	if v, ok := options[`foreignKeys`]; ok {
		for _, definition := range v.([]*assets.ForeignKeyConstraint) {
			sql = append(sql, p.GetCreateForeignKeySQL(definition, name))
		}
	}

	return sql
}
func (p *PostgreSQLPlatform) ConvertSingleBooleanValue(
	value any,
	callback func(any) any,
) any {
	if value == nil {
		return callback(nil)
	}

	_, ok1 := value.(bool)
	_, ok2 := value.(int)
	if ok1 || ok2 {
		return callback(value.(bool))
	}

	str, ok3 := value.(string)
	if !ok3 {
		return callback(true)
	}

	if slices.Contains(
		p.BooleanLiterals[`false`],
		strings.ToLower(strings.TrimSpace(str)),
	) {
		return callback(false)
	}

	if slices.Contains(
		p.BooleanLiterals[`true`],
		strings.ToLower(strings.TrimSpace(str)),
	) {
		return callback(true)
	}

	panic(fmt.Errorf(`unrecognized boolean literal, %s given`, value))
}
func (p *PostgreSQLPlatform) DoConvertBooleans(
	item any,
	callback func(any) any,
) any {
	if itemMap, ok := item.(map[string]any); ok {
		for key, value := range itemMap {
			itemMap[key] = p.ConvertSingleBooleanValue(value, callback)
		}

		return itemMap
	}

	return p.ConvertSingleBooleanValue(item, callback)
}

func (p *PostgreSQLPlatform) ConvertBooleans(item any) any {
	if !p.UseBooleanTrueFalseStrings {
		return p.AbstractPlatform.ConvertBooleans(item)
	}

	return p.DoConvertBooleans(
		item,
		func(value any) any {
			if value == nil {
				return `nil`
			}

			if value == true {
				return `true`
			} else {
				return `false`
			}
		},
	)
}
func (p *PostgreSQLPlatform) ConvertBooleansToDatabaseValue(item any) any {
	if !p.UseBooleanTrueFalseStrings {
		return p.AbstractPlatform.ConvertBooleansToDatabaseValue(item)
	}

	return p.DoConvertBooleans(
		item,
		func(value any) any {
			if value == nil {
				return nil
			}

			valueInt := value.(int)

			return &valueInt
		},
	)
}
func (p *PostgreSQLPlatform) ConvertFromBoolean(item any) *bool {
	if slices.Contains(p.BooleanLiterals[`false`], item.(string)) {
		f := false
		return &f
	}

	return p.AbstractPlatform.ConvertFromBoolean(item)
}
func (p *PostgreSQLPlatform) GetSequenceNextValSQL(sequence string) string {
	return "SELECT NEXTVAL(`" + sequence + "`)"
}
func (p *PostgreSQLPlatform) GetBooleanTypeDeclarationSQL(column map[string]any) string {
	return `BOOLEAN`
}
func (p *PostgreSQLPlatform) GetIntegerTypeDeclarationSQL(column map[string]any) string {
	return `INT` + p.GetCommonIntegerTypeDeclarationSQL(column)
}
func (p *PostgreSQLPlatform) GetBigIntTypeDeclarationSQL(column map[string]any) string {
	return `BIGINT` + p.GetCommonIntegerTypeDeclarationSQL(column)
}
func (p *PostgreSQLPlatform) GetSmallIntTypeDeclarationSQL(column map[string]any) string {
	return `SMALLINT` + p.GetCommonIntegerTypeDeclarationSQL(column)
}
func (p *PostgreSQLPlatform) GetGuidTypeDeclarationSQL(column map[string]any) string {
	return `UUID`
}
func (p *PostgreSQLPlatform) GetDateTimeTypeDeclarationSQL(column map[string]any) string {
	return `TIMESTAMP(0) WITHOUT TIME ZONE`
}
func (p *PostgreSQLPlatform) GetDateTimeTzTypeDeclarationSQL(column map[string]any) string {
	return `TIMESTAMP(0) WITH TIME ZONE`
}
func (p *PostgreSQLPlatform) GetDateTypeDeclarationSQL(column map[string]any) string {
	return `DATE`
}
func (p *PostgreSQLPlatform) GetTimeTypeDeclarationSQL(column map[string]any) string {
	return `TIME(0) WITHOUT TIME ZONE`
}
func (p *PostgreSQLPlatform) GetCommonIntegerTypeDeclarationSQL(column map[string]any) string {
	if v, ok := column[`autoincrement`]; ok && v == true {
		return ` GENERATED BY DEFAULT AS IDENTITY`
	}

	return ``
}
func (p *PostgreSQLPlatform) GetVarcharTypeDeclarationSQLSnippet(length *int) string {
	sql := `VARCHAR`

	if length != nil {
		sql += fmt.Sprintf(`(%d)`, *length)
	}

	return sql
}
func (p *PostgreSQLPlatform) GetBinaryTypeDeclarationSQLSnippet(length *int) string {
	return `BYTEA`
}
func (p *PostgreSQLPlatform) GetVarbinaryTypeDeclarationSQLSnippet(length *int) string {
	return `BYTEA`
}
func (p *PostgreSQLPlatform) GetClobTypeDeclarationSQL(column map[string]any) string {
	return `TEXT`
}
func (p *PostgreSQLPlatform) GetDateTimeTzFormatString() string {
	return `Y-m-d H:i:sO`
}
func (p *PostgreSQLPlatform) GetEmptyIdentityInsertSQL(
	quotedTableName string,
	quotedIdentifierColumnName string,
) string {
	return `INSERT ` + `INTO ` + quotedTableName + ` (` + quotedIdentifierColumnName + `) VALUES (DEFAULT)`
}
func (p *PostgreSQLPlatform) GetTruncateTableSQL(
	tableName string,
	cascade bool,
) string {
	tableIdentifier := assets.NewIdentifier(tableName)
	sql := `TRUNCATE ` + tableIdentifier.GetQuotedName(p)

	if cascade {
		sql += ` CASCADE`
	}

	return sql
}
func (p *PostgreSQLPlatform) GetDefaultColumnValueSQLSnippet() string {
	return `
SELECT pg_get_expr(adbin, adrelid)
             FROM pg_attrdef
             WHERE c.oid = pg_attrdef.adrelid
                AND pg_attrdef.adnum=a.attnum
`
}
func (p *PostgreSQLPlatform) InitializeDoctrineTypeMappings() {
	p.DoctrineTypeMapping = map[string]enums.TypesType{
		`bigint`:           enums.TypeBigint,
		`bigserial`:        enums.TypeBigint,
		`bool`:             enums.TypeBoolean,
		`boolean`:          enums.TypeBoolean,
		`bpchar`:           enums.TypeString,
		`bytea`:            enums.TypeBlob,
		`char`:             enums.TypeString,
		`date`:             enums.TypeDateMutable,
		`datetime`:         enums.TypeDatetimeMutable,
		`decimal`:          enums.TypeDecimal,
		`double`:           enums.TypeFloat,
		`double precision`: enums.TypeFloat,
		`float`:            enums.TypeFloat,
		`float4`:           enums.TypeSmallfloat,
		`float8`:           enums.TypeFloat,
		`inet`:             enums.TypeString,
		`int`:              enums.TypeInteger,
		`int2`:             enums.TypeSmallint,
		`int4`:             enums.TypeInteger,
		`int8`:             enums.TypeBigint,
		`integer`:          enums.TypeInteger,
		`interval`:         enums.TypeString,
		`json`:             enums.TypeJson,
		`jsonb`:            enums.TypeJson,
		`money`:            enums.TypeDecimal,
		`numeric`:          enums.TypeDecimal,
		`serial`:           enums.TypeInteger,
		`serial4`:          enums.TypeInteger,
		`serial8`:          enums.TypeBigint,
		`real`:             enums.TypeSmallfloat,
		`smallint`:         enums.TypeSmallint,
		`text`:             enums.TypeText,
		`time`:             enums.TypeTimeMutable,
		`timestamp`:        enums.TypeDatetimeMutable,
		`timestamptz`:      enums.TypeDatetimetzMutable,
		`timetz`:           enums.TypeTimeMutable,
		`tsvector`:         enums.TypeText,
		`uuid`:             enums.TypeGuid,
		`varchar`:          enums.TypeString,
		`year`:             enums.TypeDateMutable,
		`_varchar`:         enums.TypeString,
	}
}
func (p *PostgreSQLPlatform) CreateReservedKeywordsList() keywords.KeywordListInterface {
	return keywords.NewPostgreSQLKeywords()
}
func (p *PostgreSQLPlatform) GetBlobTypeDeclarationSQL(column map[string]any) string {
	return `BYTEA`
}

func (p *PostgreSQLPlatform) GetDefaultValueDeclarationSQL(column map[string]any) string {
	if v, ok := column[`autoincrement`]; ok && v == true {
		return ``
	}

	return p.AbstractPlatform.GetDefaultValueDeclarationSQL(column)
}

func (p *PostgreSQLPlatform) SupportsColumnCollation() bool {
	return true
}
func (p *PostgreSQLPlatform) GetJsonTypeDeclarationSQL(column map[string]any) string {
	if _, ok := column[`jsonb`]; ok {
		return `JSONB`
	}

	return `JSON`
}
func (p *PostgreSQLPlatform) CreateSchemaManager(connection *platforms.Connection) schema_managers.AbstractSchemaManagerInterface {
	return postgres_schema_manager.NewPostgreSQLSchemaManager(connection, p)
}

func (p *PostgreSQLPlatform) GetIndexDeclarationSQL(index *assets.Index) string {
	columns := index.GetColumns()

	if len(columns) == 0 {
		panic(`Incomplete definition. "columns" required.`)
	}

	return `CONSTRAINT ` + index.GetQuotedName(p) + ` ` + p.GetCreateIndexSQLFlags(index) + `(` + strings.Join(
		index.GetQuotedColumns(p),
		`, `,
	) + `)` + p.GetPartialIndexSQL(index)
}
