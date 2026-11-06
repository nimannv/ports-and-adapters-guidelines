package board

import (
	"context"
	"time"
)

// ListBoardsInput is the input DTO for listing boards.
type ListBoardsInput struct{}

// BoardItemDTO represents an item in the boards list.
type BoardItemDTO struct {
	ID          string
	Name        string
	Description string
	Columns     []ColumnDTO
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ListBoardsOutput is the output DTO containing all boards.
type ListBoardsOutput struct {
	Boards []BoardItemDTO
}

// ListBoardsUseCase defines the contract for listing all boards.
type ListBoardsUseCase interface {
	Exec(ctx context.Context, in ListBoardsInput) (ListBoardsOutput, error)
}

