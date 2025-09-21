package sql_builders

import (
	"github.com/KoNekoD/gormite/pkg/assets"
	"github.com/KoNekoD/gormite/pkg/supports_platforms_contracts"
)

type SqlBuildersPlatform interface {
	assets.AssetsPlatform
	supports_platforms_contracts.SupportsPlatform
	ModifyLimitQuery(query string, limit *int, offset int) (string, error)

	GetListViewsSQL(database string) string
	GetCreateSchemaSQL(schemaName string) string
	GetCreateTablesSQL(tables []*assets.Table) []string
	GetCreateSequenceSQL(sequence *assets.Sequence) string
	GetDropTablesSQL(tables []*assets.Table) []string
	GetDropSequenceSQL(name string) string
	GetUnionAllSQL() string
	GetUnionDistinctSQL() string
	GetUnionSelectPartSQL(subQuery string) string
	GetListDatabasesSQL() string
	GetListSequencesSQL(database string) string
	GetDummySelectSQL(expression string) string
	GetDropTableSQL(table string) string
	GetDropTemporaryTableSQL(table string) string
	GetDropIndexSQL(name string, table string) string
	GetDropConstraintSQL(name string, table string) string
	GetDropForeignKeySQL(foreignKey string, table string) string
	GetDropUniqueConstraintSQL(name string, tableName string) string
	GetCreateTableSQL(table *assets.Table) []string
	GetCreateTableInnerSQL(name string, columns []map[string]any, options map[string]any) []string
	GetCreateTableWithoutForeignKeysSQL(table *assets.Table) []string
	BuildCreateTableSQL(table *assets.Table, createForeignKeys bool) []string
	GetCommentOnTableSQL(tableName string, comment string) string
	GetCommentOnColumnSQL(tableName string, columnName string, comment string) string
	GetInlineColumnCommentSQL(comment string) string
	GetCreateTemporaryTableSnippetSQL() string
	GetAlterSequenceSQL(sequence *assets.Sequence) string
	GetCreateIndexSQL(index *assets.Index, table string) string
	GetPartialIndexSQL(index *assets.Index) string
	GetCreatePrimaryKeySQL(index *assets.Index, table string) string
	GetCreateUniqueConstraintSQL(constraint *assets.UniqueConstraint, tableName string) string
	GetDropSchemaSQL(schemaName string) string
	GetCreateForeignKeySQL(foreignKey *assets.ForeignKeyConstraint, table string) string
	GetRenameTableSQL(oldName string, newName string) string
	GetRenameIndexSQL(oldIndexName string, index *assets.Index, tableName string) []string
	GetColumnDeclarationListSQL(columns []map[string]any) string
	GetColumnDeclarationSQL(name string, column map[string]any) string
	GetDefaultValueDeclarationSQL(column map[string]any) string
	GetCheckDeclarationSQL(definition []map[string]any) string
	GetUniqueConstraintDeclarationSQL(constraint *assets.UniqueConstraint) string
	GetIndexDeclarationSQL(index *assets.Index) string
	GetForeignKeyDeclarationSQL(foreignKey *assets.ForeignKeyConstraint) string
	GetAdvancedForeignKeyOptionsSQL(foreignKey *assets.ForeignKeyConstraint) string
	GetForeignKeyReferentialActionSQL(action string) string
	GetForeignKeyBaseDeclarationSQL(foreignKey *assets.ForeignKeyConstraint) string
	GetColumnCharsetDeclarationSQL(charset string) string
	GetColumnCollationDeclarationSQL(collation string) string
	GetCurrentDateSQL() string
	GetCurrentTimeSQL() string
	GetCurrentTimestampSQL() string
	GetCreateViewSQL(name string, sql string) string
	GetDropViewSQL(name string) string
	GetSequenceNextValSQL(sequence string) string
	GetCreateDatabaseSQL(name string) string
	GetDropDatabaseSQL(name string) string
	GetTruncateTableSQL(tableName string, cascade bool) string
}
