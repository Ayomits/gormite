package sql_builders

import (
	"github.com/KoNekoD/gormite/pkg/gormite/dtos"
)

type UnionSQLBuilder interface {
	BuildSQL(query *dtos.UnionQuery) (string, error)
}
