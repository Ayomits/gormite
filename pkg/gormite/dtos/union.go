package dtos

import (
	"github.com/KoNekoD/gormite/pkg/gormite/enums"
)

type Union struct {
	Query     *QueryBuilderOrString
	UnionType *enums.UnionType
}

func NewUnion(query *QueryBuilderOrString) *Union {
	return &Union{Query: query}
}

func NewUnionWithType(
	query *QueryBuilderOrString,
	unionType enums.UnionType,
) *Union {
	return &Union{Query: query, UnionType: &unionType}
}
