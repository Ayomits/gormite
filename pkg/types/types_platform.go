package types

type TypesPlatform interface {
	GetBooleanTypeDeclarationSQL(column map[string]interface{}) string
	GetIntegerTypeDeclarationSQL(column map[string]interface{}) string
	GetBigIntTypeDeclarationSQL(column map[string]interface{}) string
	GetSmallIntTypeDeclarationSQL(column map[string]interface{}) string
	GetCommonIntegerTypeDeclarationSQL(column map[string]interface{}) string
	GetClobTypeDeclarationSQL(column map[string]interface{}) string
	GetBlobTypeDeclarationSQL(column map[string]interface{}) string
	GetDateTimeTypeDeclarationSQL(column map[string]interface{}) string
	GetDateTypeDeclarationSQL(column map[string]interface{}) string
	GetTimeTypeDeclarationSQL(column map[string]interface{}) string
	GetAsciiStringTypeDeclarationSQL(column map[string]interface{}) string
	GetBinaryTypeDeclarationSQL(column map[string]interface{}) string
	GetStringTypeDeclarationSQL(column map[string]interface{}) string
	GetDecimalTypeDeclarationSQL(column map[string]interface{}) string
	GetJsonTypeDeclarationSQL(column map[string]interface{}) string
	GetGuidTypeDeclarationSQL(column map[string]interface{}) string
	GetFloatTypeDeclarationSQL(column map[string]interface{}) string
	GetDateTimeTzTypeDeclarationSQL(column map[string]interface{}) string
	GetSmallFloatTypeDeclarationSQL(column map[string]interface{}) string
	GetCharTypeDeclarationSQLSnippet(length *int) string
	GetVarcharTypeDeclarationSQLSnippet(length *int) string
	GetBinaryTypeDeclarationSQLSnippet(length *int) string
	GetVarbinaryTypeDeclarationSQLSnippet(length *int) string
	ConvertBooleansToDatabaseValue(item any) any
	ConvertFromBoolean(item any) *bool
	GetDateFormatString() string
}
