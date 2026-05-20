package usecase

import (
	"errors"

	"go-clean-architecture-boilerplate/internal/domain"

	"go.uber.org/zap"
)

type taskStore interface {
	Create(task *domain.Task) error
	FindByID(id uint) (*domain.Task, error)
	FindAllByUserID(userID uint, page, limit int) ([]domain.Task, int64, error)
	Update(task *domain.Task) error
	Delete(id, userID uint) error
}

type TaskUsecase struct {
	taskRepo taskStore
	log      *zap.Logger
}

func NewTaskUsecase(taskRepo taskStore, log *zap.Logger) *TaskUsecase {
	return &TaskUsecase{taskRepo: taskRepo, log: log}
}

func (u *TaskUsecase) Create(userID uint, input *domain.CreateTaskInput) (*domain.Task, error) {
	task := &domain.Task{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
	}
	if input.Status != "" {
		task.Status = input.Status
	} else {
		task.Status = domain.TaskStatusPending
	}

	if err := u.taskRepo.Create(task); err != nil {
		u.log.Error("Create task: DB error", zap.Error(err))
		return nil, domain.NewInternalError(err)
	}
	return task, nil
}

func (u *TaskUsecase) GetByID(id, userID uint) (*domain.Task, error) {
	task, err := u.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewNotFoundError("task not found")
		}
		u.log.Error("GetByID: DB error", zap.Error(err))
		return nil, domain.NewInternalError(err)
	}
	// Ownership guard — users can only see their own tasks.
	if task.UserID != userID {
		return nil, domain.NewForbiddenError("you do not have permission to view this task")
	}
	return task, nil
}

func (u *TaskUsecase) GetAllByUserID(userID uint, page, limit int) ([]domain.Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	tasks, total, err := u.taskRepo.FindAllByUserID(userID, page, limit)
	if err != nil {
		u.log.Error("GetAllByUserID: DB error", zap.Error(err))
		return nil, 0, domain.NewInternalError(err)
	}
	return tasks, total, nil
}

func (u *TaskUsecase) Update(id, userID uint, input *domain.UpdateTaskInput) (*domain.Task, error) {
	task, err := u.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewNotFoundError("task not found")
		}
		u.log.Error("Update: FindByID error", zap.Error(err))
		return nil, domain.NewInternalError(err)
	}
	if task.UserID != userID {
		return nil, domain.NewForbiddenError("you do not have permission to update this task")
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Status != nil {
		task.Status = *input.Status
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}

	if err := u.taskRepo.Update(task); err != nil {
		u.log.Error("Update: save error", zap.Error(err))
		return nil, domain.NewInternalError(err)
	}
	return task, nil
}

func (u *TaskUsecase) Delete(id, userID uint) error {
	if err := u.taskRepo.Delete(id, userID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.NewNotFoundError("task not found or you don't have permission")
		}
		u.log.Error("Delete: DB error", zap.Error(err))
		return domain.NewInternalError(err)
	}
	return nil
}
