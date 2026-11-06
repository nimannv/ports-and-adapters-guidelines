package task

import "errors"

var (
	// ErrNotFound is returned by outbound storage when a task is not found.
	ErrNotFound = errors.New("task not found in repository")
)

