package enums

type TrimMode string

const (
	TrimModeUnspecified TrimMode = "UNSPECIFIED"
	TrimModeLeading     TrimMode = "LEADING"
	TrimModeTrailing    TrimMode = "TRAILING"
	TrimModeBoth        TrimMode = "BOTH"
)
