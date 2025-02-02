package types

type DateTimeType struct{ *AbstractType }

func (d *DateTimeType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetDateTimeTypeDeclarationSQL(column)
}
func (d *DateTimeType) ConvertToDatabaseValue(value any, platform TypesPlatform) any {
	panic("not implemented")
}
func (d *DateTimeType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	panic("not implemented")
}
