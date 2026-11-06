//go:build integration

package integration_test

import (
	"context"
	"errors"
	"testing"
	"time"

	gormadapter "github.com/ports-and-adapters/simple-task-management/internal/adapters/outbound/gorm"
	boardout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/board"
	taskout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/task"
)

func TestGORMRepositories_SQLiteRoundTrip(t *testing.T) {
	ctx := context.Background()
	db, err := gormadapter.OpenDB("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory db: %v", err)
	}

	boardRepo := gormadapter.NewBoardRepository(db)
	taskRepo := gormadapter.NewTaskRepository(db)

	limit := 3
	boardRec := boardout.BoardRecord{
		ID:          "b-test-1",
		Name:        "Sprint 1",
		Description: "Release MVP",
		Columns: []boardout.ColumnRecord{
			{ID: "c-1", BoardID: "b-test-1", Name: "Backlog", OrderIndex: 0},
			{ID: "c-2", BoardID: "b-test-1", Name: "In Progress", OrderIndex: 1, WipLimit: &limit},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	savedBoard, err := boardRepo.Save(ctx, boardRec)
	if err != nil {
		t.Fatalf("failed to save board: %v", err)
	}
	if savedBoard.Name != "Sprint 1" || len(savedBoard.Columns) != 2 {
		t.Fatalf("unexpected saved board: %+v", savedBoard)
	}

	// Add column
	newCol, err := boardRepo.AddColumn(ctx, "b-test-1", boardout.ColumnRecord{
		ID:         "c-3",
		Name:       "Done",
		OrderIndex: 2,
	})
	if err != nil {
		t.Fatalf("failed to add column: %v", err)
	}
	if newCol.Name != "Done" {
		t.Errorf("got col %s, want Done", newCol.Name)
	}

	// Save task
	taskRec := taskout.TaskRecord{
		ID:          "t-test-1",
		BoardID:     "b-test-1",
		ColumnID:    "c-1",
		Title:       "Task Alpha",
		Description: "Desc",
		Priority:    "HIGH",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	savedTask, err := taskRepo.Save(ctx, taskRec)
	if err != nil {
		t.Fatalf("failed to save task: %v", err)
	}
	if savedTask.Title != "Task Alpha" {
		t.Errorf("got %s, want Task Alpha", savedTask.Title)
	}

	// Count by column
	count, err := taskRepo.CountByColumn(ctx, "c-1")
	if err != nil || count != 1 {
		t.Errorf("expected count 1, got %d (err: %v)", count, err)
	}

	// Update column
	updatedTask, err := taskRepo.UpdateColumn(ctx, "t-test-1", "c-2", time.Now().UTC())
	if err != nil || updatedTask.ColumnID != "c-2" {
		t.Fatalf("failed to update column: %v", err)
	}

	// Delete task
	err = taskRepo.Delete(ctx, "t-test-1")
	if err != nil {
		t.Fatalf("failed to delete task: %v", err)
	}

	_, err = taskRepo.GetByID(ctx, "t-test-1")
	if !errors.Is(err, taskout.ErrNotFound) {
		t.Errorf("expected ErrNotFound after deletion, got %v", err)
	}
}

