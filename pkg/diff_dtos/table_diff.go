package diff_dtos

import (
	"github.com/KoNekoD/gormite/pkg/assets"
)

// TableDiff - Table Diff.
type TableDiff struct {
	oldTable            *assets.Table
	droppedForeignKeys  []*assets.ForeignKeyConstraint
	addedColumns        map[string]*assets.Column
	changedColumns      map[string]*ColumnDiff
	droppedColumns      map[string]*assets.Column
	addedIndexes        map[string]*assets.Index
	modifiedIndexes     []*assets.Index
	droppedIndexes      map[string]*assets.Index
	renamedIndexes      map[string]*assets.Index
	addedForeignKeys    []*assets.ForeignKeyConstraint
	modifiedForeignKeys []*assets.ForeignKeyConstraint
}

func NewTableDiff(
	oldTable *assets.Table,
	droppedForeignKeys []*assets.ForeignKeyConstraint,
	addedColumns map[string]*assets.Column,
	changedColumns map[string]*ColumnDiff,
	droppedColumns map[string]*assets.Column,
	addedIndexes map[string]*assets.Index,
	modifiedIndexes []*assets.Index,
	droppedIndexes map[string]*assets.Index,
	renamedIndexes map[string]*assets.Index,
	addedForeignKeys []*assets.ForeignKeyConstraint,
	modifiedForeignKeys []*assets.ForeignKeyConstraint,
) *TableDiff {
	return &TableDiff{
		oldTable:            oldTable,
		droppedForeignKeys:  droppedForeignKeys,
		addedColumns:        addedColumns,
		changedColumns:      changedColumns,
		droppedColumns:      droppedColumns,
		addedIndexes:        addedIndexes,
		modifiedIndexes:     modifiedIndexes,
		droppedIndexes:      droppedIndexes,
		renamedIndexes:      renamedIndexes,
		addedForeignKeys:    addedForeignKeys,
		modifiedForeignKeys: modifiedForeignKeys,
	}
}

// GetOldTable - Returns the old table.
func (d *TableDiff) GetOldTable() *assets.Table {
	return d.oldTable
}

func (d *TableDiff) GetAddedColumns() map[string]*assets.Column {
	return d.addedColumns
}

func (d *TableDiff) GetChangedColumns() map[string]*ColumnDiff {
	return d.changedColumns
}

func (d *TableDiff) GetDroppedColumns() map[string]*assets.Column {
	return d.droppedColumns
}

func (d *TableDiff) GetAddedIndexes() map[string]*assets.Index {
	return d.addedIndexes
}

func (d *TableDiff) GetModifiedIndexes() []*assets.Index {
	return d.modifiedIndexes
}

func (d *TableDiff) GetDroppedIndexes() map[string]*assets.Index {
	return d.droppedIndexes
}

func (d *TableDiff) GetRenamedIndexes() map[string]*assets.Index {
	return d.renamedIndexes
}

func (d *TableDiff) GetAddedForeignKeys() []*assets.ForeignKeyConstraint {
	return d.addedForeignKeys
}

func (d *TableDiff) GetModifiedForeignKeys() []*assets.ForeignKeyConstraint {
	return d.modifiedForeignKeys
}

func (d *TableDiff) GetDroppedForeignKeys() []*assets.ForeignKeyConstraint {
	return d.droppedForeignKeys
}

// IsEmpty - Returns whether the diff is empty (contains no changes).
func (d *TableDiff) IsEmpty() bool {
	return len(d.addedColumns) == 0 &&
		len(d.changedColumns) == 0 &&
		len(d.droppedColumns) == 0 &&
		len(d.addedIndexes) == 0 &&
		len(d.modifiedIndexes) == 0 &&
		len(d.droppedIndexes) == 0 &&
		len(d.renamedIndexes) == 0 &&
		len(d.addedForeignKeys) == 0 &&
		len(d.modifiedForeignKeys) == 0 &&
		len(d.droppedForeignKeys) == 0
}
