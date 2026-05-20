package handler

import (
	"net/http"

	"go-clean-api/pkg/apperror"
)

func newValidationError(err error) *apperror.AppError {
	return apperror.New(http.StatusUnprocessableEntity, err.Error())
}
