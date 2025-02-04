package dtos

type GetPortableTableDefinitionInputDto interface {
	GetSchemaName() string
	GetTableName() string
}
