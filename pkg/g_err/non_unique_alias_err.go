package g_err

import (
	"fmt"
	"strings"
)

type NonUniqueAlias struct {
	*QueryException
}

func NewNonUniqueAlias(
	alias string,
	registeredAliases []string,
) *NonUniqueAlias {
	message := fmt.Sprintf(
		"The given alias \"%s\" is not unique in FROM and JOIN clause table. "+
			"The currently registered aliases are: %s.",
		alias,
		strings.Join(registeredAliases, ", "),
	)

	return &NonUniqueAlias{
		QueryException: NewQueryException(message),
	}
}
