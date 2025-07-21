package abstract_schema_managers

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/assets"
	"github.com/KoNekoD/gormite/pkg/diff_calc"
	"github.com/KoNekoD/gormite/pkg/diff_dtos"
	"github.com/KoNekoD/gormite/pkg/dtos"
	"github.com/KoNekoD/gormite/pkg/platforms"
	"github.com/KoNekoD/gormite/pkg/schema_managers"
	"github.com/KoNekoD/ptrs/pkg/ptrs"
	"github.com/KoNekoD/smt/pkg/smt"
	"golang.org/x/exp/maps"
	"slices"
	"strings"
)

type AbstractSchemaManager struct {
	Connection *platforms.Connection

	Platform platforms.AbstractPlatformInterface

	Child schema_managers.AbstractSchemaManagerInterface
}

func NewAbstractSchemaManager(
	connection *platforms.Connection,
	platform platforms.AbstractPlatformInterface,
) *AbstractSchemaManager {
	return &AbstractSchemaManager{Connection: connection, Platform: platform}
}

func (m *AbstractSchemaManager) ListDatabases() []string {
	return smt.MapSlice(
		m.Connection.FetchAllAssociative(m.Platform.GetListDatabasesSQL()),
		func(t map[string]interface{}) string {
			return m.Child.GetPortableDatabaseDefinition(t)
		},
	)
}

func (m *AbstractSchemaManager) ListSchemaNames() []string {
	panic("not implemented")
}

func (m *AbstractSchemaManager) ListSequences() []*assets.Sequence {
	typedData := make([]dtos.ListSequencesDto, 0)

	return smt.MapSlice(
		platforms.Fetch(
			m.Connection,
			m.Platform.GetListSequencesSQL(m.getDatabase()),
			typedData,
		),
		func(t dtos.ListSequencesDto) *assets.Sequence {
			return m.Child.GetPortableSequenceDefinition(&t)
		},
	)
}

func (m *AbstractSchemaManager) ListTableColumns(table string) map[string]*assets.Column {
	database := m.getDatabase()

	return m.GetPortableTableColumnList(
		table,
		database,
		m.Child.SelectTableColumns(
			database,
			ptrs.AsPtr(m.normalizeName(table)),
		),
	)
}

func (m *AbstractSchemaManager) ListTableIndexes(table string) map[string]*assets.Index {
	database := m.getDatabase()
	table = m.normalizeName(table)

	return m.Child.GetPortableTableIndexesList(
		m.Child.SelectIndexColumns(
			database,
			&table,
		), table,
	)
}

func (m *AbstractSchemaManager) TablesExist(names []string) bool {
	names = smt.MapSlice(names, strings.ToLower)

	return len(names) == len(
		smt.SliceIntersect(
			names,
			smt.MapSlice(m.ListTableNames(), strings.ToLower),
		),
	)
}

func (m *AbstractSchemaManager) tableExists(tableName string) bool {
	return m.TablesExist([]string{tableName})
}

func (m *AbstractSchemaManager) ListTableNames() []string {
	return smt.MapSlice(
		m.Child.SelectTableNames(m.getDatabase()),
		func(t *dtos.SelectTableNamesDto) string {
			return m.Child.GetPortableTableDefinition(t)
		},
	)
}

func (m *AbstractSchemaManager) ListTables() []*assets.Table {
	database := m.getDatabase()

	tableColumnsByTable := m.fetchTableColumnsByTable(database)
	indexColumnsByTable := m.fetchIndexColumnsByTable(database)
	foreignKeyColumnsByTable := m.fetchForeignKeyColumnsByTable(database)
	tableOptionsByTable := m.Child.FetchTableOptionsByTable(database, nil)

	tables := make([]*assets.Table, 0)
	for tableName, tableColumns := range tableColumnsByTable {
		options := tableOptionsByTable[tableName]

		tables = append(
			tables, assets.NewTable(
				tableName,
				maps.Values(
					m.GetPortableTableColumnList(
						tableName,
						database,
						tableColumns,
					),
				),
				maps.Values(
					m.Child.GetPortableTableIndexesList(
						indexColumnsByTable[tableName],
						tableName,
					),
				),
				make([]*assets.UniqueConstraint, 0),
				m.getPortableTableForeignKeysList(foreignKeyColumnsByTable[tableName]),
				options.ToArray(),
			),
		)
	}

	return tables
}

