package task

import (
	"context"
	"time"
)

// MoveTaskInput is the input DTO for moving a task to another column.
type MoveTaskInput struct {
	TaskID           string
	TargetColumnID   string
	TargetColumnName string
}

// MoveTaskOutput is the output DTO after moving a task.
type MoveTaskOutput struct {
	TaskID     string
	BoardID    string
	ColumnID   string
	ColumnName string
	Title      string
	Priority   string
	UpdatedAt  time.Time
}

// MoveTaskUseCase defines the contract for moving a task between columns.
type MoveTaskUseCase interface {
	Exec(ctx context.Context, in MoveTaskInput) (MoveTaskOutput, error)
}

