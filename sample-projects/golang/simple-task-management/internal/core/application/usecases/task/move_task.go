package task

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	boardout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/board"
	taskin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/task"
	taskout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/task"
	boarddomain "github.com/ports-and-adapters/simple-task-management/internal/core/domain/board"
	taskdomain "github.com/ports-and-adapters/simple-task-management/internal/core/domain/task"
)

var _ taskin.MoveTaskUseCase = (*MoveTaskService)(nil)

type MoveTaskService struct {
	boardRepo boardout.Repository
	taskRepo  taskout.Repository
}

func NewMoveTaskService(boardRepo boardout.Repository, taskRepo taskout.Repository) *MoveTaskService {
	return &MoveTaskService{
		boardRepo: boardRepo,
		taskRepo:  taskRepo,
	}
}

func (s *MoveTaskService) Exec(ctx context.Context, in taskin.MoveTaskInput) (taskin.MoveTaskOutput, error) {
	taskRec, err := s.taskRepo.GetByID(ctx, in.TaskID)
	if err != nil {
		if errors.Is(err, taskout.ErrNotFound) {
			return taskin.MoveTaskOutput{}, fmt.Errorf("%w: task %q", taskin.ErrTaskNotFound, in.TaskID)
		}
		return taskin.MoveTaskOutput{}, fmt.Errorf("get task %q: %w", in.TaskID, err)
	}

	boardRec, err := s.boardRepo.GetByID(ctx, taskRec.BoardID)
	if err != nil {
		if errors.Is(err, boardout.ErrNotFound) {
			return taskin.MoveTaskOutput{}, fmt.Errorf("%w: parent board %q", taskin.ErrBoardNotFound, taskRec.BoardID)
		}
		return taskin.MoveTaskOutput{}, fmt.Errorf("get parent board: %w", err)
	}

	// Rehydrate domain board
	domainCols := make([]boarddomain.Column, len(boardRec.Columns))
	for i, c := range boardRec.Columns {
		domainCols[i] = boarddomain.RehydrateColumn(c.ID, c.Name, c.OrderIndex, c.WipLimit)
	}
	b := boarddomain.RehydrateBoard(boardRec.ID, boardRec.Name, boardRec.Description, domainCols, boardRec.CreatedAt, boardRec.UpdatedAt)

	targetRef := strings.TrimSpace(in.TargetColumnID)
	if targetRef == "" {
		targetRef = strings.TrimSpace(in.TargetColumnName)
	}
	if targetRef == "" {
		return taskin.MoveTaskOutput{}, fmt.Errorf("%w: target column must be specified", taskin.ErrInvalidTaskMovement)
	}

	targetCol, err := b.GetColumn(targetRef)
	if err != nil {
		return taskin.MoveTaskOutput{}, fmt.Errorf("%w: target column %q", taskin.ErrColumnNotFound, targetRef)
	}

	// If moving to same column, return current state
	if taskRec.ColumnID == targetCol.ID() {
		return taskin.MoveTaskOutput{
			TaskID:     taskRec.ID,
			BoardID:    taskRec.BoardID,
			ColumnID:   taskRec.ColumnID,
			ColumnName: targetCol.Name(),
			Title:      taskRec.Title,
			Priority:   taskRec.Priority,
			UpdatedAt:  taskRec.UpdatedAt,
		}, nil
	}

	// Check WIP limit on destination column
	if targetCol.WipLimit() != nil {
		currentCount, err := s.taskRepo.CountByColumn(ctx, targetCol.ID())
		if err != nil {
			return taskin.MoveTaskOutput{}, fmt.Errorf("count tasks in destination column: %w", err)
		}
		if targetCol.IsFull(currentCount) {
			return taskin.MoveTaskOutput{}, &taskin.WipLimitExceededError{
				Message:      fmt.Sprintf("Cannot move task to %q: column limit of %d tasks has been reached.", targetCol.Name(), *targetCol.WipLimit()),
				ColumnID:     targetCol.ID(),
				CurrentCount: currentCount,
				Limit:        *targetCol.WipLimit(),
			}
		}
	}

	// Rehydrate task domain and move
	t := taskdomain.RehydrateTask(
		taskRec.ID,
		taskRec.BoardID,
		taskRec.ColumnID,
		taskRec.Title,
		taskRec.Description,
		taskdomain.Priority(taskRec.Priority),
		taskRec.DueDate,
		taskRec.CreatedAt,
		taskRec.UpdatedAt,
	)

	now := time.Now().UTC()
	if err := t.MoveTo(targetCol.ID(), b.ID(), now); err != nil {
		return taskin.MoveTaskOutput{}, translateDomainTaskError(err)
	}

	updatedRec, err := s.taskRepo.UpdateColumn(ctx, t.ID(), targetCol.ID(), now)
	if err != nil {
		if errors.Is(err, taskout.ErrNotFound) {
			return taskin.MoveTaskOutput{}, fmt.Errorf("%w: task %q", taskin.ErrTaskNotFound, in.TaskID)
		}
		return taskin.MoveTaskOutput{}, fmt.Errorf("update task column: %w", err)
	}

	return taskin.MoveTaskOutput{
		TaskID:     updatedRec.ID,
		BoardID:    updatedRec.BoardID,
		ColumnID:   updatedRec.ColumnID,
		ColumnName: targetCol.Name(),
		Title:      updatedRec.Title,
		Priority:   updatedRec.Priority,
		UpdatedAt:  updatedRec.UpdatedAt,
	}, nil
}

