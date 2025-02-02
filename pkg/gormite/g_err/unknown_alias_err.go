package g_err

import (
	"fmt"
	"strings"
)

type UnknownAlias struct {
	*QueryException
}

func NewUnknownAlias(
	alias string,
	registeredAliases []string,
) *UnknownAlias {
	message := fmt.Sprintf(
		"The given alias \"%s\" is not part of any FROM or JOIN clause table. "+
			"The currently registered aliases are: %s.",
		alias,
		strings.Join(registeredAliases, ", "),
	)

	return &UnknownAlias{
		QueryException: NewQueryException(message),
	}
}
