package dtos

// UnionQuery - This struct should be instantiated only by QueryBuilder.
type UnionQuery struct {
	UnionParts []*Union
	OrderBy    []string
	Limit      *Limit
}

func NewUnionQuery(unionParts []*Union, orderBy []string, limit *Limit) *UnionQuery {
	return &UnionQuery{
		UnionParts: unionParts,
		OrderBy:    orderBy,
		Limit:      limit,
	}
}

func (u *UnionQuery) GetUnionParts() []*Union {
	return u.UnionParts
}

func (u *UnionQuery) GetOrderBy() []string {
	return u.OrderBy
}

func (u *UnionQuery) GetLimit() *Limit {
	return u.Limit
}
