package types

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/gormite/enums"
	"io"
)

type BinaryType struct{ *AbstractType }

func (b *BinaryType) GetSQLDeclaration(
	column map[string]interface{},
	platform TypesPlatform,
) string {
	return platform.GetBinaryTypeDeclarationSQL(column)
}
func (b *BinaryType) ConvertToPHPValue(value any, platform TypesPlatform) any {
	if value == nil {
		return nil
	}

	reader, ok := value.(io.Reader)

	if ok {
		bytes, _ := io.ReadAll(reader)
		str := string(bytes)
		return &str
	}

	str, ok := value.(string)

	if !ok {
		panic(fmt.Sprintf("Expected string or io.Reader, got %T", value))
	}

	return &str
}
func (b *BinaryType) GetBindingType() enums.ParameterType {
	return enums.ParameterTypeBinary
}
