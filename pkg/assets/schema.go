package assets

import (
	"github.com/KoNekoD/gormite/pkg/dtos"
	"golang.org/x/exp/maps"
	"strings"
)

// Schema - Object representation of a database schema.
// Different vendors have very inconsistent naming with regard to the concept
// of a "schema". Doctrine understands a schema as the entity that conceptually
// wraps a set of database objects such as tables, sequences, indexes and
// foreign keys that belong to each other into a namespace. A Doctrine Schema
// has nothing to do with the "SCHEMA" defined as in PostgreSQL, it is more
// related to the concept of "DATABASE" that exists in MySQL and PostgreSQL.
// Every asset in the doctrine schema has a name. A name consists of either a
// namespace.local name pair or just a local unqualified name.
// The abstraction layer that covers a PostgreSQL schema is the namespace of an
// database object (asset). A schema can have a name, which will be used as
// default namespace for the unqualified database objects that are created in
// the schema.
// In the case of MySQL where cross-database queries are allowed this leads to
// databases being "misinterpreted" as namespaces. This is intentional, however
// the CREATE/DROP SQL visitors will just filter this queries and do not
// execute them. Only the queries for the currently connected database are
// executed.
type Schema struct {
	*AbstractAsset

	// namespaces - The namespaces in this schema.
	namespaces map[string]string

	tables map[string]*Table

	sequences map[string]*Sequence

	schemaConfig *dtos.SchemaConfig
}

func NewSchema(
	tables []*Table,
	sequences []*Sequence,
	schemaConfig *dtos.SchemaConfig,
	namespaces []string,
) *Schema {
	v := &Schema{
		AbstractAsset: NewAbstractAsset(),
		namespaces:    make(map[string]string),
		tables:        make(map[string]*Table),
		sequences:     make(map[string]*Sequence),
	}

	if schemaConfig == nil {
		schemaConfig = dtos.NewSchemaConfig()
	}

	v.schemaConfig = schemaConfig

	name := schemaConfig.GetName()

	if name != nil {
		v.SetName(*name)
	}

	for _, namespace := range namespaces {
		v.CreateNamespace(namespace)
	}

	for _, table := range tables {
		v.addTable(table)
	}

	for _, sequence := range sequences {
		v.addSequence(sequence)
	}

	return v
}

func (s *Schema) addTable(table *Table) {
	namespaceName := table.GetNamespaceName()
	tableName := s.normalizeName(table)

	if _, ok := s.tables[tableName]; ok {
		panic("table already exists " + tableName)
	}

	if namespaceName != "" && !table.IsInDefaultNamespace(s.GetName()) && !s.HasNamespace(namespaceName) {
		s.CreateNamespace(namespaceName)
	}

	s.tables[tableName] = table
	table.SetSchemaConfig(s.schemaConfig)
}

func (s *Schema) addSequence(sequence *Sequence) {
	namespaceName := sequence.GetNamespaceName()
	seqName := s.normalizeName(sequence)

	if _, ok := s.tables[seqName]; ok {
		panic("sequence already exists " + seqName)
	}

	if namespaceName != "" && !sequence.IsInDefaultNamespace(s.GetName()) && !s.HasNamespace(namespaceName) {
		s.CreateNamespace(namespaceName)
	}

	s.sequences[seqName] = sequence
}

// GetNamespaces - Returns the namespaces of this schema.
func (s *Schema) GetNamespaces() []string {
	return maps.Values(s.namespaces)
}

// GetTables - Gets all tables of this schema.
func (s *Schema) GetTables() []*Table {
	return maps.Values(s.tables)
}

func (s *Schema) GetTable(name string) *Table {
	name = s.getFullQualifiedAssetName(name)

	table, ok := s.tables[name]
	if !ok {
		panic("table " + name + " not found")
	}

	return table
}

func (s *Schema) getFullQualifiedAssetName(name string) string {
	name = s.getUnquotedAssetName(name)

	if !strings.Contains(name, ".") {
		name = s.GetName() + "." + name
	}

	return strings.ToLower(name)
}

