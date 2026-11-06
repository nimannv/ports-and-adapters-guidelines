package board

import (
	"context"
	"time"
)

// ColumnDTO represents a column data contract at the port boundary.
type ColumnDTO struct {
	ID         string
	Name       string
	OrderIndex int
	WipLimit   *int
}

// CustomColumnInput specifies a custom initial column when creating a board.
type CustomColumnInput struct {
	Name       string
	OrderIndex int
	WipLimit   *int
}

// CreateBoardInput is the input DTO for creating a board.
type CreateBoardInput struct {
	Name          string
	Description   string
	CustomColumns []CustomColumnInput
}

// CreateBoardOutput is the output DTO after creating a board.
type CreateBoardOutput struct {
	ID          string
	Name        string
	Description string
	Columns     []ColumnDTO
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateBoardUseCase defines the contract for creating a board.
type CreateBoardUseCase interface {
	Exec(ctx context.Context, in CreateBoardInput) (CreateBoardOutput, error)
}

