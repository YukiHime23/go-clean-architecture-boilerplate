// Package response provides centralized HTTP response helpers for the API.
package response

import (
	"errors"
	"net/http"

	"go-clean-architecture-boilerplate/internal/domain"

	"github.com/gin-gonic/gin"
)

// Meta holds pagination metadata.
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// envelope is the standard envelope for all API responses.
type envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// OK sends a 200 response.
func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, envelope{Success: true, Message: message, Data: data})
}

// Created sends a 201 response.
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, envelope{Success: true, Message: message, Data: data})
}

// OKWithMeta sends a 200 response with pagination metadata.
func OKWithMeta(c *gin.Context, message string, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, envelope{Success: true, Message: message, Data: data, Meta: meta})
}

// Unauthorized sends a 401 response.
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, envelope{Success: false, Error: msg})
}

// Forbidden sends a 403 response.
func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, envelope{Success: false, Error: msg})
}

// NoContent sends a 204 response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// HandleError inspects the error type and sends the appropriate HTTP response.
// It handles *domain.AppError natively; all other errors become 500.
func HandleError(c *gin.Context, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, envelope{Success: false, Error: appErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, envelope{Success: false, Error: "an unexpected error occurred"})
}
