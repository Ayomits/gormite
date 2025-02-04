package dtos

type Join struct {
	joinType  string
	table     string
	alias     string
	condition *string
}

func NewJoin(joinType string, table string, alias string, condition *string) *Join {
	return &Join{joinType: joinType, table: table, alias: alias, condition: condition}
}

func (j *Join) GetType() string {
	return j.joinType
}

func (j *Join) GetTable() string {
	return j.table
}

func (j *Join) GetAlias() string {
	return j.alias
}

func (j *Join) GetCondition() *string {
	return j.condition
}

func NewInnerJoin(table string, alias string, condition *string) *Join {
	return NewJoin("INNER", table, alias, condition)
}

func NewLeftJoin(table string, alias string, condition *string) *Join {
	return NewJoin("LEFT", table, alias, condition)
}

func NewRightJoin(table string, alias string, condition *string) *Join {
	return NewJoin("RIGHT", table, alias, condition)
}
