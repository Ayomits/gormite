package types

type DateTimeImmutableType struct{ *AbstractType }

func NewDateTimeImmutableType() *DateTimeImmutableType {
	return &DateTimeImmutableType{AbstractType: &AbstractType{}}
}

func (d *DateTimeImmutableType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetDateTimeTypeDeclarationSQL(column)
}
func (d *DateTimeImmutableType) ConvertToDatabaseValue(
	value any,
	platform TypesPlatform,
) any {
	panic("not implemented")
}
func (d *DateTimeImmutableType) ConvertToPHPValue(
	value any,
	platform TypesPlatform,
) any {
	panic("not implemented")
}
