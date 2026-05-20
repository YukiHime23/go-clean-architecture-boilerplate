package domain

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
