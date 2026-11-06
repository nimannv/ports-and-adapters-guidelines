package task

import (
	"context"
	"time"
)

// Repository defines the outbound persistence operations for tasks.
type Repository interface {
	Save(ctx context.Context, record TaskRecord) (TaskRecord, error)
	GetByID(ctx context.Context, id string) (TaskRecord, error)
	ListByBoard(ctx context.Context, boardID string, columnID string) ([]TaskRecord, error)
	CountByColumn(ctx context.Context, columnID string) (int, error)
	UpdateColumn(ctx context.Context, id string, newColumnID string, updatedAt time.Time) (TaskRecord, error)
	Delete(ctx context.Context, id string) error
}

