package sql_builders

import (
	"github.com/KoNekoD/gormite/pkg/gormite/dtos"
	"github.com/KoNekoD/gormite/pkg/gormite/enums"
	"strings"
)

type DefaultUnionSQLBuilder struct {
	platform SqlBuildersPlatform
}

func NewDefaultUnionSQLBuilder(platform SqlBuildersPlatform) *DefaultUnionSQLBuilder {
	return &DefaultUnionSQLBuilder{platform: platform}
}

func (builder *DefaultUnionSQLBuilder) BuildSQL(unionQuery *dtos.UnionQuery) (
	string,
	error,
) {
	var parts []string
	for _, union := range unionQuery.GetUnionParts() {
		if union.UnionType != nil {
			if *union.UnionType == enums.UnionTypeAll {
				parts = append(parts, builder.platform.GetUnionAllSQL())
			} else {
				parts = append(parts, builder.platform.GetUnionDistinctSQL())
			}
		}

		parts = append(
			parts,
			builder.platform.GetUnionSelectPartSQL(union.Query.ToString()),
		)
	}

	orderBy := unionQuery.GetOrderBy()

	if len(orderBy) > 0 {
		parts = append(parts, "ORDER BY "+strings.Join(orderBy, ", "))
	}

	sql := strings.Join(parts, " ")

	limit := unionQuery.GetLimit()

	if limit.IsDefined() {
		modifiedSQL, err := builder.platform.ModifyLimitQuery(
			sql,
			limit.MaxResults,
			limit.FirstResult,
		)
		if err != nil {
			return "", err
		}
		sql = modifiedSQL
	}

	return sql, nil
}