func (m *AbstractSchemaManager) normalizeName(name string) string {
	identifier := assets.NewIdentifier(name)

	return identifier.GetName()
}

func (m *AbstractSchemaManager) fetchTableColumnsByTable(databaseName string) map[string][]*dtos.SelectTableColumnsDto {
	data := m.Child.SelectTableColumns(databaseName, nil)

	return fetchAllAssociativeGrouped(m, data)
}

func (m *AbstractSchemaManager) fetchIndexColumnsByTable(databaseName string) map[string][]*dtos.SelectIndexColumnsDto {
	data := m.Child.SelectIndexColumns(databaseName, nil)

	return fetchAllAssociativeGrouped(m, data)
}

func (m *AbstractSchemaManager) fetchForeignKeyColumnsByTable(databaseName string) map[string][]*dtos.SelectForeignKeyColumnsDto {
	data := m.Child.SelectForeignKeyColumns(databaseName, nil)

	return fetchAllAssociativeGrouped(m, data)
}

func (m *AbstractSchemaManager) IntrospectTable(name string) *assets.Table {
	columns := m.ListTableColumns(name)

	if len(columns) == 0 {
		panic("table " + name + " not found")
	}

	return assets.NewTable(
		name,
		maps.Values(columns),
		maps.Values(m.ListTableIndexes(name)),
		make([]*assets.UniqueConstraint, 0),
		m.ListTableForeignKeys(name),
		m.getTableOptions(name).ToArray(),
	)
}

func (m *AbstractSchemaManager) ListViews() []*assets.View {
	return smt.MapSlice(
		m.Connection.FetchAllAssociative(m.Platform.GetListViewsSQL(m.getDatabase())),
		func(t map[string]interface{}) *assets.View {
			return m.Child.GetPortableViewDefinition(t)
		},
	)
}

func (m *AbstractSchemaManager) ListTableForeignKeys(table string) []*assets.ForeignKeyConstraint {
	database := m.getDatabase()

	return m.getPortableTableForeignKeysList(
		m.Child.SelectForeignKeyColumns(
			database, ptrs.AsPtr(m.normalizeName(table)),
		),
	)
}

func (m *AbstractSchemaManager) getTableOptions(table string) *dtos.FetchTableOptionsByTableDto {
	normalizedName := m.normalizeName(table)

	return m.Child.FetchTableOptionsByTable(
		m.getDatabase(),
		&normalizedName,
	)[normalizedName]
}

type getPortableTableIndexesListOptionsSubDto struct {
	lengths []*int
	where   string
}

func (g *getPortableTableIndexesListOptionsSubDto) asMap() map[string]interface{} {
	res := map[string]interface{}{}

	if len(g.lengths) > 0 {
		res["lengths"] = g.lengths
	}

	if g.where != "" {
		res["where"] = g.where
	}

	return res
}

type getPortableTableIndexesListDto struct {
	name    string
	columns []string
	unique  bool
	primary bool
	flags   []string
	options *getPortableTableIndexesListOptionsSubDto
}

func (g *getPortableTableIndexesListDto) addColumn(s string) {
	g.columns = append(g.columns, s)
}

func (g *getPortableTableIndexesListDto) addOptionsLength(length *int) {
	g.options.lengths = append(g.options.lengths, length)
}

