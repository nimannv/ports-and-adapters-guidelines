package board

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

// Column represents an ordered workflow stage on a board with an optional WIP limit.
type Column struct {
	id         string
	name       string
	orderIndex int
	wipLimit   *int
}

// NewColumn creates and validates a Column entity.
func NewColumn(id, name string, orderIndex int, wipLimit *int) (Column, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return Column{}, fmt.Errorf("%w: column name cannot be empty", ErrInvalidBoard)
	}
	if wipLimit != nil && *wipLimit <= 0 {
		return Column{}, fmt.Errorf("%w: WIP limit must be greater than 0, got %d", ErrInvalidWipLimit, *wipLimit)
	}
	if id == "" {
		id = generateUUID()
	}
	return Column{
		id:         id,
		name:       trimmed,
		orderIndex: orderIndex,
		wipLimit:   wipLimit,
	}, nil
}

// RehydrateColumn reconstructs a Column from persistent storage.
func RehydrateColumn(id, name string, orderIndex int, wipLimit *int) Column {
	return Column{
		id:         id,
		name:       name,
		orderIndex: orderIndex,
		wipLimit:   wipLimit,
	}
}

func (c Column) ID() string         { return c.id }
func (c Column) Name() string       { return c.name }
func (c Column) OrderIndex() int    { return c.orderIndex }
func (c Column) WipLimit() *int     { return c.wipLimit }

// IsFull checks if adding a task would exceed this column's WIP limit.
func (c Column) IsFull(currentTaskCount int) bool {
	if c.wipLimit == nil {
		return false
	}
	return currentTaskCount >= *c.wipLimit
}

// Board represents a project or workspace containing workflow columns.
type Board struct {
	id          string
	name        string
	description string
	columns     []Column
	createdAt   time.Time
	updatedAt   time.Time
}

// NewBoard constructs a new Board aggregate and enforces business invariants.
func NewBoard(id, name, description string, customColumns []Column, now time.Time) (*Board, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, fmt.Errorf("%w: board name cannot be empty", ErrInvalidBoard)
	}

	var cols []Column
	if len(customColumns) > 0 {
		cols = make([]Column, len(customColumns))
		copy(cols, customColumns)
	} else {
		// Default columns per specification: Backlog (unlimited), In Progress (limit: 3), Done (unlimited)
		defaultLimit := 3
		c1, _ := NewColumn("", "Backlog", 0, nil)
		c2, _ := NewColumn("", "In Progress", 1, &defaultLimit)
		c3, _ := NewColumn("", "Done", 2, nil)
		cols = []Column{c1, c2, c3}
	}

	if len(cols) == 0 {
		return nil, fmt.Errorf("%w: a board must contain at least one column upon creation", ErrInvalidBoard)
	}

	if err := validateUniqueColumns(cols); err != nil {
		return nil, err
	}

	if id == "" {
		id = generateUUID()
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return &Board{
		id:          id,
		name:        trimmedName,
		description: strings.TrimSpace(description),
		columns:     cols,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// RehydrateBoard rebuilds a Board from trusted persistence data.
func RehydrateBoard(id, name, description string, columns []Column, createdAt, updatedAt time.Time) *Board {
	return &Board{
		id:          id,
		name:        name,
		description: description,
		columns:     columns,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (b *Board) ID() string             { return b.id }
func (b *Board) Name() string           { return b.name }
func (b *Board) Description() string    { return b.description }
func (b *Board) Columns() []Column      { return b.columns }
func (b *Board) CreatedAt() time.Time   { return b.createdAt }
func (b *Board) UpdatedAt() time.Time   { return b.updatedAt }

// AddColumn appends a new column to the board, validating name uniqueness and WIP limits.
func (b *Board) AddColumn(id, name string, wipLimit *int, now time.Time) (Column, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return Column{}, fmt.Errorf("%w: column name cannot be empty", ErrInvalidBoard)
	}

	for _, c := range b.columns {
		if strings.EqualFold(c.name, trimmed) {
			return Column{}, fmt.Errorf("%w: column %q already exists on board %s", ErrDuplicateColumnName, trimmed, b.id)
		}
	}

	nextOrder := len(b.columns)
	col, err := NewColumn(id, trimmed, nextOrder, wipLimit)
	if err != nil {
		return Column{}, err
	}

	b.columns = append(b.columns, col)
	if now.IsZero() {
		now = time.Now().UTC()
	}
	b.updatedAt = now
	return col, nil
}

// GetColumn finds a column by either its UUID or case-insensitive name.
func (b *Board) GetColumn(ref string) (Column, error) {
	refTrimmed := strings.TrimSpace(ref)
	for _, c := range b.columns {
		if c.id == refTrimmed || strings.EqualFold(c.name, refTrimmed) {
			return c, nil
		}
	}
	return Column{}, fmt.Errorf("%w: %q on board %s", ErrColumnNotFound, ref, b.id)
}

// HasColumn checks if a column with the given ID exists on this board.
func (b *Board) HasColumn(columnID string) bool {
	for _, c := range b.columns {
		if c.id == columnID {
			return true
		}
	}
	return false
}

func validateUniqueColumns(cols []Column) error {
	seen := make(map[string]struct{}, len(cols))
	for _, c := range cols {
		normalized := strings.ToLower(strings.TrimSpace(c.name))
		if _, exists := seen[normalized]; exists {
			return fmt.Errorf("%w: column %q", ErrDuplicateColumnName, c.name)
		}
		seen[normalized] = struct{}{}
	}
	return nil
}

func generateUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

