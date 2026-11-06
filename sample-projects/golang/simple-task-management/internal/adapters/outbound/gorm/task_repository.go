package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	taskout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/task"
	"gorm.io/gorm"
)

var _ taskout.Repository = (*TaskRepository)(nil)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Save(ctx context.Context, record taskout.TaskRecord) (taskout.TaskRecord, error) {
	model := toTaskModel(record)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return taskout.TaskRecord{}, fmt.Errorf("gorm: save task: %w", err)
	}
	return toTaskRecord(model), nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (taskout.TaskRecord, error) {
	var model TaskModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return taskout.TaskRecord{}, taskout.ErrNotFound
	}
	if err != nil {
		return taskout.TaskRecord{}, fmt.Errorf("gorm: get task %q: %w", id, err)
	}
	return toTaskRecord(model), nil
}

func (r *TaskRepository) ListByBoard(ctx context.Context, boardID string, columnID string) ([]taskout.TaskRecord, error) {
	var models []TaskModel
	q := r.db.WithContext(ctx).Where("board_id = ?", boardID)
	if columnID != "" {
		q = q.Where("column_id = ?", columnID)
	}

	if err := q.Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("gorm: list tasks on board %q: %w", boardID, err)
	}

	records := make([]taskout.TaskRecord, len(models))
	for i, m := range models {
		records[i] = toTaskRecord(m)
	}
	return records, nil
}

func (r *TaskRepository) CountByColumn(ctx context.Context, columnID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&TaskModel{}).Where("column_id = ?", columnID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("gorm: count tasks in column %q: %w", columnID, err)
	}
	return int(count), nil
}

func (r *TaskRepository) UpdateColumn(ctx context.Context, id string, newColumnID string, updatedAt time.Time) (taskout.TaskRecord, error) {
	var model TaskModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return taskout.TaskRecord{}, taskout.ErrNotFound
	}
	if err != nil {
		return taskout.TaskRecord{}, fmt.Errorf("gorm: find task for update: %w", err)
	}

	model.ColumnID = newColumnID
	model.UpdatedAt = updatedAt
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return taskout.TaskRecord{}, fmt.Errorf("gorm: update task column: %w", err)
	}

	return toTaskRecord(model), nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&TaskModel{})
	if result.Error != nil {
		return fmt.Errorf("gorm: delete task %q: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return taskout.ErrNotFound
	}
	return nil
}

