package types

type SmallFloatType struct{ *AbstractType }

func (s *SmallFloatType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetSmallFloatTypeDeclarationSQL(column)
}
func (s *SmallFloatType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	panic("not implemented")
}
