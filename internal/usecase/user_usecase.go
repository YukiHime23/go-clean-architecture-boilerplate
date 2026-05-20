package usecase

import (
	"errors"

	"go-clean-architecture-boilerplate/internal/domain"

	"go.uber.org/zap"
)

type userStore interface {
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindAll(page, limit int) ([]domain.User, int64, error)
	Update(user *domain.User) error
	Delete(id uint) error
}

type UserUsecase struct {
	userRepo userStore
	log      *zap.Logger
}

func NewUserUsecase(userRepo userStore, log *zap.Logger) *UserUsecase {
	return &UserUsecase{userRepo: userRepo, log: log}
}

func (u *UserUsecase) GetByID(id uint) (*domain.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewNotFoundError("user not found")
		}
		u.log.Error("GetByID: DB error", zap.Error(err))
		return nil, domain.NewInternalError(err)
	}
	return user, nil
}

func (u *UserUsecase) GetAll(page, limit int) ([]domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	users, total, err := u.userRepo.FindAll(page, limit)
	if err != nil {
		u.log.Error("GetAll: DB error", zap.Error(err))
		return nil, 0, domain.NewInternalError(err)
	}
	return users, total, nil
}

func (u *UserUsecase) Update(id uint, input *domain.UpdateUserInput) (*domain.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewNotFoundError("user not found")
		}
		u.log.Error("Update: FindByID error", zap.Error(err))
		return nil, domain.NewInternalError(err)
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		// Check the new email is not taken by another user.
		existing, err := u.userRepo.FindByEmail(input.Email)
		if err == nil && existing.ID != id {
			return nil, domain.NewConflictError("email already in use")
		}
		user.Email = input.Email
	}

	if err := u.userRepo.Update(user); err != nil {
		u.log.Error("Update: save error", zap.Error(err))
		return nil, domain.NewInternalError(err)
	}
	return user, nil
}

func (u *UserUsecase) Delete(id uint) error {
	if err := u.userRepo.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.NewNotFoundError("user not found")
		}
		u.log.Error("Delete: DB error", zap.Error(err))
		return domain.NewInternalError(err)
	}
	return nil
}
