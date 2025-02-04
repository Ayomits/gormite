package diff_calc

import (
	"github.com/KoNekoD/gormite/pkg/assets"
)

type DiffCalcPlatform interface {
	ColumnsEqual(column1 *assets.Column, column2 *assets.Column) bool
}
