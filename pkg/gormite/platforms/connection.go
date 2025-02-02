package platforms

import (
	"context"
	"database/sql"
	"github.com/KoNekoD/gormite/pkg/gormite/gormite_databases"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

type Connection struct {
	Platform AbstractPlatformInterface
	pg       gormite_databases.PostgresDatabaseInterface
}

func NewConnection(db any, platform AbstractPlatformInterface) *Connection {
	switch v := db.(type) {
	case gormite_databases.PostgresDatabaseInterface:
		return &Connection{Platform: platform, pg: v}
	}

	panic("unknown type")
}

func Fetch[T any](c *Connection, sqlString string, typedData T) T {
	err := c.pg.Select(sqlString).Scan(&typedData).Exec(context.Background())
	if err != nil && !errors.Is(err, sql.ErrNoRows) { // TODO: ПОПРАВИТЬ!!!
		log.Warn("error when fetching", "err", err)
	}

	return typedData
}

func (c *Connection) GetDatabasePlatform() AbstractPlatformInterface {
	return c.Platform
}

func (c *Connection) FetchAllAssociative(sql string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	err := c.pg.Get(sql).ScanCol(&result).Exec(context.Background())
	if err != nil {
		log.Warn("error when fetching", "err", err)
	}

	return result
}

func (c *Connection) GetDatabase() string {
	database := ""

	selectSQL := c.Platform.GetDummySelectSQL(c.Platform.GetCurrentDatabaseExpression())

	_ = c.pg.Get(selectSQL).ScanCol(&database).Exec(context.Background())

	return database
}
