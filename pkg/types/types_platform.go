package types

type TypesPlatform interface {
	GetBooleanTypeDeclarationSQL(column map[string]any) string
	GetIntegerTypeDeclarationSQL(column map[string]any) string
	GetBigIntTypeDeclarationSQL(column map[string]any) string
	GetSmallIntTypeDeclarationSQL(column map[string]any) string
	GetCommonIntegerTypeDeclarationSQL(column map[string]any) string
	GetClobTypeDeclarationSQL(column map[string]any) string
	GetBlobTypeDeclarationSQL(column map[string]any) string
	GetDateTimeTypeDeclarationSQL(column map[string]any) string
	GetDateTypeDeclarationSQL(column map[string]any) string
	GetTimeTypeDeclarationSQL(column map[string]any) string
	GetAsciiStringTypeDeclarationSQL(column map[string]any) string
	GetBinaryTypeDeclarationSQL(column map[string]any) string
	GetStringTypeDeclarationSQL(column map[string]any) string
	GetDecimalTypeDeclarationSQL(column map[string]any) string
	GetJsonTypeDeclarationSQL(column map[string]any) string
	GetGuidTypeDeclarationSQL(column map[string]any) string
	GetFloatTypeDeclarationSQL(column map[string]any) string
	GetDateTimeTzTypeDeclarationSQL(column map[string]any) string
	GetSmallFloatTypeDeclarationSQL(column map[string]any) string
	GetCharTypeDeclarationSQLSnippet(length *int) string
	GetVarcharTypeDeclarationSQLSnippet(length *int) string
	GetBinaryTypeDeclarationSQLSnippet(length *int) string
	GetVarbinaryTypeDeclarationSQLSnippet(length *int) string
	ConvertBooleansToDatabaseValue(item any) any
	ConvertFromBoolean(item any) *bool
	GetDateFormatString() string
}
