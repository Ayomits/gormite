package types

type StringType struct{ *AbstractType }

func NewStringType() *StringType {
	return &StringType{AbstractType: &AbstractType{}}
}

func (s *StringType) GetSQLDeclaration(column map[string]any, platform TypesPlatform) string {
	return platform.GetStringTypeDeclarationSQL(column)
}
