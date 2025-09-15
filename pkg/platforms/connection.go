package platforms

import (
	"context"
	databaseSql "database/sql"
	gdh "github.com/KoNekoD/gormite/pkg/gormite_databases_helpers"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

type Connection struct {
	db       gdh.Database
	platform AbstractPlatformInterface
}

func NewConnection(db gdh.Database, p AbstractPlatformInterface) *Connection {
	return &Connection{db: db, platform: p}
}

func Fetch[T any](c *Connection, sql string, typedData T) T {
	err := c.db.Select(sql).Scan(&typedData).Exec(context.Background())
	if err != nil && !errors.Is(err, databaseSql.ErrNoRows) {
		log.Warn("error when fetching", "err", err)
	}

	return typedData
}

func (c *Connection) GetDatabasePlatform() AbstractPlatformInterface {
	return c.platform
}

func (c *Connection) FetchAllAssociative(sql string) []map[string]any {
	result := make([]map[string]any, 0)

	err := c.db.Get(sql).ScanCol(&result).Exec(context.Background())
	if err != nil {
		log.Warn("error when fetching", "err", err)
	}

	return result
}

func (c *Connection) GetDatabase() string {
	database := ""

	sql := c.platform.GetDummySelectSQL(c.platform.GetCurrentDatabaseExpression())

	err := c.db.Get(sql).ScanCol(&database).Exec(context.Background())
	if err != nil {
		panic(err)
	}

	return database
}
