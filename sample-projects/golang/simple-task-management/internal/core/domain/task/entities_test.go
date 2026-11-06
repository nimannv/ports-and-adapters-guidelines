package task_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ports-and-adapters/simple-task-management/internal/core/domain/task"
)

func TestNewTask_Valid(t *testing.T) {
	now := time.Now().UTC()
	taskItem, err := task.NewTask("", "b-1", "c-1", "Implement Login", "Add auth", task.PriorityHigh, nil, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if taskItem.Title() != "Implement Login" {
		t.Errorf("got %q, want %q", taskItem.Title(), "Implement Login")
	}
	if taskItem.Priority() != task.PriorityHigh {
		t.Errorf("got %s, want HIGH", taskItem.Priority())
	}
	if taskItem.BoardID() != "b-1" || taskItem.ColumnID() != "c-1" {
		t.Errorf("board/column ID mismatch")
	}
}

func TestNewTask_InvalidFields(t *testing.T) {
	now := time.Now()
	_, err := task.NewTask("", "", "c-1", "Title", "", task.PriorityLow, nil, now)
	if !errors.Is(err, task.ErrInvalidTask) {
		t.Errorf("expected ErrInvalidTask for empty board, got %v", err)
	}

	_, err = task.NewTask("", "b-1", "", "Title", "", task.PriorityLow, nil, now)
	if !errors.Is(err, task.ErrInvalidTask) {
		t.Errorf("expected ErrInvalidTask for empty column, got %v", err)
	}

	_, err = task.NewTask("", "b-1", "c-1", "   ", "", task.PriorityLow, nil, now)
	if !errors.Is(err, task.ErrInvalidTask) {
		t.Errorf("expected ErrInvalidTask for empty title, got %v", err)
	}
}

func TestTask_MoveTo(t *testing.T) {
	now := time.Now()
	taskItem, _ := task.NewTask("t-1", "b-1", "col-1", "Task 1", "", task.PriorityMedium, nil, now)

	// Move to column on same board
	err := taskItem.MoveTo("col-2", "b-1", now.Add(time.Minute))
	if err != nil {
		t.Fatalf("unexpected error moving task: %v", err)
	}
	if taskItem.ColumnID() != "col-2" {
		t.Errorf("expected col-2, got %s", taskItem.ColumnID())
	}

	// Move to different board should fail
	err = taskItem.MoveTo("col-3", "b-different", now)
	if !errors.Is(err, task.ErrInvalidTaskMovement) {
		t.Errorf("expected ErrInvalidTaskMovement, got %v", err)
	}
}

func TestParsePriority(t *testing.T) {
	tests := []struct {
		input   string
		want    task.Priority
		wantErr bool
	}{
		{"low", task.PriorityLow, false},
		{"MEDIUM", task.PriorityMedium, false},
		{"High", task.PriorityHigh, false},
		{"urgent", task.PriorityUrgent, false},
		{"", task.PriorityMedium, false},
		{"UNKNOWN", "", true},
	}

	for _, tt := range tests {
		got, err := task.ParsePriority(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParsePriority(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("ParsePriority(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

