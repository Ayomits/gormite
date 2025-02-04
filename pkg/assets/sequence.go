package assets

import (
	"fmt"
)

// Sequence - Sequence structure.
type Sequence struct {
	*AbstractAsset

	allocationSize int

	initialValue int

	cache *int
}

type SequenceOption func(*Sequence)

func WithAllocationSize(allocationSize int) SequenceOption {
	return func(s *Sequence) {
		s.allocationSize = allocationSize
	}
}

func WithInitialValue(initialValue int) SequenceOption {
	return func(s *Sequence) {
		s.initialValue = initialValue
	}
}

func WithCache(cache int) SequenceOption {
	return func(s *Sequence) {
		s.cache = &cache
	}
}

func NewSequence(name string, options ...SequenceOption) *Sequence {
	v := &Sequence{
		AbstractAsset:  NewAbstractAsset(),
		allocationSize: 1,
		initialValue:   1,
	}

	v.SetName(name)

	for _, option := range options {
		option(v)
	}

	return v
}

func (s *Sequence) GetAllocationSize() int {
	return s.allocationSize
}

func (s *Sequence) GetInitialValue() int {
	return s.initialValue
}

func (s *Sequence) GetCache() *int {
	return s.cache
}

func (s *Sequence) SetAllocationSize(allocationSize int) *Sequence {
	s.allocationSize = allocationSize

	return s
}

func (s *Sequence) SetInitialValue(initialValue int) *Sequence {
	s.initialValue = initialValue

	return s
}

func (s *Sequence) SetCache(cache int) *Sequence {
	s.cache = &cache

	return s
}

// isAutoIncrementsFor - Checks if this sequence is an autoincrement sequence for a given table.
// This is used inside the comparator to not report sequences as missing,
// when the "from" schema implicitly creates the sequences.
func (s *Sequence) IsAutoIncrementsFor(table *Table) bool {
	primaryKey := table.GetPrimaryKey()

	if primaryKey == nil {
		return false
	}

	pkColumns := primaryKey.GetColumns()

	if len(pkColumns) != 1 {
		return false
	}

	column := table.GetColumn(pkColumns[0])

	if !column.GetAutoincrement() {
		return false
	}

	sequenceName := s.GetShortestName(table.GetNamespaceName())
	tableName := table.GetShortestName(table.GetNamespaceName())
	tableSequenceName := fmt.Sprintf("%s__%s__seq", tableName, column.GetShortestName(table.GetNamespaceName()))

	return tableSequenceName == sequenceName
}

func (s *Sequence) Clone() *Sequence {
	cloned := *s

	return &cloned
}
