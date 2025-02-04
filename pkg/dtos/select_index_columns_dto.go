package dtos

type SelectIndexColumnsDto struct {
	TableName    string  `db:"table_name"`
	SchemaName   string  `db:"schema_name"`
	RelName      string  `db:"relname"`
	IndisUnique  bool    `db:"indisunique"`
	IndisPrimary bool    `db:"indisprimary"`
	Indkey       string  `db:"indkey"`
	Indrelid     *string `db:"indrelid"`
	Where        *string `db:"where"`
}

func (s *SelectIndexColumnsDto) GetSchemaName() string {
	return s.SchemaName
}

func (s *SelectIndexColumnsDto) GetTableName() string {
	return s.TableName
}
