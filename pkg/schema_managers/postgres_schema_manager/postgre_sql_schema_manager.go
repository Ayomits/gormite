package postgres_schema_manager

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/assets"
	"github.com/KoNekoD/gormite/pkg/dtos"
	"github.com/KoNekoD/gormite/pkg/platforms"
	"github.com/KoNekoD/gormite/pkg/schema_managers/abstract_schema_managers"
	"github.com/KoNekoD/gormite/pkg/types"
	"github.com/KoNekoD/gormite/pkg/utils"
	"github.com/KoNekoD/ptrs/pkg/ptrs"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type PostgreSQLSchemaManager struct {
	*abstract_schema_managers.AbstractSchemaManager
	currentSchema *string
}

func NewPostgreSQLSchemaManager(
	connection *platforms.Connection,
	platform platforms.AbstractPlatformInterface,
) *PostgreSQLSchemaManager {
	v := &PostgreSQLSchemaManager{
		AbstractSchemaManager: abstract_schema_managers.NewAbstractSchemaManager(
			connection,
			platform,
		),
	}

	v.Child = v

	return v
}

type ListSchemaNamesDto struct {
	SchemaName string `db:"schema_name"`
}

func (m *PostgreSQLSchemaManager) ListSchemaNames() []string {
	typedData := make([]ListSchemaNamesDto, 0)

	return utils.MapSlice(
		platforms.Fetch(
			m.Connection,
			`
SELECT schema_name
FROM   information_schema.schemata
WHERE  schema_name NOT LIKE 'pg\_%'
AND    schema_name != 'information_schema'
`,
			typedData,
		), func(t ListSchemaNamesDto) string {
			return t.SchemaName
		},
	)
}

func (m *PostgreSQLSchemaManager) CreateSchemaConfig() *dtos.SchemaConfig {
	config := m.AbstractSchemaManager.CreateSchemaConfig()

	config.SetName(m.getCurrentSchema())

	return config
}

func (m *PostgreSQLSchemaManager) getCurrentSchema() *string {
	if m.currentSchema == nil {
		m.currentSchema = m.determineCurrentSchema()
	}

	return m.currentSchema
}

type determineCurrentSchemaDto struct {
	SchemaName string `db:"schema_name"`
}

func (m *PostgreSQLSchemaManager) determineCurrentSchema() *string {
	dto := platforms.Fetch(
		m.Connection,
		`SELECT current_schema() AS schema_name`,
		determineCurrentSchemaDto{},
	)

	if dto.SchemaName == "" {
		return nil
	}

	return &dto.SchemaName
}

func (m *PostgreSQLSchemaManager) GetPortableTableForeignKeyDefinition(tableForeignKey *dtos.SelectForeignKeyColumnsDto) *assets.ForeignKeyConstraint {
	var onUpdate, onDelete *string

	onUpdateRegex := regexp.MustCompile(`ON UPDATE ([a-zA-Z0-9]+( (NULL|ACTION|DEFAULT))?)`)
	onDeleteRegex := regexp.MustCompile(`ON DELETE ([a-zA-Z0-9]+( (NULL|ACTION|DEFAULT))?)`)

	// Check for ON UPDATE condition
	if match := onUpdateRegex.FindStringSubmatch(tableForeignKey.Condef); len(match) > 0 {
		onUpdate = &match[1]
	}

	// Check for ON DELETE condition
	if match := onDeleteRegex.FindStringSubmatch(tableForeignKey.Condef); len(match) > 0 {
		onDelete = &match[1]
	}

	foreignKeyRegex := regexp.MustCompile(`FOREIGN KEY \((.+)\) REFERENCES (.+)\((.+)\)`)

	// Parse the FOREIGN KEY constraint
	if match := foreignKeyRegex.FindStringSubmatch(tableForeignKey.Condef); len(match) == 4 {
		localColumns := strings.Split(match[1], ",")
		for i := range localColumns {
			localColumns[i] = strings.TrimSpace(localColumns[i])
		}

		foreignColumns := strings.Split(match[3], ",")
		for i := range foreignColumns {
			foreignColumns[i] = strings.TrimSpace(foreignColumns[i])
		}

		foreignTable := match[2]

		return assets.NewForeignKeyConstraint(
			tableForeignKey.Conname,
			localColumns,
			foreignTable,
			foreignColumns,
			map[string]interface{}{
				"onUpdate": onUpdate,
				"onDelete": onDelete,
			},
		)
	}

	panic(fmt.Errorf("invalid foreign key definition"))
}

