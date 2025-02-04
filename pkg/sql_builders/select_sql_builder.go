package sql_builders

import (
	"github.com/KoNekoD/gormite/pkg/dtos"
)

// SelectSQLBuilder - The SQL builder should be instantiated only by database platforms.
type SelectSQLBuilder interface {
	BuildSQL(query *dtos.SelectQuery) (string, error)
}
