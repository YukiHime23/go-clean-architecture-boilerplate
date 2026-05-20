package task

import (
	"go-clean-api/internal/domain/entity"
	"go-clean-api/pkg/apperror"
)

type taskFinder interface {
	FindByID(id uint) (*entity.Task, error)
}

type taskLister interface {
	FindAllByUserID(userID uint) ([]*entity.Task, error)
}

type taskCreatorUpdater interface {
	Create(task *entity.Task) error
	Update(task *entity.Task) error
}

type taskDeleter interface {
	Delete(id uint) error
}

type taskRepo interface {
	taskFinder
	taskLister
	taskCreatorUpdater
	taskDeleter
}

type CreateInput struct {
	Title       string            `json:"title" binding:"required,min=1"`
	Description string            `json:"description"`
	Status      entity.TaskStatus `json:"status" binding:"omitempty,oneof=pending in_progress done"`
}

type UpdateInput struct {
	Title       string            `json:"title" binding:"omitempty,min=1"`
	Description string            `json:"description"`
	Status      entity.TaskStatus `json:"status" binding:"omitempty,oneof=pending in_progress done"`
}

type UseCase struct {
	taskRepo taskRepo
}

func New(taskRepo taskRepo) *UseCase {
	return &UseCase{taskRepo: taskRepo}
}

func (u *UseCase) Create(userID uint, input CreateInput) (*entity.Task, error) {
	status := input.Status
	if status == "" {
		status = entity.TaskStatusPending
	}
	task := &entity.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      status,
		UserID:      userID,
	}
	if err := u.taskRepo.Create(task); err != nil {
		return nil, apperror.ErrInternal
	}
	return task, nil
}

func (u *UseCase) GetAll(userID uint) ([]*entity.Task, error) {
	return u.taskRepo.FindAllByUserID(userID)
}

func (u *UseCase) GetByID(userID, taskID uint) (*entity.Task, error) {
	task, err := u.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, err
	}
	if task.UserID != userID {
		return nil, apperror.ErrForbidden
	}
	return task, nil
}

func (u *UseCase) Update(userID, taskID uint, input UpdateInput) (*entity.Task, error) {
	task, err := u.GetByID(userID, taskID)
	if err != nil {
		return nil, err
	}
	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Status != "" {
		task.Status = input.Status
	}
	if err := u.taskRepo.Update(task); err != nil {
		return nil, apperror.ErrInternal
	}
	return task, nil
}

func (u *UseCase) Delete(userID, taskID uint) error {
	if _, err := u.GetByID(userID, taskID); err != nil {
		return err
	}
	return u.taskRepo.Delete(taskID)
}
