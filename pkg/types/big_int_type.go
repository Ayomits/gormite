package types

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/enums"
	"strconv"
)

type BigIntType struct{ *AbstractType }

func (b *BigIntType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetBigIntTypeDeclarationSQL(column)
}
func (b *BigIntType) GetBindingType() enums.ParameterType {
	return enums.ParameterTypeString
}

func (b *BigIntType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	if value == nil {
		return value
	}
	v, ok := value.(int)

	if ok {
		return v
	}

	str, ok := value.(string)

	if !ok {
		panic(fmt.Sprintf("Expected int or string, got %T", value))
	}

	v, err := strconv.Atoi(str)

	if err != nil {
		panic(fmt.Sprintf("Convert to int error %s", value))
	}

	return value
}
