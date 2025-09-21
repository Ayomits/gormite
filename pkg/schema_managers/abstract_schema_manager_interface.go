package schema_managers

import (
	"github.com/KoNekoD/gormite/pkg/assets"
	"github.com/KoNekoD/gormite/pkg/dtos"
)

type AbstractSchemaManagerInterface interface {
	IntrospectSchema() *assets.Schema
	ListTables() []*assets.Table
	ListSequences() []*assets.Sequence
	CreateSchemaConfig() *dtos.SchemaConfig
	ListSchemaNames() []string

	GetPortableDatabaseDefinition(row map[string]any) string
	GetPortableSequenceDefinition(sequence *dtos.ListSequencesDto) *assets.Sequence
	GetPortableTableColumnDefinition(tableColumn *dtos.SelectTableColumnsDto) *assets.Column
	GetPortableViewDefinition(view map[string]any) *assets.View
	GetPortableTableForeignKeyDefinition(tableForeignKey *dtos.SelectForeignKeyColumnsDto) *assets.ForeignKeyConstraint

	GetPortableTableDefinition(table dtos.GetPortableTableDefinitionInputDto) string

	SelectTableNames(databaseName string) []*dtos.SelectTableNamesDto
	SelectTableColumns(
		databaseName string,
		tableName *string,
	) []*dtos.SelectTableColumnsDto
	SelectIndexColumns(
		databaseName string,
		tableName *string,
	) []*dtos.SelectIndexColumnsDto
	SelectForeignKeyColumns(
		databaseName string,
		tableName *string,
	) []*dtos.SelectForeignKeyColumnsDto

	FetchTableOptionsByTable(
		databaseName string,
		tableName *string,
	) map[string]*dtos.FetchTableOptionsByTableDto

	GetPortableTableIndexesList(
		tableIndexes []*dtos.SelectIndexColumnsDto,
		tableName string,
	) map[string]*assets.Index
}
