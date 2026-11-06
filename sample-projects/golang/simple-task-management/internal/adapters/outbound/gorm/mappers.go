package gorm

import (
	boardout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/board"
	taskout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/task"
)

func toBoardRecord(m BoardModel) boardout.BoardRecord {
	cols := make([]boardout.ColumnRecord, len(m.Columns))
	for i, c := range m.Columns {
		cols[i] = toColumnRecord(c)
	}
	return boardout.BoardRecord{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Columns:     cols,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toColumnRecord(c ColumnModel) boardout.ColumnRecord {
	return boardout.ColumnRecord{
		ID:         c.ID,
		BoardID:    c.BoardID,
		Name:       c.Name,
		OrderIndex: c.OrderIndex,
		WipLimit:   c.WipLimit,
	}
}

func toBoardModel(r boardout.BoardRecord) BoardModel {
	cols := make([]ColumnModel, len(r.Columns))
	for i, c := range r.Columns {
		cols[i] = ColumnModel{
			ID:         c.ID,
			BoardID:    r.ID,
			Name:       c.Name,
			OrderIndex: c.OrderIndex,
			WipLimit:   c.WipLimit,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		}
	}
	return BoardModel{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Columns:     cols,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toTaskRecord(m TaskModel) taskout.TaskRecord {
	return taskout.TaskRecord{
		ID:          m.ID,
		BoardID:     m.BoardID,
		ColumnID:    m.ColumnID,
		Title:       m.Title,
		Description: m.Description,
		Priority:    m.Priority,
		DueDate:     m.DueDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toTaskModel(r taskout.TaskRecord) TaskModel {
	return TaskModel{
		ID:          r.ID,
		BoardID:     r.BoardID,
		ColumnID:    r.ColumnID,
		Title:       r.Title,
		Description: r.Description,
		Priority:    r.Priority,
		DueDate:     r.DueDate,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

