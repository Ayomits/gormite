package assets

import (
	"golang.org/x/exp/maps"
	"strings"
)

// UniqueConstraint - Class for a unique constraint.
type UniqueConstraint struct {
	*AbstractAsset

	// columns - Asset identifier instances of the column names the unique constraint is associated with.
	columns map[string]*Identifier

	// flags - Platform specific flags
	flags map[string]bool

	options map[string]string
}

type UniqueConstraintOption func(u *UniqueConstraint)

func WithFlags(flags []string) UniqueConstraintOption {
	return func(u *UniqueConstraint) {
		for _, flag := range flags {
			u.AddFlag(flag)
		}
	}
}

func WithOptions(options map[string]string) UniqueConstraintOption {
	return func(u *UniqueConstraint) {
		u.options = options
	}
}

func NewUniqueConstraint(name string, columns []string, options ...UniqueConstraintOption) *UniqueConstraint {
	v := &UniqueConstraint{
		AbstractAsset: NewAbstractAsset(),
		columns:       make(map[string]*Identifier),
		flags:         make(map[string]bool),
		options:       make(map[string]string),
	}
	for _, option := range options {
		option(v)
	}

	v.SetName(name)

	for _, column := range columns {
		v.addColumn(column)
	}

	return v
}

// GetColumns - Returns the names of the referencing table columns the constraint is associated with.
func (u *UniqueConstraint) GetColumns() []string {
	return maps.Keys(u.columns)
}

// GetQuotedColumns - Returns the quoted representation of the column names the constraint is associated with.
// But only if they were defined with one or a column name
// is a keyword reserved by the platform.
// Otherwise, the plain unquoted value as inserted is returned.
func (u *UniqueConstraint) GetQuotedColumns(platform AssetsPlatform) []string {
	columns := make([]string, 0, len(u.columns))

	for _, column := range u.columns {
		columns = append(columns, column.GetQuotedName(platform))
	}

	return columns
}

func (u *UniqueConstraint) GetUnquotedColumns() []string {
	columns := make([]string, 0, len(u.columns))

	for _, column := range u.GetColumns() {
		columns = append(columns, u.trimQuotes(column))
	}

	return columns
}

// GetFlags - Returns platform specific flags for unique constraint.
func (u *UniqueConstraint) GetFlags() []string {
	return maps.Keys(u.flags)
}

// AddFlag - Adds flag for a unique constraint that translates to platform specific handling.
func (u *UniqueConstraint) AddFlag(flag string) *UniqueConstraint {
	u.flags[strings.ToLower(flag)] = true

	return u
}

// HasFlag - Does this unique constraint have a specific flag?
func (u *UniqueConstraint) HasFlag(flag string) bool {
	_, ok := u.flags[strings.ToLower(flag)]

	return ok
}

// RemoveFlag - Removes a flag.
func (u *UniqueConstraint) RemoveFlag(flag string) {
	delete(u.flags, strings.ToLower(flag))
}

func (u *UniqueConstraint) HasOption(name string) bool {
	_, ok := u.options[strings.ToLower(name)]

	return ok
}

func (u *UniqueConstraint) GetOption(name string) string {
	v, ok := u.options[strings.ToLower(name)]

	if !ok {
		panic("option " + name + " not found")
	}

	return v
}

func (u *UniqueConstraint) GetOptions() map[string]string {
	return u.options
}

func (u *UniqueConstraint) addColumn(column string) {
	u.columns[column] = NewIdentifier(column)
}
