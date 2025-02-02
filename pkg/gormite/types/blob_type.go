package types

import (
	"github.com/KoNekoD/gormite/pkg/gormite/enums"
)

type BlobType struct{ *AbstractType }

func (b *BlobType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetBlobTypeDeclarationSQL(column)
}

func (b *BlobType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	if value == nil {
		return nil
	}

	panic("Not implemented, needed to write im temp dir and read then")
}

func (b *BlobType) GetBindingType() enums.ParameterType {
	return enums.ParameterTypeLargeObject
}
