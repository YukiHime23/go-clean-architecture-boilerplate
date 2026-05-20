package apperror

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

var (
	ErrNotFound      = New(http.StatusNotFound, "resource not found")
	ErrUnauthorized  = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden     = New(http.StatusForbidden, "forbidden")
	ErrBadRequest    = New(http.StatusBadRequest, "bad request")
	ErrConflict      = New(http.StatusConflict, "resource already exists")
	ErrInternal      = New(http.StatusInternalServerError, "internal server error")
)
