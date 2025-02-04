package dtos

type PortableTableIndexesDto struct {
	KeyName    string
	ColumnName string
	NonUnique  bool
	Primary    bool
	Where      *string
}
