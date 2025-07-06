package types

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/enums"
	"strconv"
)

type BigintType struct{ *AbstractType }

func NewBigintType() *BigintType {
	return &BigintType{AbstractType: &AbstractType{}}
}

func (b *BigintType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetBigIntTypeDeclarationSQL(column)
}

func (b *BigintType) GetBindingType() enums.ParameterType {
	return enums.ParameterTypeString
}

func (b *BigintType) ConvertToPHPValue(value any, platform TypesPlatform) any {
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
