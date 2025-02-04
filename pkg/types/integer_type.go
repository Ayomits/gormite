package types

import (
	"github.com/KoNekoD/gormite/pkg/enums"
	"strconv"
)

type IntegerType struct{ *AbstractType }

func NewIntegerType() *IntegerType {
	return &IntegerType{AbstractType: &AbstractType{}}
}

func (i *IntegerType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetIntegerTypeDeclarationSQL(column)
}
func (i *IntegerType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	if value == nil {
		return nil
	}

	if v, ok := value.(string); ok {
		if v == "" {
			return nil
		}

		v, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		return &v
	}

	panic("unknown type")
}
func (i *IntegerType) GetBindingType() enums.ParameterType {
	return enums.ParameterTypeInteger
}
