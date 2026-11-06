package task

import (
	"errors"
	"fmt"
)

var (
	// ErrTaskNotFound indicates that a referenced task does not exist.
	ErrTaskNotFound = errors.New("task not found")

	// ErrInvalidTask indicates that task fields failed validation invariants.
	ErrInvalidTask = errors.New("invalid task")

	// ErrInvalidTaskMovement indicates that moving a task violates board boundaries.
	ErrInvalidTaskMovement = errors.New("invalid task movement")
)

// WipLimitExceededError is returned when a destination column cannot accept more tasks.
type WipLimitExceededError struct {
	ColumnName   string
	ColumnID     string
	CurrentCount int
	Limit        int
}

func (e *WipLimitExceededError) Error() string {
	return fmt.Sprintf("cannot move task to %q: column limit of %d tasks has been reached", e.ColumnName, e.Limit)
}

