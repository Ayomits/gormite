package dtos

type SelectTableNamesDto struct {
	TableName  string `db:"table_name"`
	SchemaName string `db:"schema_name"`
}

func (s *SelectTableNamesDto) GetSchemaName() string {
	return s.SchemaName
}

func (s *SelectTableNamesDto) GetTableName() string {
	return s.TableName
}
