package docs_example

import (
	"context"
	"github.com/KoNekoD/go-snaps/snaps"
	"github.com/KoNekoD/gormite/pkg/diff_calc"
	"github.com/KoNekoD/gormite/pkg/gormite_databases"
	gdh "github.com/KoNekoD/gormite/pkg/gormite_databases_helpers"
	"github.com/KoNekoD/gormite/pkg/local_schema"
	"github.com/KoNekoD/gormite/pkg/platforms"
	"github.com/KoNekoD/gormite/pkg/platforms/postgres_platform"
	"github.com/KoNekoD/gormite/pkg/schema_managers/postgres_schema_manager"
	_ "github.com/KoNekoD/gormite/tests/docs_example/resources"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"testing"
)

type mockQuery struct{}

func (m mockQuery) Scan(dest ...interface{}) gdh.QueryInterface { return m }

func (m mockQuery) ScanCol(dest ...interface{}) gdh.QueryInterface { return m }

func (m mockQuery) Exec(ctx context.Context) error { return nil }

type mockPostgresDatabase struct{}

func (m mockPostgresDatabase) Exec(ctx context.Context, sql string, arguments ...any) (gdh.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m mockPostgresDatabase) Select(sql string, args ...interface{}) gdh.QueryInterface {
	return &mockQuery{}
}

func (m mockPostgresDatabase) Get(sql string, args ...interface{}) gdh.QueryInterface {
	return &mockQuery{}
}

func (m mockPostgresDatabase) Query(ctx context.Context, sql string, args ...any) (gdh.Rows, error) {
	panic("not implemented")
}
func (m mockPostgresDatabase) GetNamedArgs(args any) any {
	panic("not implemented")
}

func TestDiffRunnerTest(t *testing.T) {
	d := gormite_databases.PostgresDatabase{}

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
