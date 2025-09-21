package types

import (
	"strconv"
)

type DecimalType struct{ *AbstractType }

func NewDecimalType() *DecimalType {
	return &DecimalType{AbstractType: &AbstractType{}}
}

func (d *DecimalType) GetSQLDeclaration(column map[string]any, platform TypesPlatform) string {
	return platform.GetDecimalTypeDeclarationSQL(column)
}

func (d *DecimalType) ConvertToPHPValue(value any, _ TypesPlatform) any {
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
