package diff_calc

import (
	"github.com/KoNekoD/gormite/pkg/gormite/assets"
)

type DiffCalcPlatform interface {
	ColumnsEqual(column1 *assets.Column, column2 *assets.Column) bool
}
