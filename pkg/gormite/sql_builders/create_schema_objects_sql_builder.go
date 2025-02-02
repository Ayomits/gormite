package sql_builders

import (
	"github.com/KoNekoD/gormite/pkg/gormite/assets"
)

type CreateSchemaObjectsSQLBuilder struct {
	platform SqlBuildersPlatform
}

func NewCreateSchemaObjectsSQLBuilder(platform SqlBuildersPlatform) *CreateSchemaObjectsSQLBuilder {
	return &CreateSchemaObjectsSQLBuilder{platform: platform}
}

func (b *CreateSchemaObjectsSQLBuilder) BuildSQL(schema *assets.Schema) []string {
	results := make([]string, 0)

	results = append(
		results,
		b.buildNamespaceStatements(schema.GetNamespaces())...,
	)
	results = append(
		results,
		b.buildSequenceStatements(schema.GetSequences())...,
	)
	results = append(results, b.buildTableStatements(schema.GetTables())...)

	return results
}

func (b *CreateSchemaObjectsSQLBuilder) buildNamespaceStatements(namespaces []string) []string {
	statements := make([]string, 0, len(namespaces))

	if b.platform.SupportsSchemas() {
		for _, namespace := range namespaces {
			statements = append(
				statements,
				b.platform.GetCreateSchemaSQL(namespace),
			)
		}
	}

	return statements
}

func (b *CreateSchemaObjectsSQLBuilder) buildTableStatements(tables []*assets.Table) []string {
	return b.platform.GetCreateTablesSQL(tables)
}

func (b *CreateSchemaObjectsSQLBuilder) buildSequenceStatements(sequences map[string]*assets.Sequence) []string {
	statements := make([]string, 0, len(sequences))

	for _, sequence := range sequences {
		statements = append(
			statements,
			b.platform.GetCreateSequenceSQL(sequence),
		)
	}

	return statements
}
