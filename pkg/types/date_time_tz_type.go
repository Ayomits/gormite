package types

type DateTimeTzType struct{ *AbstractType }

func (d *DateTimeTzType) GetSQLDeclaration(column map[string]any, platform TypesPlatform) string {
	return platform.GetDateTimeTzTypeDeclarationSQL(column)
}
func (d *DateTimeTzType) ConvertToDatabaseValue(value any, platform TypesPlatform) any {
	panic("not implemented")
}
func (d *DateTimeTzType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	panic("not implemented")
}
