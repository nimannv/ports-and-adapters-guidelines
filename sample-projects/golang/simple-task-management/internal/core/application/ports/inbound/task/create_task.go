package task

import (
	"context"
	"time"
)

// CreateTaskInput is the input DTO for creating a new task.
type CreateTaskInput struct {
	BoardID     string
	Column      string // Can be column ID or column Name
	Title       string
	Description string
	Priority    string
	DueDate     *time.Time
}

// CreateTaskOutput is the output DTO returned after creating a task.
type CreateTaskOutput struct {
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

// CreateTaskUseCase defines the contract for creating a task.
type CreateTaskUseCase interface {
	Exec(ctx context.Context, in CreateTaskInput) (CreateTaskOutput, error)
}