func (m *PostgreSQLSchemaManager) GetPortableDatabaseDefinition(row map[string]interface{}) string {
	return row["datname"].(string)
}

func (m *PostgreSQLSchemaManager) GetPortableSequenceDefinition(sequence *dtos.ListSequencesDto) *assets.Sequence {
	sequenceName := ""
	if sequence.Schemaname != "public" {
		sequenceName = sequence.Schemaname + "." + sequence.Relname
	} else {
		sequenceName = sequence.Relname
	}

	return assets.NewSequence(
		sequenceName,
		assets.WithAllocationSize(sequence.GetIncrementBy()),
		assets.WithInitialValue(sequence.GetMinValue()),
	)
}

func (m *PostgreSQLSchemaManager) GetPortableTableColumnDefinition(tableColumn *dtos.SelectTableColumnsDto) *assets.Column {
	var length *int

	if slices.Contains([]string{"varchar", "bpchar"}, tableColumn.Type) {
		matches := regexp.MustCompile(`\((\d*)\)`).FindStringSubmatch(tableColumn.CompleteType)
		if len(matches) == 2 {
			lenInt, _ := strconv.Atoi(matches[1])
			length = ptrs.AsPtr(lenInt)
		}
	}

	autoincrement := tableColumn.Attidentity == 'd'

	matches := make([]string, 0)

	_ = tableColumn.Default
	_ = tableColumn.CompleteType

	if tableColumn.Default != nil {
		if matches = regexp.MustCompile(`^['(](.*)[')]::`).FindStringSubmatch(*tableColumn.Default); len(matches) == 2 {
			tableColumn.Default = &matches[1]
		} else if matches = regexp.MustCompile(`^NULL::`).FindStringSubmatch(*tableColumn.Default); len(matches) == 2 {
			tableColumn.Default = nil
		}
	}

	//if length != nil &&  *length == -1 && nil != tableColumn.atttypmod {
	//	length = tableColumn.atttypmod
	//}

	if length != nil && *length <= 0 {
		length = nil
	}

	fixed := false

	var precision *int
	var scale *int
	var jsonb *bool

	dbType := strings.ToLower(tableColumn.Type)

	if tableColumn.DomainType != nil && *tableColumn.DomainType != "" && !m.Platform.HasDoctrineTypeMappingFor(dbType) {
		dbType = strings.ToLower(*tableColumn.DomainType)
		tableColumn.CompleteType = *tableColumn.DomainCompleteType
	}

	typeMapping := m.Platform.GetDoctrineTypeMapping(dbType)

	switch dbType {
	case "smallint", "int2", "int", "int4", "integer", "bigint", "int8":
		length = nil

	case "bool", "boolean":
		if tableColumn.Default != nil && *tableColumn.Default == "true" {
			tableColumn.Default = ptrs.AsPtr(`true`)
		}
		if tableColumn.Default != nil && *tableColumn.Default == "false" {
			tableColumn.Default = ptrs.AsPtr(`false`)
		}
		length = nil
	case "json", "text", "_varchar", "varchar":
		tableColumn.Default = m.parseDefaultExpression(tableColumn.Default)
	case "char", "bpchar":
		fixed = true
	case "float", "float4", "float8", "double", "double precision", "real", "decimal", "money", "numeric":
		re := regexp.MustCompile(`([A-Za-z]+\(([0-9]+),([0-9]+)\))`)
		matches := re.FindStringSubmatch(tableColumn.CompleteType)

		if len(matches) > 2 {
			precisionInt, _ := strconv.Atoi(matches[2])
			precision = ptrs.AsPtr(precisionInt)
			scaleInt, _ := strconv.Atoi(matches[3])
			scale = ptrs.AsPtr(scaleInt)
			length = nil
		}
	case "year":
		length = nil
	case "jsonb":
		jsonb = ptrs.AsPtr(true)
	}

	if tableColumn.Default != nil {
		re := regexp.MustCompile(`'([^']+)'::`)
		if matches = re.FindStringSubmatch(*tableColumn.Default); len(matches) == 2 {
			tableColumn.Default = ptrs.AsPtr(matches[1])
		}
	}

	options := make([]assets.ColumnOption, 0)
	options = append(options, assets.WithColumnLength(length))

	if tableColumn.IsNotnull {
		options = append(options, assets.WithColumnNotNull())
	}

	if tableColumn.Default != nil {
		options = append(
			options,
			assets.WithColumnDefault(*tableColumn.Default),
		)
	}

	options = append(options, assets.WithColumnPrecision(precision))

	options = append(options, assets.WithColumnScale(scale))

	if fixed {
		options = append(options, assets.WithColumnFixed())
	}

	if autoincrement {
		options = append(options, assets.WithColumnAutoIncrement())
	}

	if tableColumn.Comment != nil {
		options = append(
			options,
			assets.WithColumnComment(*tableColumn.Comment),
		)
	}

	column := assets.NewColumn(
		tableColumn.Field,
		types.GetType(typeMapping),
		options...,
	)

	if tableColumn.Collation != nil {
		column.SetPlatformOption("collation", *tableColumn.Collation)
	}

	if _, ok := column.GetColumnType().(*types.JsonType); ok {
		column.SetPlatformOption("jsonb", jsonb)
	}

	return column
}

