package task

import (
	"context"
	"errors"
	"fmt"
	"strings"

	boardout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/board"
	taskin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/task"
	taskout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/task"
	boarddomain "github.com/ports-and-adapters/simple-task-management/internal/core/domain/board"
)

var _ taskin.ListTasksUseCase = (*ListTasksService)(nil)

type ListTasksService struct {
	boardRepo boardout.Repository
	taskRepo  taskout.Repository
}

func NewListTasksService(boardRepo boardout.Repository, taskRepo taskout.Repository) *ListTasksService {
	return &ListTasksService{
		boardRepo: boardRepo,
		taskRepo:  taskRepo,
	}
}

func (s *ListTasksService) Exec(ctx context.Context, in taskin.ListTasksInput) (taskin.ListTasksOutput, error) {
	boardRec, err := s.boardRepo.GetByID(ctx, in.BoardID)
	if err != nil {
		if errors.Is(err, boardout.ErrNotFound) {
			return taskin.ListTasksOutput{}, fmt.Errorf("%w: board %q", taskin.ErrBoardNotFound, in.BoardID)
		}
		return taskin.ListTasksOutput{}, fmt.Errorf("get board %q: %w", in.BoardID, err)
	}

	// Build map of column id -> name
	columnNameByID := make(map[string]string, len(boardRec.Columns))
	for _, c := range boardRec.Columns {
		columnNameByID[c.ID] = c.Name
	}

	filterColumnID := strings.TrimSpace(in.ColumnID)
	if filterColumnID == "" && strings.TrimSpace(in.ColumnName) != "" {
		// Resolve column name
		domainCols := make([]boarddomain.Column, len(boardRec.Columns))
		for i, c := range boardRec.Columns {
			domainCols[i] = boarddomain.RehydrateColumn(c.ID, c.Name, c.OrderIndex, c.WipLimit)
		}
		b := boarddomain.RehydrateBoard(boardRec.ID, boardRec.Name, boardRec.Description, domainCols, boardRec.CreatedAt, boardRec.UpdatedAt)
		col, err := b.GetColumn(in.ColumnName)
		if err != nil {
			return taskin.ListTasksOutput{}, fmt.Errorf("%w: column %q", taskin.ErrColumnNotFound, in.ColumnName)
		}
		filterColumnID = col.ID()
	}

	tasks, err := s.taskRepo.ListByBoard(ctx, in.BoardID, filterColumnID)
	if err != nil {
		return taskin.ListTasksOutput{}, fmt.Errorf("list tasks: %w", err)
	}

	outTasks := make([]taskin.TaskItemDTO, len(tasks))
	for i, t := range tasks {
		cName := columnNameByID[t.ColumnID]
		if cName == "" {
			cName = t.ColumnID
		}
		outTasks[i] = taskin.TaskItemDTO{
			ID:          t.ID,
			BoardID:     t.BoardID,
			ColumnID:    t.ColumnID,
			ColumnName:  cName,
			Title:       t.Title,
			Description: t.Description,
			Priority:    t.Priority,
			DueDate:     t.DueDate,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}

	return taskin.ListTasksOutput{
		BoardID: in.BoardID,
		Tasks:   outTasks,
	}, nil
}