// normalizeName - The normalized name is qualified and lower-cased. Lower-casing is
// actually wrong, but we have to do it to keep our sanity. If you are
// using database objects that only differentiate in the casing (FOO vs
// Foo) then you will NOT be able to use Doctrine Schema abstraction.
// Every non-namespaced element is prefixed with this schema name.
func (s *Schema) normalizeName(asset AbstractAssetInterface) string {
	name := asset.GetName()

	if asset.GetNamespaceName() == "" {
		name = s.GetName() + "." + name
	}

	return strings.ToLower(name)
}

// getUnquotedAssetName - Returns the unquoted representation of a given asset name.
func (s *Schema) getUnquotedAssetName(name string) string {
	if s.isIdentifierQuoted(name) {
		name = s.trimQuotes(name)
	}
	return name
}

// HasNamespace - Does this schema have a namespace with the given name?
func (s *Schema) HasNamespace(name string) bool {
	name = strings.ToLower(s.getUnquotedAssetName(name))

	_, ok := s.namespaces[name]

	return ok
}

func (s *Schema) HasTable(name string) bool {
	name = s.getFullQualifiedAssetName(name)
	_, ok := s.tables[name]
	return ok
}

func (s *Schema) HasSequence(name string) bool {
	name = s.getFullQualifiedAssetName(name)

	_, ok := s.sequences[name]

	return ok
}

func (s *Schema) GetSequence(name string) *Sequence {
	name = s.getFullQualifiedAssetName(name)

	v, ok := s.sequences[name]

	if !ok {
		panic("sequence " + name + " not found")
	}

	return v
}

func (s *Schema) GetSequences() map[string]*Sequence {
	return s.sequences
}

// CreateNamespace - Creates a new namespace.
func (s *Schema) CreateNamespace(name string) *Schema {
	unquotedName := strings.ToLower(s.getUnquotedAssetName(name))

	if _, ok := s.namespaces[unquotedName]; ok {
		panic("namespace already exists " + unquotedName)
	}

	s.namespaces[unquotedName] = name

	return s
}

// CreateTable - Creates a new table.
func (s *Schema) CreateTable(name string) *Table {
	table := NewTable(
		name,
		make([]*Column, 0),
		make([]*Index, 0),
		make([]*UniqueConstraint, 0),
		make([]*ForeignKeyConstraint, 0),
		make(map[string]interface{}),
	)

	s.addTable(table)

	for option, value := range s.schemaConfig.GetDefaultTableOptions() {
		table.AddOption(option, value)
	}

	return table
}

// RenameTable - Renames a table.
func (s *Schema) RenameTable(oldName string, newName string) *Schema {
	table := s.GetTable(oldName)
	table.SetName(newName)

	s.dropTable(oldName)
	s.addTable(table)

	return s
}

// dropTable - Drops a table from the schema.
func (s *Schema) dropTable(name string) *Schema {
	name = s.getFullQualifiedAssetName(name)
	s.GetTable(name)
	delete(s.tables, name)

	return s
}

// CreateSequence - Creates a new sequence.
// allocationSize and initialValue is 1 by default
func (s *Schema) CreateSequence(
	name string,
	allocationSize int,
	initialValue int,
) *Sequence {
	seq := NewSequence(
		name,
		WithAllocationSize(allocationSize),
		WithInitialValue(initialValue),
	)
	s.addSequence(seq)

	return seq
}

func (s *Schema) dropSequence(name string) *Schema {
	name = s.getFullQualifiedAssetName(name)
	delete(s.sequences, name)

	return s
}

// Clone - Cloning a Schema triggers a deep clone of all related assets.
func (s *Schema) Clone() {
	for k, table := range s.tables {
		s.tables[k] = table.Clone()
	}

	for k, sequence := range s.sequences {
		s.sequences[k] = sequence.Clone()
	}
}
