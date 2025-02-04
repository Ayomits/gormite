package assets

// Identifier - An abstraction class for an asset identifier.
// Wraps identifier names like column names in indexes / foreign keys
// in an abstract class for proper quotation capabilities.
type Identifier struct {
	*AbstractAsset
}

// newIdentifier - Creates a new identifier.
// identifier - Identifier name to wrap.
// quote - by default false - Whether to force quoting the given identifier.
func newIdentifier(identifier string, quote bool) *Identifier {
	v := &Identifier{AbstractAsset: NewAbstractAsset()}

	v.SetName(identifier)

	if !quote || v.IsQuoted() {
		return v
	}

	v.SetName(`"` + v.GetName() + `"`)

	return v
}

func NewIdentifier(identifier string) *Identifier {
	return newIdentifier(identifier, false)
}

func NewIdentifierWithQuote(identifier string) *Identifier {
	return newIdentifier(identifier, true)
}
