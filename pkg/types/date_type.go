package types

type DateType struct{ *AbstractType }

func (d *DateType) GetSQLDeclaration(column map[string]any, platform TypesPlatform) string {
	return platform.GetDateTypeDeclarationSQL(column)
}
func (d *DateType) ConvertToDatabaseValue(value any, platform TypesPlatform) any {
	panic("not implemented")
}
func (d *DateType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	panic("not implemented")
}
