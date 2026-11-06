package task

import (
	"context"
	"errors"
	"fmt"

	taskin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/task"
	taskout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/task"
)

var _ taskin.DeleteTaskUseCase = (*DeleteTaskService)(nil)

type DeleteTaskService struct {
	taskRepo taskout.Repository
}

func NewDeleteTaskService(taskRepo taskout.Repository) *DeleteTaskService {
	return &DeleteTaskService{taskRepo: taskRepo}
}

func (s *DeleteTaskService) Exec(ctx context.Context, in taskin.DeleteTaskInput) (taskin.DeleteTaskOutput, error) {
	err := s.taskRepo.Delete(ctx, in.TaskID)
	if err != nil {
		if errors.Is(err, taskout.ErrNotFound) {
			return taskin.DeleteTaskOutput{}, fmt.Errorf("%w: task %q", taskin.ErrTaskNotFound, in.TaskID)
		}
		return taskin.DeleteTaskOutput{}, fmt.Errorf("delete task: %w", err)
	}

	return taskin.DeleteTaskOutput{
		TaskID:  in.TaskID,
		Deleted: true,
	}, nil
}

