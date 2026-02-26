package domain

// AuthUsecase defines the contract for authentication business logic.
type AuthUsecase interface {
	Register(input *RegisterInput) (*User, error)
	Login(input *LoginInput) (token string, err error)
}

// RegisterInput is a clean struct for registration logic.
type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

// LoginInput is a clean struct for login logic.
type LoginInput struct {
	Email    string
	Password string
}
