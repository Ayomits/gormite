package types

type TextType struct{ *AbstractType }

func NewTextType() *TextType {
	return &TextType{AbstractType: &AbstractType{}}
}

func (t *TextType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetClobTypeDeclarationSQL(column)
}
func (t *TextType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	return value
}
