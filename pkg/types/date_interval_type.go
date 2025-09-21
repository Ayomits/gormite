package types

type DateIntervalType struct{ *AbstractType }

func (d *DateIntervalType) GetSQLDeclaration(column map[string]any, platform TypesPlatform) string {
	column["length"] = 255

	return platform.GetStringTypeDeclarationSQL(column)
}
func (d *DateIntervalType) ConvertToDatabaseValue(value any, platform TypesPlatform) any {
	panic("not implemented")
}
func (d *DateIntervalType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	panic("not implemented")
}
