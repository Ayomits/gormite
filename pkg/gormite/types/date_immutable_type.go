package types

import (
	"fmt"
	"time"
)

type DateImmutableType struct{ *AbstractType }

func (d *DateImmutableType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetDateTypeDeclarationSQL(column)
}
func (d *DateImmutableType) ConvertToDatabaseValue(value any, platform TypesPlatform) *string {
	if value == nil {
		return nil
	}

	dateTime, ok := value.(*time.Time)

	if ok {
		str := dateTime.Format(platform.GetDateFormatString())
		return &str
	}

	panic(fmt.Sprintf("Expected time.Time, got %T", value))
}
func (d *DateImmutableType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	if value == nil {
		return nil
	}

	dateTime, ok := value.(*time.Time)
	if ok {
		return dateTime
	}

	str, ok := value.(string)
	if !ok {
		panic(fmt.Sprintf("Expected string or time.Time, got %T", value))
	}

	dateTimeParsed, err := time.Parse(platform.GetDateFormatString(), str)

	if err != nil {
		panic(fmt.Sprintf("Convert to time error %s", value))
	}

	return &dateTimeParsed
}
