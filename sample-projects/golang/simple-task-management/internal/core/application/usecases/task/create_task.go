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

var _ taskin.CreateTaskUseCase = (*CreateTaskService)(nil)

type CreateTaskService struct {
	boardRepo boardout.Repository
	taskRepo  taskout.Repository
}

func NewCreateTaskService(boardRepo boardout.Repository, taskRepo taskout.Repository) *CreateTaskService {
	return &CreateTaskService{
		boardRepo: boardRepo,
		taskRepo:  taskRepo,
	}
}

func (s *CreateTaskService) Exec(ctx context.Context, in taskin.CreateTaskInput) (taskin.CreateTaskOutput, error) {
	boardRec, err := s.boardRepo.GetByID(ctx, in.BoardID)
	if err != nil {
		if errors.Is(err, boardout.ErrNotFound) {
			return taskin.CreateTaskOutput{}, fmt.Errorf("%w: board %q", taskin.ErrBoardNotFound, in.BoardID)
		}
		return taskin.CreateTaskOutput{}, fmt.Errorf("get board %q: %w", in.BoardID, err)
	}

	// Rehydrate domain board to resolve column
	domainCols := make([]boarddomain.Column, len(boardRec.Columns))
	for i, c := range boardRec.Columns {
		domainCols[i] = boarddomain.RehydrateColumn(c.ID, c.Name, c.OrderIndex, c.WipLimit)
	}
	b := boarddomain.RehydrateBoard(boardRec.ID, boardRec.Name, boardRec.Description, domainCols, boardRec.CreatedAt, boardRec.UpdatedAt)

	colRef := strings.TrimSpace(in.Column)
	if colRef == "" {
		// Default to first column if empty
		if len(b.Columns()) > 0 {
			colRef = b.Columns()[0].ID()
		}
	}

	targetCol, err := b.GetColumn(colRef)
	if err != nil {
		return taskin.CreateTaskOutput{}, fmt.Errorf("%w: column %q", taskin.ErrColumnNotFound, colRef)
	}

	// Check WIP limit on destination column
	if targetCol.WipLimit() != nil {
		currentCount, err := s.taskRepo.CountByColumn(ctx, targetCol.ID())
		if err != nil {
			return taskin.CreateTaskOutput{}, fmt.Errorf("count tasks in column: %w", err)
		}
		if targetCol.IsFull(currentCount) {
			return taskin.CreateTaskOutput{}, &taskin.WipLimitExceededError{
				Message:      fmt.Sprintf("Cannot create task in %q: column limit of %d tasks has been reached.", targetCol.Name(), *targetCol.WipLimit()),
				ColumnID:     targetCol.ID(),
				CurrentCount: currentCount,
				Limit:        *targetCol.WipLimit(),
			}
		}
	}

	prio, err := taskdomain.ParsePriority(in.Priority)
	if err != nil {
		return taskin.CreateTaskOutput{}, fmt.Errorf("%w: %v", taskin.ErrInvalidTask, err)
	}

	now := time.Now().UTC()
	t, err := taskdomain.NewTask("", b.ID(), targetCol.ID(), in.Title, in.Description, prio, in.DueDate, now)
	if err != nil {
		return taskin.CreateTaskOutput{}, translateDomainTaskError(err)
	}

	rec := taskout.TaskRecord{
		ID:          t.ID(),
		BoardID:     t.BoardID(),
		ColumnID:    t.ColumnID(),
		Title:       t.Title(),
		Description: t.Description(),
		Priority:    string(t.Priority()),
		DueDate:     t.DueDate(),
		CreatedAt:   t.CreatedAt(),
		UpdatedAt:   t.UpdatedAt(),
	}

	saved, err := s.taskRepo.Save(ctx, rec)
	if err != nil {
		return taskin.CreateTaskOutput{}, fmt.Errorf("save task: %w", err)
	}

	return taskin.CreateTaskOutput{
		ID:          saved.ID,
		BoardID:     saved.BoardID,
		ColumnID:    saved.ColumnID,
		ColumnName:  targetCol.Name(),
		Title:       saved.Title,
		Description: saved.Description,
		Priority:    saved.Priority,
		DueDate:     saved.DueDate,
		CreatedAt:   saved.CreatedAt,
		UpdatedAt:   saved.UpdatedAt,
	}, nil
}

func translateDomainTaskError(err error) error {
	switch {
	case errors.Is(err, taskdomain.ErrInvalidTask):
		return fmt.Errorf("%w: %v", taskin.ErrInvalidTask, err)
	case errors.Is(err, taskdomain.ErrInvalidTaskMovement):
		return fmt.Errorf("%w: %v", taskin.ErrInvalidTaskMovement, err)
	case errors.Is(err, taskdomain.ErrTaskNotFound):
		return fmt.Errorf("%w: %v", taskin.ErrTaskNotFound, err)
	default:
		var wip *taskdomain.WipLimitExceededError
		if errors.As(err, &wip) {
			return &taskin.WipLimitExceededError{
				Message:      wip.Error(),
				ColumnID:     wip.ColumnID,
				CurrentCount: wip.CurrentCount,
				Limit:        wip.Limit,
			}
		}
		return err
	}
}

