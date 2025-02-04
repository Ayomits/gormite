package docs_example

import (
	"context"
	"github.com/KoNekoD/go-snaps/snaps"
	"github.com/KoNekoD/gormite/pkg/diff_calc"
	"github.com/KoNekoD/gormite/pkg/gormite_databases"
	"github.com/KoNekoD/gormite/pkg/local_schema"
	"github.com/KoNekoD/gormite/pkg/platforms"
	"github.com/KoNekoD/gormite/pkg/platforms/postgres_platform"
	"github.com/KoNekoD/gormite/pkg/schema_managers/postgres_schema_manager"
	_ "github.com/KoNekoD/gormite/test/docs_example/resources"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"testing"
)

type mockQuery struct{}

func (m mockQuery) Scan(dest ...interface{}) gormite_databases.QueryInterface { return m }

func (m mockQuery) ScanCol(dest ...interface{}) gormite_databases.QueryInterface { return m }

func (m mockQuery) Exec(ctx context.Context) error { return nil }

type mockPostgresDatabase struct{}

func (m mockPostgresDatabase) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m mockPostgresDatabase) Select(sql string, args ...interface{}) gormite_databases.QueryInterface {
	return &mockQuery{}
}

func (m mockPostgresDatabase) Get(sql string, args ...interface{}) gormite_databases.QueryInterface {
	return &mockQuery{}
}

func TestDiffRunnerTest(t *testing.T) {
	d := &mockPostgresDatabase{}

	p := postgres_platform.NewPostgreSQLPlatform()

	manager := postgres_schema_manager.NewPostgreSQLSchemaManager(
		platforms.NewConnection(
			d,
			postgres_platform.NewPostgreSQLPlatform(),
		), p,
	)

	oldSchema := manager.IntrospectSchema()

	newSchema, err := local_schema.IntrospectLocalSchema("resources/gormite.yaml")
	if err != nil {
		panic(errors.Wrap(err, "failed to introspect local schema"))
	}

	// Remove implicit indexes
	for _, table := range oldSchema.GetTables() {
		table.ClearImplicitIndexes()
	}
	for _, table := range newSchema.GetTables() {
		table.ClearImplicitIndexes()
	}

	c := diff_calc.NewComparator(p)

	diff := c.CompareSchemas(oldSchema, newSchema)
	diffDown := c.CompareSchemas(newSchema, oldSchema)

	if diff.IsEmpty() {
		panic(errors.New("No changes detected"))
	}

	// TODO: Test snapshots for sql
	//up := manager.AlterSchema(diff)
	//down := manager.AlterSchema(diffDown)

	snaps.MatchStandaloneSnapshot(t, diff)
	snaps.MatchStandaloneSnapshot(t, diffDown)
}
