package dtos

type QueryBuilderInterface interface {
	ToString() string
}

type QueryBuilderOrString struct {
	QueryBuilder QueryBuilderInterface
	String       *string
}

func (u *QueryBuilderOrString) ToString() string {
	if u.String != nil {
		return *u.String
	}
	return u.QueryBuilder.ToString()
}
