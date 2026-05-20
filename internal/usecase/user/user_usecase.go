package user

import (
	"go-clean-api/internal/domain/entity"
)

type userFinder interface {
	FindByID(id uint) (*entity.User, error)
}

type userLister interface {
	FindAll() ([]*entity.User, error)
}

type userUpdater interface {
	Update(user *entity.User) error
}

type userRepo interface {
	userFinder
	userLister
	userUpdater
}

type UpdateInput struct {
	Name  string `json:"name" binding:"omitempty,min=2"`
	Email string `json:"email" binding:"omitempty,email"`
}

type UseCase struct {
	userRepo userRepo
}

func New(userRepo userRepo) *UseCase {
	return &UseCase{userRepo: userRepo}
}

func (u *UseCase) GetMe(id uint) (*entity.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *UseCase) UpdateMe(id uint, input UpdateInput) (*entity.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UseCase) GetAll() ([]*entity.User, error) {
	return u.userRepo.FindAll()
}
