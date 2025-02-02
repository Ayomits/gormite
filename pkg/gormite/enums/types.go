package enums

// TypesType - Default built-in types provided by Doctrine DBAL.
type TypesType string

const (
	TypeAsciiString         TypesType = "ascii_string"
	TypeBigint              TypesType = "bigint"
	TypeBinary              TypesType = "binary"
	TypeBlob                TypesType = "blob"
	TypeBoolean             TypesType = "boolean"
	TypeDateMutable         TypesType = "date_mutable"
	TypeDateImmutable       TypesType = "date_immutable"
	TypeDateinterval        TypesType = "dateinterval"
	TypeDatetimeMutable     TypesType = "datetime_mutable"
	TypeDatetimeImmutable   TypesType = "datetime_immutable"
	TypeDatetimetzMutable   TypesType = "datetimetz_mutable"
	TypeDatetimetzImmutable TypesType = "datetimetz_immutable"
	TypeDecimal             TypesType = "decimal"
	TypeFloat               TypesType = "float"
	TypeGuid                TypesType = "guid"
	TypeInteger             TypesType = "integer"
	TypeJson                TypesType = "json"
	TypeSimpleArray         TypesType = "simple_array"
	TypeSmallfloat          TypesType = "smallfloat"
	TypeSmallint            TypesType = "smallint"
	TypeString              TypesType = "string"
	TypeText                TypesType = "text"
	TypeTimeMutable         TypesType = "time_mutable"
	TypeTimeImmutable       TypesType = "time_immutable"
)
