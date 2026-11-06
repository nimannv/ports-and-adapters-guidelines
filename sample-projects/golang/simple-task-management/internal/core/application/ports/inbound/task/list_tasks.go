package task

import (
	"context"
	"time"
)

// TaskItemDTO represents a task in a list response.
type TaskItemDTO struct {
	ID          string
	BoardID     string
	ColumnID    string
	ColumnName  string
	Title       string
	Description string
	Priority    string
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ListTasksInput is the input DTO for listing tasks on a board.
type ListTasksInput struct {
	BoardID    string
	ColumnID   string
	ColumnName string
}

// ListTasksOutput is the output DTO containing tasks for a board.
type ListTasksOutput struct {
	BoardID string
	Tasks   []TaskItemDTO
}

// ListTasksUseCase defines the contract for listing tasks on a board.
type ListTasksUseCase interface {
	Exec(ctx context.Context, in ListTasksInput) (ListTasksOutput, error)
}

