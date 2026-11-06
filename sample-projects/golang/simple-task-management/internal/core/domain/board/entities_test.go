package board_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ports-and-adapters/simple-task-management/internal/core/domain/board"
)

func TestNewBoard_DefaultColumns(t *testing.T) {
	b, err := board.NewBoard("", "Sprint 1", "Release MVP", nil, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if b.Name() != "Sprint 1" {
		t.Errorf("got name %q, want %q", b.Name(), "Sprint 1")
	}
	if len(b.Columns()) != 3 {
		t.Fatalf("got %d columns, want 3", len(b.Columns()))
	}
	if b.Columns()[0].Name() != "Backlog" || b.Columns()[0].WipLimit() != nil {
		t.Errorf("expected Backlog to be unlimited")
	}
	if b.Columns()[1].Name() != "In Progress" || b.Columns()[1].WipLimit() == nil || *b.Columns()[1].WipLimit() != 3 {
		t.Errorf("expected In Progress to have wip limit 3")
	}
	if b.Columns()[2].Name() != "Done" || b.Columns()[2].WipLimit() != nil {
		t.Errorf("expected Done to be unlimited")
	}
}

func TestNewBoard_EmptyName(t *testing.T) {
	_, err := board.NewBoard("", "   ", "", nil, time.Now())
	if !errors.Is(err, board.ErrInvalidBoard) {
		t.Errorf("expected ErrInvalidBoard, got %v", err)
	}
}

func TestNewBoard_DuplicateColumns(t *testing.T) {
	col1, _ := board.NewColumn("", "Review", 0, nil)
	col2, _ := board.NewColumn("", "review", 1, nil)
	_, err := board.NewBoard("", "Sprint", "", []board.Column{col1, col2}, time.Now())
	if !errors.Is(err, board.ErrDuplicateColumnName) {
		t.Errorf("expected ErrDuplicateColumnName, got %v", err)
	}
}

func TestBoard_AddColumn(t *testing.T) {
	b, err := board.NewBoard("", "Sprint", "", nil, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	limit := 2
	col, err := b.AddColumn("", "QA Review", &limit, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if col.Name() != "QA Review" {
		t.Errorf("got name %q, want %q", col.Name(), "QA Review")
	}
	if col.WipLimit() == nil || *col.WipLimit() != 2 {
		t.Errorf("expected wip limit 2")
	}

	// Adding duplicate column should fail
	_, err = b.AddColumn("", "qa review", nil, time.Now())
	if !errors.Is(err, board.ErrDuplicateColumnName) {
		t.Errorf("expected ErrDuplicateColumnName, got %v", err)
	}

	// Invalid WIP limit <= 0
	invalidLimit := 0
	_, err = b.AddColumn("", "Dev", &invalidLimit, time.Now())
	if !errors.Is(err, board.ErrInvalidWipLimit) {
		t.Errorf("expected ErrInvalidWipLimit, got %v", err)
	}
}

func TestBoard_GetColumn(t *testing.T) {
	b, _ := board.NewBoard("", "Sprint", "", nil, time.Now())

	col, err := b.GetColumn("in progress")
	if err != nil {
		t.Fatalf("failed to find column: %v", err)
	}
	if col.Name() != "In Progress" {
		t.Errorf("got column %s, want In Progress", col.Name())
	}

	_, err = b.GetColumn("non-existent")
	if !errors.Is(err, board.ErrColumnNotFound) {
		t.Errorf("expected ErrColumnNotFound, got %v", err)
	}
}

