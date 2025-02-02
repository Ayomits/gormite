package types

import (
	"github.com/KoNekoD/gormite/pkg/gormite/enums"
)

type SmallIntType struct{ *AbstractType }

func (s *SmallIntType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetSmallIntTypeDeclarationSQL(column)
}
func (s *SmallIntType) ConvertToPHPValue(
	value any,
	platform TypesPlatform,
) any {
	panic("not implemented")
}
func (s *SmallIntType) GetBindingType() enums.ParameterType {
	return enums.ParameterTypeInteger
}
