package board

import "context"

// AddColumnInput is the input DTO for adding a column to an existing board.
type AddColumnInput struct {
	BoardID  string
	Name     string
	WipLimit *int
}

// AddColumnOutput is the output DTO after successfully adding a column.
type AddColumnOutput struct {
	BoardID string
	Column  ColumnDTO
}

// AddColumnUseCase defines the contract for adding a column to a board.
type AddColumnUseCase interface {
	Exec(ctx context.Context, in AddColumnInput) (AddColumnOutput, error)
}

