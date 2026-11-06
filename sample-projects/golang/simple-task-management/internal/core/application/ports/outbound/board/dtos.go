package board

import "time"

// ColumnRecord is the persistence DTO for a board's column.
type ColumnRecord struct {
	ID         string
	BoardID    string
	Name       string
	OrderIndex int
	WipLimit   *int
}

// BoardRecord is the persistence DTO for a board and its columns.
type BoardRecord struct {
	ID          string
	Name        string
	Description string
	Columns     []ColumnRecord
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

