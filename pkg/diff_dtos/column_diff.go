package diff_dtos

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/assets"
	"strings"
)

// ColumnDiff - Represents the change of a column.
type ColumnDiff struct {
	oldColumn *assets.Column
	newColumn *assets.Column
}

func NewColumnDiff(
	oldColumn *assets.Column,
	newColumn *assets.Column,
) *ColumnDiff {
	return &ColumnDiff{oldColumn: oldColumn, newColumn: newColumn}
}

func (c *ColumnDiff) CountChangedProperties() int {
	sumOfBool := func(b ...bool) int {
		var count int
		for _, b := range b {
			if b {
				count++
			}
		}
		return count
	}

	return sumOfBool(
		c.HasUnsignedChanged(),
		c.HasAutoIncrementChanged(),
		c.HasDefaultChanged(),
		c.HasFixedChanged(),
		c.HasPrecisionChanged(),
		c.HasScaleChanged(),
		c.HasLengthChanged(),
		c.HasNotNullChanged(),
		c.HasNameChanged(),
		c.HasTypeChanged(),
		c.HasCommentChanged(),
	)
}

func (c *ColumnDiff) GetOldColumn() *assets.Column {
	return c.oldColumn
}

func (c *ColumnDiff) GetNewColumn() *assets.Column {
	return c.newColumn
}

func (c *ColumnDiff) HasNameChanged() bool {
	// Column names are case-insensitive
	return !strings.EqualFold(c.oldColumn.GetName(), c.newColumn.GetName())
}

func (c *ColumnDiff) HasTypeChanged() bool {
	oldType := fmt.Sprintf("%T", c.oldColumn.GetColumnType())
	newType := fmt.Sprintf("%T", c.newColumn.GetColumnType())

	return oldType != newType
}

func (c *ColumnDiff) HasDefaultChanged() bool {
	oldDefault := c.oldColumn.GetColumnDefault()
	newDefault := c.newColumn.GetColumnDefault()

	// Null values need to be checked additionally as they tell whether to create or drop a default value.
	// null != 0, null != false, null != '' etc. This affects platform's table alteration SQL generation.
	if (newDefault == nil) != (oldDefault == nil) {
		return true
	}

	// Remaining cases: oldDefault,newDefault all non-nil, or oldDefault,newDefault both nil
	if oldDefault == nil || newDefault == nil {
		return oldDefault != newDefault
	}

	return *oldDefault != *newDefault
}

func (c *ColumnDiff) HasLengthChanged() bool {
	oldLength := c.oldColumn.GetLength()
	newLength := c.newColumn.GetLength()

	if oldLength == nil || newLength == nil {
		return oldLength != newLength
	}

	return *oldLength != *newLength
}

func (c *ColumnDiff) HasUnsignedChanged() bool {
	return c.oldColumn.GetUnsigned() != c.newColumn.GetUnsigned()
}

func (c *ColumnDiff) HasAutoIncrementChanged() bool {
	return c.oldColumn.GetAutoincrement() != c.newColumn.GetAutoincrement()
}

func (c *ColumnDiff) HasFixedChanged() bool {
	return c.oldColumn.GetFixed() != c.newColumn.GetFixed()
}

func (c *ColumnDiff) HasPrecisionChanged() bool {
	return c.oldColumn.GetPrecision() != c.newColumn.GetPrecision()
}

func (c *ColumnDiff) HasScaleChanged() bool {
	return c.oldColumn.GetScale() != c.newColumn.GetScale()
}

func (c *ColumnDiff) HasNotNullChanged() bool {
	return c.oldColumn.GetNotNull() != c.newColumn.GetNotNull()
}

func (c *ColumnDiff) HasCommentChanged() bool {
	return c.oldColumn.GetComment() != c.newColumn.GetComment()
}
