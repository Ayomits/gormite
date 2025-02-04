package types

import (
	"fmt"
)

type DecimalType struct{ *AbstractType }

func (d *DecimalType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetDecimalTypeDeclarationSQL(column)
}
func (d *DecimalType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	switch v := value.(type) {
	case float32, float64, int:
		return fmt.Sprintf("%f", v)
	}
	return value
}
