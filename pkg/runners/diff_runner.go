package runners

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/diff_calc"
	"github.com/KoNekoD/gormite/pkg/gormite_databases"
	"github.com/KoNekoD/gormite/pkg/local_schema"
	"github.com/KoNekoD/gormite/pkg/platforms"
	"github.com/KoNekoD/gormite/pkg/platforms/postgres_platform"
	"github.com/KoNekoD/gormite/pkg/schema_managers/postgres_schema_manager"
	"github.com/pkg/errors"
	"os"
	"time"
)

const GooseMigrationTemplate = `-- +goose Up
-- +goose StatementBegin
%s
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
%s
-- +goose StatementEnd
`

type MigrationToolType string

const (
	MigrationToolTypeGoose   MigrationToolType = "goose"
	MigrationToolTypeMigrate MigrationToolType = "migrate"
)

type DiffRunnerOptions struct {
	Tool string
	Dsn  string
}

type DiffRunner struct{ opts DiffRunnerOptions }

func NewDiffRunner(opts DiffRunnerOptions) *DiffRunner {
	return &DiffRunner{opts: opts}
}

func (r *DiffRunner) Run() error {
	d := gormite_databases.NewPostgresDatabase(r.opts.Dsn)

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
		return errors.Wrap(err, "failed to introspect local schema")
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
		return errors.New("No changes detected")
	}

	up := manager.AlterSchema(diff)
	down := manager.AlterSchema(diffDown)

	switch r.opts.Tool {
	case string(MigrationToolTypeMigrate):
		key := time.Now().Format("20060102150405")
		file := fmt.Sprintf("migrations/%s_gen.up.sql", key)
		fileDown := fmt.Sprintf("migrations/%s_gen.down.sql", key)
		if err := os.WriteFile(file, []byte(up), 0644); err != nil {
			return errors.Wrap(err, "Cannot write up sql file")
		}
		if err := os.WriteFile(fileDown, []byte(down), 0644); err != nil {
			return errors.Wrap(err, "Cannot write down sql file")
		}
	case string(MigrationToolTypeGoose):
		key := time.Now().Format("20060102150405")
		file := fmt.Sprintf("migrations/%s_gen.sql", key)
		if err := os.WriteFile(
			file,
			[]byte(fmt.Sprintf(GooseMigrationTemplate, up, down)),
			0644,
		); err != nil {
			return errors.Wrap(err, "Cannot write migration file")
		}
	}

	return nil
}
