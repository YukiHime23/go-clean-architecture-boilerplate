package repository

import (
	"errors"
	"time"

	"go-clean-architecture-boilerplate/internal/domain"

	"gorm.io/gorm"
)

// taskModel reflects the database schema for the tasks table.
type taskModel struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	UserID      uint   `gorm:"not null;index"`
	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(20);default:'pending'"`
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (taskModel) TableName() string {
	return "tasks"
}

func (m *taskModel) ToDomain() *domain.Task {
	return &domain.Task{
		ID:          m.ID,
		UserID:      m.UserID,
		Title:       m.Title,
		Description: m.Description,
		Status:      domain.TaskStatus(m.Status),
		DueDate:     m.DueDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func fromTaskDomain(t *domain.Task) *taskModel {
	return &taskModel{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	_ = db.AutoMigrate(&taskModel{})
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *domain.Task) error {
	m := fromTaskDomain(task)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	task.ID = m.ID
	return nil
}

func (r *TaskRepository) FindByID(id uint) (*domain.Task, error) {
	var m taskModel
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *TaskRepository) FindAllByUserID(userID uint, page, limit int) ([]domain.Task, int64, error) {
	var models []taskModel
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&taskModel{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	tasks := make([]domain.Task, len(models))
	for i, m := range models {
		tasks[i] = *m.ToDomain()
	}
	return tasks, total, nil
}

func (r *TaskRepository) Update(task *domain.Task) error {
	m := fromTaskDomain(task)
	return r.db.Save(m).Error
}

func (r *TaskRepository) Delete(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&taskModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
