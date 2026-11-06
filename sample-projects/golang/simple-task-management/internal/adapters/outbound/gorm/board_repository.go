package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	boardout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/board"
	"gorm.io/gorm"
)

var _ boardout.Repository = (*BoardRepository)(nil)

type BoardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Save(ctx context.Context, record boardout.BoardRecord) (boardout.BoardRecord, error) {
	model := toBoardModel(record)

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing BoardModel
		if err := tx.Where("id = ?", model.ID).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(&model).Error
			}
			return err
		}

		if err := tx.Model(&existing).Updates(map[string]any{
			"name":        model.Name,
			"description": model.Description,
			"updated_at":  model.UpdatedAt,
		}).Error; err != nil {
			return err
		}

		// Update or insert columns
		for _, c := range model.Columns {
			c.BoardID = model.ID
			if err := tx.Save(&c).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return boardout.BoardRecord{}, fmt.Errorf("gorm: save board: %w", err)
	}

	return r.GetByID(ctx, record.ID)
}

func (r *BoardRepository) GetByID(ctx context.Context, id string) (boardout.BoardRecord, error) {
	var model BoardModel
	err := r.db.WithContext(ctx).
		Preload("Columns", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Where("id = ?", id).
		First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return boardout.BoardRecord{}, boardout.ErrNotFound
	}
	if err != nil {
		return boardout.BoardRecord{}, fmt.Errorf("gorm: get board %q: %w", id, err)
	}

	return toBoardRecord(model), nil
}

func (r *BoardRepository) ListAll(ctx context.Context) ([]boardout.BoardRecord, error) {
	var models []BoardModel
	err := r.db.WithContext(ctx).
		Preload("Columns", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Order("created_at ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("gorm: list boards: %w", err)
	}

	records := make([]boardout.BoardRecord, len(models))
	for i, m := range models {
		records[i] = toBoardRecord(m)
	}
	return records, nil
}

func (r *BoardRepository) AddColumn(ctx context.Context, boardID string, col boardout.ColumnRecord) (boardout.ColumnRecord, error) {
	now := time.Now().UTC()
	model := ColumnModel{
		ID:         col.ID,
		BoardID:    boardID,
		Name:       col.Name,
		OrderIndex: col.OrderIndex,
		WipLimit:   col.WipLimit,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var b BoardModel
		if err := tx.Where("id = ?", boardID).First(&b).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return boardout.ErrNotFound
			}
			return err
		}

		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		return tx.Model(&b).Update("updated_at", now).Error
	})
	if err != nil {
		if errors.Is(err, boardout.ErrNotFound) {
			return boardout.ColumnRecord{}, boardout.ErrNotFound
		}
		return boardout.ColumnRecord{}, fmt.Errorf("gorm: add column: %w", err)
	}

	return toColumnRecord(model), nil
}

