package types

type TimeImmutableType struct{ *AbstractType }

func (t *TimeImmutableType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetTimeTypeDeclarationSQL(column)
}
func (t *TimeImmutableType) ConvertToDatabaseValue(
	value any,
	platform TypesPlatform,
) any {
	panic("not implemented")
}
func (t *TimeImmutableType) ConvertToPHPValue(
	value any,
	platform TypesPlatform,
) any {
	panic("not implemented")
}
