package types

import (
	"strconv"
)

type FloatType struct{ *AbstractType }

func (f *FloatType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetFloatTypeDeclarationSQL(column)
}
func (f *FloatType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	if value == nil {
		return nil
	}

	if v, ok := value.(string); ok {
		if v == "" {
			return nil
		}
		v, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		return &v
	}

	panic("unknown type")
}
