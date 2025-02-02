package dtos

import (
	"github.com/KoNekoD/gormite/pkg/gormite/enums"
)

type ForUpdate struct {
	conflictResolutionMode enums.ConflictResolutionMode
}

func NewForUpdate(conflictResolutionMode enums.ConflictResolutionMode) *ForUpdate {
	return &ForUpdate{conflictResolutionMode: conflictResolutionMode}
}

func (f *ForUpdate) GetConflictResolutionMode() enums.ConflictResolutionMode {
	return f.conflictResolutionMode
}