func (m *AbstractSchemaManager) GetPortableTableIndexesList(
	tableIndexes []*dtos.PortableTableIndexesDto,
	tableName string,
) map[string]*assets.Index {
	result := make(map[string]*getPortableTableIndexesListDto)

	for _, tableIndex := range tableIndexes {

		indexName := tableIndex.KeyName
		keyName := indexName
		if tableIndex.Primary {
			keyName = "primary"
		}

		keyName = strings.ToLower(keyName)

		if _, ok := result[keyName]; !ok {
			options := &getPortableTableIndexesListOptionsSubDto{
				lengths: make(
					[]*int,
					0,
				),
			}
			if tableIndex.Where != nil {
				options.where = *tableIndex.Where
			}

			result[keyName] = &getPortableTableIndexesListDto{
				name:    indexName,
				columns: make([]string, 0),
				unique:  !tableIndex.NonUnique,
				primary: tableIndex.Primary,
				flags:   make([]string, 0),
				options: options,
			}
		}

		result[keyName].addColumn(tableIndex.ColumnName)
		//result[keyName].addOptionsLength(nil) // tableIndex.Length
	}

	indexes := make(map[string]*assets.Index)

	for indexKey, data := range result {
		indexes[indexKey] = assets.NewIndex(
			data.name,
			data.columns,
			data.unique,
			data.primary,
			data.flags,
			data.options.asMap(),
		)
	}

	return indexes
}

func (m *AbstractSchemaManager) getPortableTableForeignKeysList(tableForeignKeys []*dtos.SelectForeignKeyColumnsDto) []*assets.ForeignKeyConstraint {
	list := make([]*assets.ForeignKeyConstraint, 0)

	for _, value := range tableForeignKeys {
		list = append(list, m.Child.GetPortableTableForeignKeyDefinition(value))
	}

	return list
}

func (m *AbstractSchemaManager) IntrospectSchema() *assets.Schema {
	s := m.Child

	schemaNames := make([]string, 0)

	if m.Platform.SupportsSchemas() {
		schemaNames = s.ListSchemaNames()
	}

	sequences := make([]*assets.Sequence, 0)

	if m.Platform.SupportsSequences() {
		sequences = s.ListSequences()
	}

	tables := s.ListTables()

	// Remove schema_migrations
	for i, table := range tables {
		if slices.Contains(
			[]string{"schema_migrations", "goose_db_version"},
			table.GetName(),
		) {
			tables = slices.Delete(tables, i, i+1)
			break
		}
	}

	return assets.NewSchema(
		tables,
		sequences,
		s.CreateSchemaConfig(),
		schemaNames,
	)
}

func (m *AbstractSchemaManager) CreateSchemaConfig() *dtos.SchemaConfig {
	schemaConfig := dtos.NewSchemaConfig()
	schemaConfig.SetMaxIdentifierLength(m.Platform.GetMaxIdentifierLength())

	return schemaConfig
}

func (m *AbstractSchemaManager) getDatabase() string {
	return m.Connection.GetDatabase()
}

func (m *AbstractSchemaManager) CreateComparator() *diff_calc.Comparator {
	return diff_calc.NewComparator(m.Platform)
}

func fetchAllAssociativeGrouped[T dtos.GetPortableTableDefinitionInputDto](
	m *AbstractSchemaManager,
	typedData []T,
) map[string][]T {
	data := make(map[string][]T)

	for _, row := range typedData {
		tableName := m.Child.GetPortableTableDefinition(row)

		if _, ok := data[tableName]; !ok {
			data[tableName] = make([]T, 0)
		}

		data[tableName] = append(data[tableName], row)
	}

	return data
}

func (m *AbstractSchemaManager) GetPortableTableColumnList(
	table string,
	database string,
	tableColumns []*dtos.SelectTableColumnsDto,
) map[string]*assets.Column {
	list := make(map[string]*assets.Column)

	for _, tableColumn := range tableColumns {
		column := m.Child.GetPortableTableColumnDefinition(tableColumn)

		name := strings.ToLower(column.GetQuotedName(m.Platform))
		list[name] = column
	}

	return list
}

func (m *AbstractSchemaManager) AlterSchema(schemaDiff *diff_dtos.SchemaDiff) string {
	comment := "-- THIS FILE WAS GENERATED BY GORMITE, EDIT IT IF YOU WANT <3"

	mappedRows := make([]string, 0)

	for _, sql := range m.Platform.GetAlterSchemaSQL(schemaDiff) {
		mappedRows = append(mappedRows, fmt.Sprintf("%s;", sql))
	}

	mappedRows = append([]string{comment, ""}, mappedRows...)

	return strings.Join(mappedRows, "\n")
}
