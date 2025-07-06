package types

import "strconv"

type SmallFloatType struct{ *AbstractType }

func NewSmallFloatType() *SmallFloatType {
	return &SmallFloatType{AbstractType: &AbstractType{}}
}

func (s *SmallFloatType) GetSQLDeclaration(
	column map[string]any,
	platform TypesPlatform,
) string {
	return platform.GetSmallFloatTypeDeclarationSQL(column)
}
func (s *SmallFloatType) ConvertToPHPValue(value any, _ TypesPlatform) any {
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
		v2 := float32(v)
		return &v2
	}

	panic("unknown type")
}
