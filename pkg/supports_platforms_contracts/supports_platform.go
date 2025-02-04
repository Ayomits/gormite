package supports_platforms_contracts

type SupportsPlatform interface {
	SupportsColumnLengthIndexes() bool
	SupportsSchemas() bool
	SupportsSequences() bool
	SupportsIdentityColumns() bool
	SupportsPartialIndexes() bool
	SupportsSavepoints() bool
	SupportsReleaseSavepoints() bool
	SupportsInlineColumnComments() bool
	SupportsCommentOnStatement() bool
	SupportsColumnCollation() bool
}
