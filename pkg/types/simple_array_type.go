package types

type SimpleArrayType struct{ *AbstractType }

func (s *SimpleArrayType) GetSQLDeclaration(column map[string]interface{}, platform TypesPlatform) string {
	return platform.GetClobTypeDeclarationSQL(column)
}
func (s *SimpleArrayType) ConvertToDatabaseValue(value any, platform TypesPlatform) *string {
	panic("not implemented")
}
