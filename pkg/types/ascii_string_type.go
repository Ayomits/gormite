package types

import (
	"github.com/KoNekoD/gormite/pkg/enums"
)

type AsciiStringType struct{ *StringType }

func (a *AsciiStringType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetAsciiStringTypeDeclarationSQL(column)
}

func (a *AsciiStringType) GetBindingType() enums.ParameterType {
	return enums.ParameterTypeAscii
}
