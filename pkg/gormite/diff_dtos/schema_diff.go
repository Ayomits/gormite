package diff_dtos

import (
	"github.com/KoNekoD/gormite/pkg/gormite/assets"
)

// SchemaDiff - Differences between two schemas.
type SchemaDiff struct {
	createdSchemas   []string
	droppedSchemas   []string
	createdTables    []*assets.Table
	alteredTables    []*TableDiff
	droppedTables    []*assets.Table
	createdSequences []*assets.Sequence
	alteredSequences []*assets.Sequence
	droppedSequences []*assets.Sequence
}

func NewSchemaDiff(
	createdSchemas []string,
	droppedSchemas []string,
	createdTables []*assets.Table,
	alteredTables []*TableDiff,
	droppedTables []*assets.Table,
	createdSequences []*assets.Sequence,
	alteredSequences []*assets.Sequence,
	droppedSequences []*assets.Sequence,
) *SchemaDiff {
	return &SchemaDiff{
		createdSchemas:   createdSchemas,
		droppedSchemas:   droppedSchemas,
		createdTables:    createdTables,
		alteredTables:    alteredTables,
		droppedTables:    droppedTables,
		createdSequences: createdSequences,
		alteredSequences: alteredSequences,
		droppedSequences: droppedSequences,
	}
}

func (s *SchemaDiff) GetCreatedSchemas() []string {
	return s.createdSchemas
}

func (s *SchemaDiff) GetDroppedSchemas() []string {
	return s.droppedSchemas
}

func (s *SchemaDiff) GetCreatedTables() []*assets.Table {
	return s.createdTables
}

func (s *SchemaDiff) GetAlteredTables() []*TableDiff {
	return s.alteredTables
}

func (s *SchemaDiff) GetDroppedTables() []*assets.Table {
	return s.droppedTables
}

func (s *SchemaDiff) GetCreatedSequences() []*assets.Sequence {
	return s.createdSequences
}

func (s *SchemaDiff) GetAlteredSequences() []*assets.Sequence {
	return s.alteredSequences
}

func (s *SchemaDiff) GetDroppedSequences() []*assets.Sequence {
	return s.droppedSequences
}

// IsEmpty - Returns whether the diff is empty (contains no changes).
func (s *SchemaDiff) IsEmpty() bool {
	return len(s.createdSchemas) == 0 &&
		len(s.droppedSchemas) == 0 &&
		len(s.createdTables) == 0 &&
		len(s.alteredTables) == 0 &&
		len(s.droppedTables) == 0 &&
		len(s.createdSequences) == 0 &&
		len(s.alteredSequences) == 0 &&
		len(s.droppedSequences) == 0
}