func (m *PostgreSQLSchemaManager) GetPortableViewDefinition(view map[string]interface{}) *assets.View {
	return assets.NewView(
		view["schemaname"].(string)+"."+view["viewname"].(string),
		view["definition"].(string),
	)
}

func (m *PostgreSQLSchemaManager) GetPortableTableIndexesList(
	tableIndexes []*dtos.SelectIndexColumnsDto,
	tableName string,
) map[string]*assets.Index {
	buffer := make([]*dtos.PortableTableIndexesDto, 0)

	for _, row := range tableIndexes {
		colNumbers := strings.Split(row.Indkey, " ")

		columnNameSql := fmt.Sprintf(
			"SELECT attnum, attname FROM pg_attribute WHERE attrelid=%s AND attnum IN (%s) ORDER BY attnum ASC",
			*row.Indrelid,
			strings.Join(colNumbers, " ,"),
		)

		indexColumns := platforms.Fetch(
			m.Connection,
			columnNameSql,
			make([]dtos.GetColNameDto, 0),
		)
		for _, colNum := range colNumbers {
			for _, colRow := range indexColumns {
				if colNum != colRow.Attnum {
					continue
				}

				buffer = append(
					buffer, &dtos.PortableTableIndexesDto{
						KeyName:    row.RelName,
						ColumnName: strings.TrimSpace(colRow.Attname),
						NonUnique:  !row.IndisUnique,
						Primary:    row.IndisPrimary,
						Where:      row.Where,
					},
				)
			}
		}
	}

	return m.AbstractSchemaManager.GetPortableTableIndexesList(
		buffer,
		tableName,
	)
}

func (m *PostgreSQLSchemaManager) GetPortableTableDefinition(table dtos.GetPortableTableDefinitionInputDto) string {
	currentSchema := m.getCurrentSchema()

	if table.GetSchemaName() == *currentSchema {
		return table.GetTableName()
	}

	return table.GetSchemaName() + "." + table.GetTableName()
}

