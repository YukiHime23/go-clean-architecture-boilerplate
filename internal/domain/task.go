package domain

import "time"

// TaskStatus represents the valid states a task can be in.
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

// Task represents the core task entity.
// It is tag-free to remain agnostic of external frameworks.
type Task struct {
	ID          uint
	UserID      uint
	Title       string
	Description string
	Status      TaskStatus
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TaskRepository defines the contract for task data persistence.
type TaskRepository interface {
	Create(task *Task) error
	FindByID(id uint) (*Task, error)
	FindAllByUserID(userID uint, page, limit int) ([]Task, int64, error)
	Update(task *Task) error
	Delete(id, userID uint) error
}

// TaskUsecase defines the contract for task business logic.
type TaskUsecase interface {
	Create(userID uint, input *CreateTaskInput) (*Task, error)
	GetByID(id, userID uint) (*Task, error)
	GetAllByUserID(userID uint, page, limit int) ([]Task, int64, error)
	Update(id, userID uint, input *UpdateTaskInput) (*Task, error)
	Delete(id, userID uint) error
}

// CreateTaskInput is a clean struct for creating a task.
type CreateTaskInput struct {
	Title       string
	Description string
	Status      TaskStatus
	DueDate     *time.Time
}

// UpdateTaskInput is a clean struct for updating a task.
type UpdateTaskInput struct {
	Title       *string
	Description *string
	Status      *TaskStatus
	DueDate     *time.Time
}
