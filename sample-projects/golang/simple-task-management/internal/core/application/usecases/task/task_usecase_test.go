package task_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	boardout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/board"
	taskin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/task"
	taskout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/task"
	taskuc "github.com/ports-and-adapters/simple-task-management/internal/core/application/usecases/task"
)

type fakeBoardRepo struct {
	boards map[string]boardout.BoardRecord
}

func (f *fakeBoardRepo) Save(_ context.Context, r boardout.BoardRecord) (boardout.BoardRecord, error) {
	f.boards[r.ID] = r
	return r, nil
}

func (f *fakeBoardRepo) GetByID(_ context.Context, id string) (boardout.BoardRecord, error) {
	r, ok := f.boards[id]
	if !ok {
		return boardout.BoardRecord{}, boardout.ErrNotFound
	}
	return r, nil
}

func (f *fakeBoardRepo) ListAll(_ context.Context) ([]boardout.BoardRecord, error) {
	var res []boardout.BoardRecord
	for _, b := range f.boards {
		res = append(res, b)
	}
	return res, nil
}

func (f *fakeBoardRepo) AddColumn(_ context.Context, boardID string, col boardout.ColumnRecord) (boardout.ColumnRecord, error) {
	b, ok := f.boards[boardID]
	if !ok {
		return boardout.ColumnRecord{}, boardout.ErrNotFound
	}
	b.Columns = append(b.Columns, col)
	f.boards[boardID] = b
	return col, nil
}

type fakeTaskRepo struct {
	mu    sync.Mutex
	tasks map[string]taskout.TaskRecord
}

func newFakeTaskRepo() *fakeTaskRepo {
	return &fakeTaskRepo{tasks: make(map[string]taskout.TaskRecord)}
}

func (f *fakeTaskRepo) Save(_ context.Context, r taskout.TaskRecord) (taskout.TaskRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tasks[r.ID] = r
	return r, nil
}

func (f *fakeTaskRepo) GetByID(_ context.Context, id string) (taskout.TaskRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	r, ok := f.tasks[id]
	if !ok {
		return taskout.TaskRecord{}, taskout.ErrNotFound
	}
	return r, nil
}

func (f *fakeTaskRepo) ListByBoard(_ context.Context, boardID string, columnID string) ([]taskout.TaskRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var res []taskout.TaskRecord
	for _, t := range f.tasks {
		if t.BoardID == boardID {
			if columnID == "" || t.ColumnID == columnID {
				res = append(res, t)
			}
		}
	}
	return res, nil
}

func (f *fakeTaskRepo) CountByColumn(_ context.Context, columnID string) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	count := 0
	for _, t := range f.tasks {
		if t.ColumnID == columnID {
			count++
		}
	}
	return count, nil
}

func (f *fakeTaskRepo) UpdateColumn(_ context.Context, id string, newColumnID string, updatedAt time.Time) (taskout.TaskRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	r, ok := f.tasks[id]
	if !ok {
		return taskout.TaskRecord{}, taskout.ErrNotFound
	}
	r.ColumnID = newColumnID
	r.UpdatedAt = updatedAt
	f.tasks[id] = r
	return r, nil
}

func (f *fakeTaskRepo) Delete(_ context.Context, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.tasks[id]; !ok {
		return taskout.ErrNotFound
	}
	delete(f.tasks, id)
	return nil
}

