package task

import "context"

// DeleteTaskInput is the input DTO for deleting a task.
type DeleteTaskInput struct {
	TaskID string
}

// DeleteTaskOutput is the output DTO after deleting a task.
type DeleteTaskOutput struct {
	TaskID  string
	Deleted bool
}

// DeleteTaskUseCase defines the contract for deleting a task.
type DeleteTaskUseCase interface {
	Exec(ctx context.Context, in DeleteTaskInput) (DeleteTaskOutput, error)
}

