package sql_builders

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/dtos"
	"github.com/KoNekoD/gormite/pkg/enums"
	"github.com/pkg/errors"
	"strings"
)

type NotSupportedError struct {
	Feature string
}

func (e *NotSupportedError) Error() string {
	return fmt.Sprintf("%s is not supported", e.Feature)
}

type DefaultSelectSQLBuilder struct {
	platform      SqlBuildersPlatform
	forUpdateSQL  *string
	skipLockedSQL *string
}

func NewDefaultSelectSQLBuilder(
	platform SqlBuildersPlatform,
	forUpdateSQL, skipLockedSQL string,
) *DefaultSelectSQLBuilder {
	return &DefaultSelectSQLBuilder{
		platform:      platform,
		forUpdateSQL:  &forUpdateSQL,
		skipLockedSQL: &skipLockedSQL,
	}
}

func (b *DefaultSelectSQLBuilder) BuildSQL(selectQuery *dtos.SelectQuery) (
	string,
	error,
) {
	var parts []string
	parts = append(parts, "SELECT")

	if selectQuery.Distinct {
		parts = append(parts, "DISTINCT")
	}

	parts = append(parts, strings.Join(selectQuery.Columns, ", "))

	from := selectQuery.GetFrom()

	if len(from) > 0 {
		parts = append(parts, "FROM "+strings.Join(from, ", "))
	}

	where := selectQuery.GetWhere()

	if where != nil {
		parts = append(parts, "WHERE "+*where)
	}

	groupBy := selectQuery.GetGroupBy()

	if len(groupBy) > 0 {
		parts = append(parts, "GROUP BY "+strings.Join(groupBy, ", "))
	}

	having := selectQuery.GetHaving()

	if having != nil {
		parts = append(parts, "HAVING "+*having)
	}

	orderBy := selectQuery.GetOrderBy()

	if len(orderBy) > 0 {
		parts = append(parts, "ORDER BY "+strings.Join(orderBy, ", "))
	}

	sql := strings.Join(parts, " ")

	limit := selectQuery.GetLimit()

	if limit.IsDefined() {
		modifiedSQL, err := b.platform.ModifyLimitQuery(
			sql,
			limit.MaxResults,
			limit.FirstResult,
		)
		if err != nil {
			return "", errors.Wrap(err, "failed to modify limit query")
		}
		sql = modifiedSQL
	}

	forUpdate := selectQuery.GetForUpdate()

	if forUpdate != nil {
		if b.forUpdateSQL == nil {
			return "", &NotSupportedError{Feature: "FOR UPDATE"}
		}

		sql += " " + *b.forUpdateSQL

		if forUpdate.GetConflictResolutionMode() == enums.ConflictResolutionModeSkipLocked {
			if b.skipLockedSQL == nil {
				return "", &NotSupportedError{Feature: "SKIP LOCKED"}
			}

			sql += " " + *b.skipLockedSQL
		}
	}

	return sql, nil
}
