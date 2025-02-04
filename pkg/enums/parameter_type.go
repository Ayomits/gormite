package enums

// ParameterType - Statement parameter type.
type ParameterType string

var (
	// ParameterTypeNull - Represents the SQL NULL data type.
	ParameterTypeNull ParameterType = "NULL"

	// ParameterTypeInteger - Represents the SQL INTEGER data type.
	ParameterTypeInteger ParameterType = "INTEGER"

	// ParameterTypeString - Represents the SQL CHAR, VARCHAR, or other string data type.
	ParameterTypeString ParameterType = "STRING"

	// ParameterTypeLargeObject - Represents the SQL large object data type.
	ParameterTypeLargeObject ParameterType = "LARGE_OBJECT"

	// ParameterTypeBoolean - Represents a boolean data type.
	ParameterTypeBoolean ParameterType = "BOOLEAN"

	// ParameterTypeBinary - Represents a binary string data type.
	ParameterTypeBinary ParameterType = "BINARY"

	// ParameterTypeAscii - Represents an ASCII string data type
	ParameterTypeAscii ParameterType = "ASCII"
)
