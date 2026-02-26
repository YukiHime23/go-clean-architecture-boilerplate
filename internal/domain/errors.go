package domain

import (
	"errors"
	"net/http"
)

// AppError is a structured application error that carries an HTTP status code
// alongside a developer-facing and a user-facing message.
type AppError struct {
	Code    int    // HTTP status code
	Message string // User-facing message
	Err     error  // Underlying error (for logging)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// --- Sentinel Errors ---
// These are used to distinguish well-known error kinds without exposing HTTP details
// to the inner layers. Delivery layer maps these to AppError.

var (
	ErrNotFound           = errors.New("resource not found")
	ErrConflict           = errors.New("resource already exists")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrBadRequest         = errors.New("bad request")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalServer     = errors.New("internal server error")
)

// --- Constructor helpers ---

func NewNotFoundError(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg, Err: ErrNotFound}
}

func NewConflictError(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg, Err: ErrConflict}
}

func NewUnauthorizedError(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg, Err: ErrUnauthorized}
}

func NewForbiddenError(msg string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: msg, Err: ErrForbidden}
}

func NewBadRequestError(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg, Err: ErrBadRequest}
}

func NewInternalError(err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: "an unexpected error occurred", Err: err}
}