func (m *PostgreSQLSchemaManager) SelectTableColumns(
	databaseName string,
	tableName *string,
) []*dtos.SelectTableColumnsDto {
	sql := "SELECT "

	if tableName == nil {
		sql += "c.relname AS table_name, n.nspname AS schema_name,"
	}

	sql += fmt.Sprintf(
		`
	           a.attnum,
	           quote_ident(a.attname) AS field,
	           t.typname AS type,
	           format_type(a.atttypid, a.atttypmod) AS complete_type,
	           (SELECT tc.collcollate FROM pg_catalog.pg_collation tc WHERE tc.oid = a.attcollation) AS collation,
	           (SELECT t1.typname FROM pg_catalog.pg_type t1 WHERE t1.oid = t.typbasetype) AS domain_type,
	           (SELECT format_type(t2.typbasetype, t2.typtypmod) FROM
	             pg_catalog.pg_type t2 WHERE t2.typtype = 'd' AND t2.oid = a.atttypid) AS domain_complete_type,
	           a.attnotnull AS isnotnull,
	           a.attidentity,
	           (SELECT 't'
	            FROM pg_index
	            WHERE c.oid = pg_index.indrelid
	               AND pg_index.indkey[0] = a.attnum
	               AND pg_index.indisprimary = 't'
	           ) AS pri,
	           (%s) AS default,
	           (SELECT pg_description.description
	               FROM pg_description WHERE pg_description.objoid = c.oid AND a.attnum = pg_description.objsubid
	           ) AS comment
	           FROM pg_attribute a
	               INNER JOIN pg_class c
	                   ON c.oid = a.attrelid
	               INNER JOIN pg_type t
	                   ON t.oid = a.atttypid
	               INNER JOIN pg_namespace n
	                   ON n.oid = c.relnamespace
	               LEFT JOIN pg_depend d
	                   ON d.objid = c.oid
	                       AND d.deptype = 'e'
	                       AND d.classid = (SELECT oid FROM pg_class WHERE relname = 'pg_class')
	           `, m.Platform.GetDefaultColumnValueSQLSnippet(),
	)

	conditions := make([]string, 0)
	conditions = append(conditions, "a.attnum > 0")
	conditions = append(conditions, "c.relkind = 'r'")
	conditions = append(conditions, "d.refobjid IS NULL")
	conditions = append(conditions, m.buildQueryConditions(tableName)...)

	sql += " WHERE " + strings.Join(conditions, " AND ") + " ORDER BY a.attnum"

	typedData := make([]dtos.SelectTableColumnsDto, 0)

	return utils.MapSlice(
		platforms.Fetch(m.Connection, sql, typedData),
		ptrs.AsPtr,
	)
}

func (m *PostgreSQLSchemaManager) SelectIndexColumns(
	databaseName string,
	tableName *string,
) []*dtos.SelectIndexColumnsDto {
	sql := "SELECT"

	if tableName == nil {
		sql += " tc.relname AS table_name, tn.nspname AS schema_name,"
	}

	sql += `
	quote_ident(ic.relname) AS relname,
		i.indisunique,
		i.indisprimary,
		i.indkey,
		i.indrelid,
		pg_get_expr(indpred, indrelid) AS "where"
	FROM pg_index i
	JOIN pg_class AS tc ON tc.oid = i.indrelid
	JOIN pg_namespace tn ON tn.oid = tc.relnamespace
	JOIN pg_class AS ic ON ic.oid = i.indexrelid
	WHERE ic.oid IN (
		SELECT indexrelid
	FROM pg_index i, pg_class c, pg_namespace n
	`

	conditions := make([]string, 0)
	conditions = append(conditions, "c.oid = i.indrelid")
	conditions = append(conditions, "c.relnamespace = n.oid")
	conditions = append(conditions, m.buildQueryConditions(tableName)...)

	sql += " WHERE " + strings.Join(conditions, " AND ") + ")"

	return utils.MapSlice(
		platforms.Fetch(
			m.Connection,
			sql,
			make([]dtos.SelectIndexColumnsDto, 0),
		), ptrs.AsPtr,
	)
}

func (m *PostgreSQLSchemaManager) SelectTableNames(databaseName string) []*dtos.SelectTableNamesDto {
	sql := `
SELECT quote_ident(table_name) AS table_name,
	table_schema AS schema_name
FROM information_schema.tables
WHERE table_catalog = '` + databaseName + `'
AND table_schema NOT LIKE 'pg\_%'
AND table_schema != 'information_schema'
AND table_name != 'geometry_columns'
AND table_name != 'spatial_ref_sys'
AND table_type = 'BASE TABLE'
	`

	return utils.MapSlice(
		platforms.Fetch(
			m.Connection,
			sql,
			make([]dtos.SelectTableNamesDto, 0),
		), ptrs.AsPtr,
	)
}

