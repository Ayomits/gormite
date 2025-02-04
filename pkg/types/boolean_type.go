package types

import (
	"github.com/KoNekoD/gormite/pkg/enums"
)

type BooleanType struct{ *AbstractType }

func NewBooleanType() *BooleanType {
	return &BooleanType{AbstractType: &AbstractType{}}
}

func (b *BooleanType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetBooleanTypeDeclarationSQL(column)
}
func (b *BooleanType) ConvertToDatabaseValue(
	value any,
	platform TypesPlatform,
) any {
	return platform.ConvertBooleansToDatabaseValue(value)
}
func (b *BooleanType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	return platform.ConvertFromBoolean(value)
}
func (b *BooleanType) GetBindingType() enums.ParameterType {
	return enums.ParameterTypeBoolean
}
