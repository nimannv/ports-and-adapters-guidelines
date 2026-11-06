package board

import "errors"

var (
	// ErrInvalidBoard indicates invalid board attributes (e.g. empty name or no columns).
	ErrInvalidBoard = errors.New("invalid board")

	// ErrDuplicateColumnName indicates a column name conflict within a board.
	ErrDuplicateColumnName = errors.New("duplicate column name")

	// ErrColumnNotFound indicates that the requested column does not exist on the board.
	ErrColumnNotFound = errors.New("column not found")

	// ErrInvalidWipLimit indicates an invalid WIP limit (must be > 0 if specified).
	ErrInvalidWipLimit = errors.New("invalid wip limit")

	// ErrColumnNotEmpty indicates a column cannot be deleted because it still contains tasks.
	ErrColumnNotEmpty = errors.New("column is not empty")
)

