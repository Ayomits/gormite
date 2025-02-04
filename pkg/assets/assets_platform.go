package assets

import (
	"github.com/KoNekoD/gormite/pkg/keywords"
	"github.com/KoNekoD/gormite/pkg/supports_platforms_contracts"
)

type AssetsPlatform interface {
	supports_platforms_contracts.SupportsPlatform
	GetReservedKeywordsList() keywords.KeywordListInterface
	QuoteIdentifier(identifier string) string
}
