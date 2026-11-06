package task

import "time"

// TaskRecord is the persistence DTO for a task.
type TaskRecord struct {
	ID          string
	BoardID     string
	ColumnID    string
	Title       string
	Description string
	Priority    string
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

