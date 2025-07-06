package assets

import (
	"github.com/KoNekoD/gormite/pkg/types"
)

type ColumnOption func(c *Column)

func WithColumnLength(length *int) ColumnOption {
	return func(c *Column) {
		c.length = length
	}
}

func WithColumnNotNull() ColumnOption {
	return func(c *Column) {
		c.notnull = true
	}
}

func WithColumnDefault(columnDefault string) ColumnOption {
	return func(c *Column) {
		c.columnDefault = &columnDefault
	}
}

func WithColumnPrecision(precision *int) ColumnOption {
	return func(c *Column) {
		c.precision = precision
	}
}

func WithColumnScale(scale *int) ColumnOption {
	return func(c *Column) {
		c.scale = scale
	}
}

func WithColumnFixed() ColumnOption {
	return func(c *Column) {
		c.fixed = true
	}
}

func WithColumnAutoIncrement() ColumnOption {
	return func(c *Column) {
		c.autoincrement = true
	}
}

func WithColumnComment(comment string) ColumnOption {
	return func(c *Column) {
		c.comment = comment
	}
}

// Column - Object representation of a database column.
type Column struct {
	*AbstractAsset

	columnType       types.AbstractTypeInterface
	length           *int
	precision        *int
	scale            *int
	unsigned         bool
	fixed            bool
	notnull          bool
	columnDefault    *string
	autoincrement    bool
	platformOptions  map[string]any
	columnDefinition *string
	comment          string
}

// NewColumn - Creates a new Column.
func NewColumn(
	name string,
	Type types.AbstractTypeInterface,
	options ...ColumnOption,
) *Column {
	c := &Column{
		AbstractAsset:   NewAbstractAsset(),
		platformOptions: make(map[string]any),
	}

	c.SetName(name)
	c.SetType(Type)
	c.SetOptions(options)

	return c
}

func (c *Column) SetOptions(options []ColumnOption) *Column {
	for _, option := range options {
		option(c)
	}

	return c
}

func (c *Column) SetType(Type types.AbstractTypeInterface) *Column {
	c.columnType = Type

	return c
}

func (c *Column) SetLength(length int) *Column {
	c.length = &length
	return c
}

func (c *Column) SetPrecision(precision int) *Column {
	c.precision = &precision
	return c
}

func (c *Column) SetScale(scale int) *Column {
	c.scale = &scale
	return c
}

func (c *Column) SetUnsigned() *Column {
	c.unsigned = true
	return c
}

func (c *Column) SetFixed() *Column {
	c.fixed = true
	return c
}

func (c *Column) SetNotNull() *Column {
	c.notnull = true
	return c
}

func (c *Column) SetColumnDefault(columnDefault string) *Column {
	c.columnDefault = &columnDefault
	return c
}

func (c *Column) SetPlatformOptions(platformOptions map[string]interface{}) *Column {
	c.platformOptions = platformOptions
	return c
}

func (c *Column) SetPlatformOption(name string, value interface{}) *Column {
	c.platformOptions[name] = value
	return c
}

func (c *Column) SetColumnDefinition(columnDefinition string) *Column {
	c.columnDefinition = &columnDefinition
	return c
}

func (c *Column) GetColumnType() types.AbstractTypeInterface {
	return c.columnType
}

func (c *Column) GetLength() *int {
	return c.length
}

func (c *Column) GetPrecision() *int {
	return c.precision
}

func (c *Column) GetScale() *int {
	return c.scale
}

func (c *Column) GetUnsigned() bool {
	return c.unsigned
}

func (c *Column) GetFixed() bool {
	return c.fixed
}

func (c *Column) GetNotNull() bool {
	return c.notnull
}

func (c *Column) GetColumnDefault() *string {
	return c.columnDefault
}

func (c *Column) GetPlatformOptions() map[string]interface{} {
	return c.platformOptions
}

func (c *Column) HasPlatformOption(name string) bool {
	_, ok := c.platformOptions[name]
	return ok
}

func (c *Column) GetPlatformOption(name string) interface{} {
	return c.platformOptions[name]
}

func (c *Column) GetColumnDefinition() *string {
	return c.columnDefinition
}

func (c *Column) GetAutoincrement() bool {
	return c.autoincrement
}

func (c *Column) SetAutoincrement() *Column {
	c.autoincrement = true
	return c
}

func (c *Column) SetComment(comment string) *Column {
	c.comment = comment
	return c
}

func (c *Column) GetComment() string {
	return c.comment
}

func (c *Column) ToArray() map[string]interface{} {
	columnDefault := ""
	if c.columnDefault != nil {
		columnDefault = *c.columnDefault
	}

	data := map[string]interface{}{
		"name":             c.name,
		"type":             c.columnType,
		"default":          columnDefault,
		"notnull":          c.notnull,
		"length":           c.length,
		"precision":        c.precision,
		"scale":            c.scale,
		"fixed":            c.fixed,
		"unsigned":         c.unsigned,
		"autoincrement":    c.autoincrement,
		"columnDefinition": c.columnDefinition,
		"comment":          c.comment,
	}

	for k, v := range c.platformOptions {
		data[k] = v
	}

	return data
}

func (c *Column) Clone() *Column {
	cloned := *c

	return &cloned
}
