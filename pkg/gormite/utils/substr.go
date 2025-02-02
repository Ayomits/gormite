package utils

// Substr - Returns a substring of a string.
// NOTE: this isn't multi-Unicode-codepoint aware, like specifying skintone or
//
//	gender of an emoji: https://unicode.org/emoji/charts/full-emoji-modifiers.html
func Substr(input string, offset int, length int) string {
	asRunes := []rune(input)

	if offset >= len(asRunes) {
		return ""
	}

	if offset+length > len(asRunes) {
		length = len(asRunes) - offset
	}

	return string(asRunes[offset : offset+length])
}
