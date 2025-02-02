package dtos

// SelectQuery - This class should be instantiated only by QueryBuilder.
type SelectQuery struct {
	Distinct  bool
	Columns   []string
	From      []string
	Where     *string
	GroupBy   []string
	Having    *string
	OrderBy   []string
	Limit     *Limit
	ForUpdate *ForUpdate
}

func NewSelectQuery(
	distinct bool,
	columns []string,
	from []string,
	where *string,
	groupBy []string,
	having *string,
	orderBy []string,
	limit *Limit,
	forUpdate *ForUpdate,
) *SelectQuery {
	return &SelectQuery{
		Distinct:  distinct,
		Columns:   columns,
		From:      from,
		Where:     where,
		GroupBy:   groupBy,
		Having:    having,
		OrderBy:   orderBy,
		Limit:     limit,
		ForUpdate: forUpdate,
	}
}

func (s *SelectQuery) IsDistinct() bool {
	return s.Distinct
}

func (s *SelectQuery) GetColumns() []string {
	return s.Columns
}

func (s *SelectQuery) GetFrom() []string {
	return s.From
}

func (s *SelectQuery) GetWhere() *string {
	return s.Where
}

func (s *SelectQuery) GetGroupBy() []string {
	return s.GroupBy
}

func (s *SelectQuery) GetHaving() *string {
	return s.Having
}

func (s *SelectQuery) GetOrderBy() []string {
	return s.OrderBy
}

func (s *SelectQuery) GetLimit() *Limit {
	return s.Limit
}

func (s *SelectQuery) GetForUpdate() *ForUpdate {
	return s.ForUpdate
}
