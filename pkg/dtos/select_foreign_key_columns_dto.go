package dtos

type SelectForeignKeyColumnsDto struct {
	TableName  string `db:"table_name"`
	SchemaName string `db:"schema_name"`
	Conname    string `db:"conname"`
	Condef     string `db:"condef"`
}

func (s *SelectForeignKeyColumnsDto) GetSchemaName() string {
	return s.SchemaName
}

func (s *SelectForeignKeyColumnsDto) GetTableName() string {
	return s.TableName
}
