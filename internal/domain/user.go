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

// UpdateUserInput is a clean Go struct for business logic input.
type UpdateUserInput struct {
	Name  string
	Email string
}
