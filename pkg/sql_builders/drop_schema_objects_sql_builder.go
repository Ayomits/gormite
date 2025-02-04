package sql_builders

import (
	"github.com/KoNekoD/gormite/pkg/assets"
)

type DropSchemaObjectsSQLBuilder struct {
	platform SqlBuildersPlatform
}

func NewDropSchemaObjectsSQLBuilder(platform SqlBuildersPlatform) *DropSchemaObjectsSQLBuilder {
	return &DropSchemaObjectsSQLBuilder{platform: platform}
}

func (b *DropSchemaObjectsSQLBuilder) BuildSQL(schema *assets.Schema) []string {
	results := make([]string, 0)

	results = append(
		results,
		b.buildSequenceStatements(schema.GetSequences())...,
	)
	results = append(results, b.buildTableStatements(schema.GetTables())...)

	return results
}

func (b *DropSchemaObjectsSQLBuilder) buildTableStatements(tables []*assets.Table) []string {
	return b.platform.GetDropTablesSQL(tables)
}

func (b *DropSchemaObjectsSQLBuilder) buildSequenceStatements(sequences map[string]*assets.Sequence) []string {
	statements := make([]string, 0, len(sequences))

	for _, sequence := range sequences {
		statements = append(
			statements,
			b.platform.GetDropSequenceSQL(sequence.GetQuotedName(b.platform)),
		)
	}

	return statements
}
