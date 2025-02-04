package dtos

type From struct {
	table string
	alias *string
}

func NewFrom(table string, alias *string) *From {
	return &From{table: table, alias: alias}
}

func (f *From) GetTable() string {
	return f.table
}

func (f *From) GetAlias() *string {
	return f.alias
}
