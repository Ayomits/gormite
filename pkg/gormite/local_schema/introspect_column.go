package local_schema

import (
	"github.com/KoNekoD/gormite/pkg/gormite/assets"
	"github.com/KoNekoD/gormite/pkg/gormite/types"
)

type columnData struct {
	ColumnName   string
	IsPrimaryKey bool
	IsForeignKey bool
	IsNotNull    bool

	TypeName string

	IsUnique          bool
	UniqueName        *string
	IsUniqueCondition bool
	UniqueCondition   *string

	IsIndex          bool
	IndexName        *string
	IsIndexCondition bool
	IndexCondition   *string

	Length       int
	DefaultValue *string

	ColumnType types.AbstractTypeInterface

	Options []assets.ColumnOption
}
