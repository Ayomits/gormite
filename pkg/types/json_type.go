package types

import (
	"encoding/json"
)

type JsonType struct{ *AbstractType }

func NewJsonType() *JsonType {
	return &JsonType{AbstractType: &AbstractType{}}
}

func (j *JsonType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetJsonTypeDeclarationSQL(column)
}
func (j *JsonType) ConvertToDatabaseValue(value any, platform TypesPlatform) *string {
	if value == nil {
		return nil
	}

	v, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	str := string(v)
	return &str
}
func (j *JsonType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	if value == nil || value == "" {
		return nil
	}

	if v, ok := value.(string); ok {
		if v == "" {
			return nil
		}
		return v
	}

	var data []map[string]interface{}

	if err := json.Unmarshal([]byte(value.(string)), &data); err != nil {
		panic(err)
	}
	return data
}
