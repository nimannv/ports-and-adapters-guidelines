package board

import "errors"

var (
	// ErrBoardNotFound is returned when a requested board is not found.
	ErrBoardNotFound = errors.New("board not found")

	// ErrColumnNotFound is returned when a requested column is not found on the board.
	ErrColumnNotFound = errors.New("column not found")

	// ErrDuplicateColumnName is returned when a column with the same name already exists on the board.
	ErrDuplicateColumnName = errors.New("duplicate column name")

	// ErrInvalidWipLimit is returned when an invalid WIP limit is supplied.
	ErrInvalidWipLimit = errors.New("invalid wip limit")

	// ErrInvalidBoard is returned when board parameters fail domain validation.
	ErrInvalidBoard = errors.New("invalid board")
)

