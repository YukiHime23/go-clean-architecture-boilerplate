package dto

import (
	"time"
)

// TaskResponse is the public representation of a task.
type TaskResponse struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTaskRequest contains the fields required to create a new task.
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=3,max=255"`
	Description string     `json:"description" binding:"omitempty,max=1000"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress done"`
	DueDate     *time.Time `json:"due_date" binding:"omitempty"`
}

// UpdateTaskRequest contains the optional fields for updating a task.
type UpdateTaskRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=3,max=255"`
	Description *string    `json:"description" binding:"omitempty,max=1000"`
	Status      *string    `json:"status" binding:"omitempty,oneof=pending in_progress done"`
	DueDate     *time.Time `json:"due_date" binding:"omitempty"`
}
