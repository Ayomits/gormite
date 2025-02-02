package dtos

// SchemaConfig - Configuration for a Schema.
type SchemaConfig struct {
	maxIdentifierLength int
	name                *string
	defaultTableOptions map[string]interface{}
}

func NewSchemaConfig() *SchemaConfig {
	return &SchemaConfig{
		maxIdentifierLength: 63,
		name:                nil,
		defaultTableOptions: make(map[string]interface{}),
	}
}

func (s *SchemaConfig) SetMaxIdentifierLength(maxIdentifierLength int) {
	s.maxIdentifierLength = maxIdentifierLength
}

func (s *SchemaConfig) GetMaxIdentifierLength() int {
	return s.maxIdentifierLength
}

// GetName - Gets the default namespace of schema objects.
func (s *SchemaConfig) GetName() *string {
	return s.name
}

// SetName - Sets the default namespace name of schema objects.
func (s *SchemaConfig) SetName(name *string) {
	s.name = name
}

// GetDefaultTableOptions - Gets the default options that are passed to Table instances created with
// Schema#createTable().
func (s *SchemaConfig) GetDefaultTableOptions() map[string]interface{} {
	return s.defaultTableOptions
}

func (s *SchemaConfig) SetDefaultTableOptions(defaultTableOptions map[string]interface{}) {
	s.defaultTableOptions = defaultTableOptions
}
