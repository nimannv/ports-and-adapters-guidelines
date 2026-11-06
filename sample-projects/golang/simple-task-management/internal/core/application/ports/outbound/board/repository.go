package board

import "context"

// Repository defines the outbound storage capability for boards and columns.
type Repository interface {
	Save(ctx context.Context, record BoardRecord) (BoardRecord, error)
	GetByID(ctx context.Context, id string) (BoardRecord, error)
	ListAll(ctx context.Context) ([]BoardRecord, error)
	AddColumn(ctx context.Context, boardID string, column ColumnRecord) (ColumnRecord, error)
}

