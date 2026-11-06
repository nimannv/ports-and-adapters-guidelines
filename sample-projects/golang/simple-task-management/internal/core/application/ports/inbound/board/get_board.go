package board

import (
	"context"
	"time"
)

// GetBoardInput is the input DTO for retrieving a board.
type GetBoardInput struct {
	BoardID string
}

// GetBoardOutput is the output DTO for retrieving a board.
type GetBoardOutput struct {
	ID          string
	Name        string
	Description string
	Columns     []ColumnDTO
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GetBoardUseCase defines the contract for retrieving a board by ID.
type GetBoardUseCase interface {
	Exec(ctx context.Context, in GetBoardInput) (GetBoardOutput, error)
}

