package types

type GuidType struct{ *StringType }

func (g *GuidType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetGuidTypeDeclarationSQL(column)
}
