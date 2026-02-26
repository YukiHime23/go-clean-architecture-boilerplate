package domain

import "time"

// User represents the core user entity.
// It is tag-free to remain agnostic of external frameworks (DB/JSON/API).
type User struct {
	ID        uint
	Name      string
	Email     string
	Password  string // Hashed password
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRepository defines the contract for user data persistence.
type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll(page, limit int) ([]User, int64, error)
	Update(user *User) error
	Delete(id uint) error
}

// UserUsecase defines the contract for user business logic.
type UserUsecase interface {
	GetByID(id uint) (*User, error)
	GetAll(page, limit int) ([]User, int64, error)
	Update(id uint, input *UpdateUserInput) (*User, error)
	Delete(id uint) error
}

// UpdateUserInput is a clean Go struct for business logic input.
type UpdateUserInput struct {
	Name  string
	Email string
}
