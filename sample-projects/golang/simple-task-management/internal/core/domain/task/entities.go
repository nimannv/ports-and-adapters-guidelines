package task

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

// Priority defines the urgency level of a task.
type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
	PriorityUrgent Priority = "URGENT"
)

// ParsePriority parses a string into a valid Priority, returning an ErrInvalidTask if unknown.
func ParsePriority(s string) (Priority, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case string(PriorityLow):
		return PriorityLow, nil
	case string(PriorityMedium), "":
		return PriorityMedium, nil
	case string(PriorityHigh):
		return PriorityHigh, nil
	case string(PriorityUrgent):
		return PriorityUrgent, nil
	default:
		return "", fmt.Errorf("%w: invalid priority %q, must be one of LOW, MEDIUM, HIGH, URGENT", ErrInvalidTask, s)
	}
}

// Task represents an item of work tracked on a board.
type Task struct {
	id          string
	boardID     string
	columnID    string
	title       string
	description string
	priority    Priority
	dueDate     *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

// NewTask constructs a new validated Task entity.
func NewTask(
	id, boardID, columnID, title, description string,
	priority Priority,
	dueDate *time.Time,
	now time.Time,
) (*Task, error) {
	trimmedTitle := strings.TrimSpace(title)
	if trimmedTitle == "" {
		return nil, fmt.Errorf("%w: task title cannot be empty", ErrInvalidTask)
	}
	trimmedBoardID := strings.TrimSpace(boardID)
	if trimmedBoardID == "" {
		return nil, fmt.Errorf("%w: board_id cannot be empty", ErrInvalidTask)
	}
	trimmedColumnID := strings.TrimSpace(columnID)
	if trimmedColumnID == "" {
		return nil, fmt.Errorf("%w: column_id cannot be empty", ErrInvalidTask)
	}

	if priority == "" {
		priority = PriorityMedium
	}

	if id == "" {
		id = generateUUID()
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return &Task{
		id:          id,
		boardID:     trimmedBoardID,
		columnID:    trimmedColumnID,
		title:       trimmedTitle,
		description: strings.TrimSpace(description),
		priority:    priority,
		dueDate:     dueDate,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// RehydrateTask rebuilds a Task from trusted persistent storage.
func RehydrateTask(
	id, boardID, columnID, title, description string,
	priority Priority,
	dueDate *time.Time,
	createdAt, updatedAt time.Time,
) *Task {
	return &Task{
		id:          id,
		boardID:     boardID,
		columnID:    columnID,
		title:       title,
		description: description,
		priority:    priority,
		dueDate:     dueDate,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (t *Task) ID() string             { return t.id }
func (t *Task) BoardID() string        { return t.boardID }
func (t *Task) ColumnID() string       { return t.columnID }
func (t *Task) Title() string          { return t.title }
func (t *Task) Description() string    { return t.description }
func (t *Task) Priority() Priority     { return t.priority }
func (t *Task) DueDate() *time.Time    { return t.dueDate }
func (t *Task) CreatedAt() time.Time   { return t.createdAt }
func (t *Task) UpdatedAt() time.Time   { return t.updatedAt }

// MoveTo updates the task's column location, enforcing that columns belong to the same board.
func (t *Task) MoveTo(newColumnID string, targetBoardID string, now time.Time) error {
	if targetBoardID != t.boardID {
		return fmt.Errorf("%w: tasks can only be moved between columns that belong to the same board", ErrInvalidTaskMovement)
	}
	trimmedColumn := strings.TrimSpace(newColumnID)
	if trimmedColumn == "" {
		return fmt.Errorf("%w: target column cannot be empty", ErrInvalidTaskMovement)
	}

	t.columnID = trimmedColumn
	if now.IsZero() {
		now = time.Now().UTC()
	}
	t.updatedAt = now
	return nil
}

func generateUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

