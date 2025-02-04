package local_schema

import (
	"github.com/KoNekoD/gormite/pkg/utils"
	"strings"
)

func getName(store *store, name string) string {
	if value, ok := store.namesMap[name]; ok {
		name = strings.TrimSpace(value)
	}
	return utils.ToSnakeCase(name)
}
