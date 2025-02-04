package dtos

type SelectTableColumnsDto struct {
	TableName          string  `db:"table_name"`
	SchemaName         string  `db:"schema_name"`
	Attnum             int     `db:"attnum"`
	Field              string  `db:"field"`
	Type               string  `db:"type"`
	CompleteType       string  `db:"complete_type"`
	Collation          *string `db:"collation"`
	DomainType         *string `db:"domain_type"`
	DomainCompleteType *string `db:"domain_complete_type"`
	IsNotnull          bool    `db:"isnotnull"`
	Attidentity        rune    `db:"attidentity"`
	Pri                *string `db:"pri"`
	Default            *string `db:"default"`
	Comment            *string `db:"comment"`
}

func (s *SelectTableColumnsDto) GetSchemaName() string {
	return s.SchemaName
}

func (s *SelectTableColumnsDto) GetTableName() string {
	return s.TableName
}
