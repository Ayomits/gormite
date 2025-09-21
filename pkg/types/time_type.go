package types

import (
	"time"
)

type TimeType struct{ *AbstractType }

func (t *TimeType) GetSQLDeclaration(column map[string]any, platform TypesPlatform) string {
	return platform.GetTimeTypeDeclarationSQL(column)
}
func (t *TimeType) ConvertToDatabaseValue(value any, platform TypesPlatform) any {
	if value == nil {
		return nil
	}

	if v, ok := value.(*time.Time); ok {
		return v.Format("2006-01-02 15:04:05")
	}

	panic("unknown type")
}
func (t *TimeType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	if value == nil {
		return value
	}
	if v, ok := value.(time.Time); ok {
		return v
	}

	panic("unknown type")
}
