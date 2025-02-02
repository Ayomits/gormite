package keywords

type KeywordListInterface interface {
	IsKeyword(word string) bool
	InitializeKeywords()

	// GetKeywords - Returns the list of keywords.
	GetKeywords() []string
}
