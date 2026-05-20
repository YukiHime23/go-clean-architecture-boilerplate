package mysql

import (
	"go-clean-api/internal/domain/entity"
	"go-clean-api/pkg/apperror"

	"gorm.io/gorm"
)

type taskModel struct {
	gorm.Model
	Title       string            `gorm:"size:500;not null"`
	Description string            `gorm:"type:text"`
	Status      entity.TaskStatus `gorm:"type:enum('pending','in_progress','done');default:'pending'"`
	UserID      uint              `gorm:"not null;index"`
}

func (taskModel) TableName() string { return "tasks" }

func toTaskEntity(m *taskModel) *entity.Task {
	return &entity.Task{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Status:      m.Status,
		UserID:      m.UserID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) FindByID(id uint) (*entity.Task, error) {
	var m taskModel
	if err := r.db.First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return toTaskEntity(&m), nil
}

func (r *taskRepository) FindAllByUserID(userID uint) ([]*entity.Task, error) {
	var models []taskModel
	if err := r.db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	tasks := make([]*entity.Task, len(models))
	for i, m := range models {
		m := m
		tasks[i] = toTaskEntity(&m)
	}
	return tasks, nil
}

func (r *taskRepository) Create(task *entity.Task) error {
	m := &taskModel{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		UserID:      task.UserID,
	}
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	task.ID = m.ID
	task.CreatedAt = m.CreatedAt
	task.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *taskRepository) Update(task *entity.Task) error {
	return r.db.Model(&taskModel{}).Where("id = ?", task.ID).
		Updates(map[string]interface{}{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
		}).Error
}

func (r *taskRepository) Delete(id uint) error {
	return r.db.Delete(&taskModel{}, id).Error
}
