package board

import "errors"

var (
	// ErrNotFound is returned by outbound storage when a board is not found.
	ErrNotFound = errors.New("board not found in repository")

	// ErrDuplicate is returned by outbound storage when a unique constraint is violated.
	ErrDuplicate = errors.New("board or column already exists in repository")
)

