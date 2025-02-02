package platforms

import (
	"github.com/KoNekoD/gormite/pkg/gormite/assets"
	"github.com/KoNekoD/gormite/pkg/gormite/diff_calc"
	"github.com/KoNekoD/gormite/pkg/gormite/diff_dtos"
	"github.com/KoNekoD/gormite/pkg/gormite/enums"
	"github.com/KoNekoD/gormite/pkg/gormite/keywords"
	"github.com/KoNekoD/gormite/pkg/gormite/schema_managers"
	"github.com/KoNekoD/gormite/pkg/gormite/sql_builders"
	"github.com/KoNekoD/gormite/pkg/gormite/supports_platforms_contracts"
	"github.com/KoNekoD/gormite/pkg/gormite/types"
)

type AbstractPlatformInterface interface {
	types.TypesPlatform
	assets.AssetsPlatform
	supports_platforms_contracts.SupportsPlatform
	sql_builders.SqlBuildersPlatform
	diff_calc.DiffCalcPlatform
	InitializeDoctrineTypeMappings()
	GetLocateExpression(string string, substring string, start *string) string
	GetDateDiffExpression(date1 string, date2 string) string
	GetDateArithmeticIntervalExpression(
		date string,
		operator string,
		interval string,
		unit enums.DateIntervalUnit,
	) string
	GetCurrentDatabaseExpression() string
	CreateReservedKeywordsList() keywords.KeywordListInterface
	CreateSchemaManager(connection *Connection) schema_managers.AbstractSchemaManagerInterface

	// From struct
	CreateSelectSQLBuilder() sql_builders.SelectSQLBuilder
	CreateUnionSQLBuilder() sql_builders.UnionSQLBuilder
	GetMaxIdentifierLength() int
	HasDoctrineTypeMappingFor(dbType string) bool
	GetDoctrineTypeMapping(dbType string) enums.TypesType
	GetDefaultColumnValueSQLSnippet() string
	QuoteStringLiteral(str string) string

	InitializeAllDoctrineTypeMappings()
	ExtractLength(data map[string]interface{}) *int
	RegisterDoctrineTypeMapping(dbType string, doctrineType string)
	GetRegexpExpression() string
	GetLengthExpression(string string) string
	GetModExpression(dividend string, divisor string) string
	GetTrimExpression(str string, mode enums.TrimMode, char *string) string
	GetSubstringExpression(string string, start string, length *string) string
	GetConcatExpression(string ...string) string
	GetDateAddSecondsExpression(date string, seconds string) string
	GetDateSubSecondsExpression(date string, seconds string) string
	GetDateAddMinutesExpression(date string, minutes string) string
	GetDateSubMinutesExpression(date string, minutes string) string
	GetDateAddHourExpression(date string, hours string) string
	GetDateSubHourExpression(date string, hours string) string
	GetDateAddDaysExpression(date string, days string) string
	GetDateSubDaysExpression(date string, days string) string
	GetDateAddWeeksExpression(date string, weeks string) string
	GetDateSubWeeksExpression(date string, weeks string) string
	GetDateAddMonthExpression(date string, months string) string
	GetDateSubMonthExpression(date string, months string) string
	GetDateAddQuartersExpression(date string, quarters string) string
	GetDateSubQuartersExpression(date string, quarters string) string
	GetDateAddYearsExpression(date string, years string) string
	GetDateSubYearsExpression(date string, years string) string
	MultiplyInterval(interval string, multiplier int) string
	GetBitAndComparisonExpression(value1 string, value2 string) string
	GetBitOrComparisonExpression(value1 string, value2 string) string
	GetCreateIndexSQLFlags(index *assets.Index) string
	QuoteSingleIdentifier(str string) string
	GetTemporaryTableName(tableName string) string
	ConvertBooleans(item any) any
	GetDateTimeFormatString() string
	GetDateTimeTzFormatString() string
	GetTimeFormatString() string
	DoModifyLimitQuery(query string, limit *int, offset int) string
	CreateSavePoint(savepoint string) string
	ReleaseSavePoint(savepoint string) string
	RollbackSavePoint(savepoint string) string
	EscapeStringForLike(inputString string, escapeChar string) string
	ColumnToArray(column *assets.Column) map[string]interface{}
	GetLikeWildcardCharacters() string
	GetPreAlterTableIndexForeignKeySQL(diff *diff_dtos.TableDiff) []string
	GetPostAlterTableIndexForeignKeySQL(diff *diff_dtos.TableDiff) []string
	GetAlterTableSQL(diff *diff_dtos.TableDiff) []string
	GetAlterSchemaSQL(diff *diff_dtos.SchemaDiff) []string
}
