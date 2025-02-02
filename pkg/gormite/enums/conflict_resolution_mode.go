package enums

type ConflictResolutionMode string

const (
	// ConflictResolutionModeOrdinary - Wait for the row to be unlocked
	ConflictResolutionModeOrdinary ConflictResolutionMode = "ORDINARY"

	// ConflictResolutionModeSkipLocked - Skip the row if it is locked
	ConflictResolutionModeSkipLocked ConflictResolutionMode = "SKIP_LOCKED"
)
