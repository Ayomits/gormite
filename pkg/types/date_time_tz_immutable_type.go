package types

type DateTimeTzImmutableType struct{ *AbstractType }

func (d *DateTimeTzImmutableType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetDateTimeTzTypeDeclarationSQL(column)
}
func (d *DateTimeTzImmutableType) ConvertToDatabaseValue(
	value any,
	platform TypesPlatform,
) any {
	panic("not implemented")
}
func (d *DateTimeTzImmutableType) ConvertToPHPValue(
	value any,
	platform TypesPlatform,
) any {
	panic("not implemented")
}