func setupTestBoard() (boardout.BoardRecord, *fakeBoardRepo) {
	limit := 2
	b := boardout.BoardRecord{
		ID:          "board-1",
		Name:        "Sprint 1",
		Description: "Desc",
		Columns: []boardout.ColumnRecord{
			{ID: "col-1", BoardID: "board-1", Name: "Backlog", OrderIndex: 0},
			{ID: "col-2", BoardID: "board-1", Name: "In Progress", OrderIndex: 1, WipLimit: &limit},
			{ID: "col-3", BoardID: "board-1", Name: "Done", OrderIndex: 2},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	bRepo := &fakeBoardRepo{boards: map[string]boardout.BoardRecord{"board-1": b}}
	return b, bRepo
}

func TestCreateTaskService(t *testing.T) {
	_, bRepo := setupTestBoard()
	tRepo := newFakeTaskRepo()
	svc := taskuc.NewCreateTaskService(bRepo, tRepo)

	out, err := svc.Exec(context.Background(), taskin.CreateTaskInput{
		BoardID:  "board-1",
		Column:   "Backlog",
		Title:    "Setup CI",
		Priority: "HIGH",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Title != "Setup CI" || out.ColumnName != "Backlog" || out.Priority != "HIGH" {
		t.Errorf("task fields mismatch: %+v", out)
	}

	// Unknown board
	_, err = svc.Exec(context.Background(), taskin.CreateTaskInput{
		BoardID: "unknown-board",
		Column:  "Backlog",
		Title:   "Task",
	})
	if !errors.Is(err, taskin.ErrBoardNotFound) {
		t.Errorf("expected ErrBoardNotFound, got %v", err)
	}

	// Unknown column
	_, err = svc.Exec(context.Background(), taskin.CreateTaskInput{
		BoardID: "board-1",
		Column:  "NonExistentColumn",
		Title:   "Task",
	})
	if !errors.Is(err, taskin.ErrColumnNotFound) {
		t.Errorf("expected ErrColumnNotFound, got %v", err)
	}
}

func TestMoveTaskService_WipLimitEnforcement(t *testing.T) {
	_, bRepo := setupTestBoard()
	tRepo := newFakeTaskRepo()
	createSvc := taskuc.NewCreateTaskService(bRepo, tRepo)
	moveSvc := taskuc.NewMoveTaskService(bRepo, tRepo)

	// In Progress has wipLimit = 2
	t1, _ := createSvc.Exec(context.Background(), taskin.CreateTaskInput{BoardID: "board-1", Column: "Backlog", Title: "Task 1"})
	t2, _ := createSvc.Exec(context.Background(), taskin.CreateTaskInput{BoardID: "board-1", Column: "Backlog", Title: "Task 2"})
	t3, _ := createSvc.Exec(context.Background(), taskin.CreateTaskInput{BoardID: "board-1", Column: "Backlog", Title: "Task 3"})

	// Move task 1 to In Progress (count becomes 1)
	_, err := moveSvc.Exec(context.Background(), taskin.MoveTaskInput{TaskID: t1.ID, TargetColumnName: "In Progress"})
	if err != nil {
		t.Fatalf("unexpected error moving t1: %v", err)
	}

	// Move task 2 to In Progress (count becomes 2)
	_, err = moveSvc.Exec(context.Background(), taskin.MoveTaskInput{TaskID: t2.ID, TargetColumnName: "In Progress"})
	if err != nil {
		t.Fatalf("unexpected error moving t2: %v", err)
	}

	// Move task 3 to In Progress (should exceed limit of 2)
	_, err = moveSvc.Exec(context.Background(), taskin.MoveTaskInput{TaskID: t3.ID, TargetColumnName: "In Progress"})
	if err == nil {
		t.Fatalf("expected WIP limit error, got nil")
	}

	var wipErr *taskin.WipLimitExceededError
	if !errors.As(err, &wipErr) {
		t.Fatalf("expected *taskin.WipLimitExceededError, got %T: %v", err, err)
	}
	if wipErr.CurrentCount != 2 || wipErr.Limit != 2 {
		t.Errorf("wip error details mismatch: count=%d, limit=%d", wipErr.CurrentCount, wipErr.Limit)
	}
}

func TestDeleteTaskService(t *testing.T) {
	_, bRepo := setupTestBoard()
	tRepo := newFakeTaskRepo()
	createSvc := taskuc.NewCreateTaskService(bRepo, tRepo)
	deleteSvc := taskuc.NewDeleteTaskService(tRepo)

	t1, _ := createSvc.Exec(context.Background(), taskin.CreateTaskInput{BoardID: "board-1", Column: "Backlog", Title: "Task 1"})

	out, err := deleteSvc.Exec(context.Background(), taskin.DeleteTaskInput{TaskID: t1.ID})
	if err != nil || !out.Deleted {
		t.Fatalf("failed to delete task: %v", err)
	}

	// Deleting again should fail with ErrTaskNotFound
	_, err = deleteSvc.Exec(context.Background(), taskin.DeleteTaskInput{TaskID: t1.ID})
	if !errors.Is(err, taskin.ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