func (m *PostgreSQLSchemaManager) SelectForeignKeyColumns(
	databaseName string,
	tableName *string,
) []*dtos.SelectForeignKeyColumnsDto {
	sql := "SELECT"

	if tableName == nil {
		sql += " tc.relname AS table_name, tn.nspname AS schema_name,"
	}

	sql += `
	quote_ident(r.conname) as conname,
		pg_get_constraintdef(r.oid, true) as condef
	FROM pg_constraint r
	JOIN pg_class AS tc ON tc.oid = r.conrelid
	JOIN pg_namespace tn ON tn.oid = tc.relnamespace
	WHERE r.conrelid IN
	(
		SELECT c.oid
	FROM pg_class c, pg_namespace n
	`

	conditions := make([]string, 0)
	conditions = append(conditions, "n.oid = c.relnamespace")
	conditions = append(conditions, m.buildQueryConditions(tableName)...)

	sql += " WHERE " + strings.Join(
		conditions,
		" AND ",
	) + ") AND r.contype = 'f'"

	return utils.MapSlice(
		platforms.Fetch(
			m.Connection,
			sql,
			make([]dtos.SelectForeignKeyColumnsDto, 0),
		), ptrs.AsPtr,
	)
}

func (m *PostgreSQLSchemaManager) FetchTableOptionsByTable(
	databaseName string,
	tableName *string,
) map[string]*dtos.FetchTableOptionsByTableDto {
	sql := `
	SELECT c.relname,
		CASE c.relpersistence WHEN 'u' THEN true ELSE false END as unlogged,
		obj_description(c.oid, 'pg_class') AS comment
	FROM pg_class c
	INNER JOIN pg_namespace n
	ON n.oid = c.relnamespace
	`

	conditions := make([]string, 0)
	conditions = append(conditions, "c.relkind = 'r'")
	conditions = append(conditions, m.buildQueryConditions(tableName)...)

	sql += " WHERE " + strings.Join(conditions, " AND ")

	result := make(map[string]*dtos.FetchTableOptionsByTableDto)

	utils.MapSlice(
		platforms.Fetch(
			m.Connection,
			sql,
			make([]dtos.FetchTableOptionsByTableDto, 0),
		),
		func(t dtos.FetchTableOptionsByTableDto) *dtos.FetchTableOptionsByTableDto {
			if result[t.Relname] != nil {
				panic(fmt.Sprintf("duplicate table name: %s", t.Relname))
			}
			result[t.Relname] = &t
			return &t
		},
	)

	return result
}

func (m *PostgreSQLSchemaManager) buildQueryConditions(tableName *string) []string {
	conditions := make([]string, 0)

	if tableName != nil {
		tableNameStr := *tableName
		if strings.Contains(*tableName, ".") {
			parts := strings.Split(*tableName, ".")
			schemaName := parts[0]
			tableNameStr = parts[1]
			conditions = append(
				conditions,
				"n.nspname = "+m.Platform.QuoteStringLiteral(schemaName),
			)
		} else {
			conditions = append(
				conditions,
				"n.nspname = ANY(current_schemas(false))",
			)
		}

		identifier := assets.NewIdentifier(tableNameStr)
		conditions = append(
			conditions,
			"c.relname = "+m.Platform.QuoteStringLiteral(identifier.GetName()),
		)
	}

	conditions = append(
		conditions,
		"n.nspname NOT IN ('pg_catalog', 'information_schema', 'pg_toast')",
	)

	return conditions
}

func (m *PostgreSQLSchemaManager) parseDefaultExpression(defaultExpression *string) *string {
	if defaultExpression == nil {
		return nil
	}

	return ptrs.AsPtr(strings.ReplaceAll(*defaultExpression, "''", "'"))
}
