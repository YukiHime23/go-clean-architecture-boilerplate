package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-clean-api/pkg/apperror"
)

type envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, code int, data interface{}) {
	c.JSON(code, envelope{Success: true, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, envelope{Success: true, Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if e, ok := err.(*apperror.AppError); ok {
		appErr = e
	} else {
		appErr = apperror.ErrInternal
	}
	c.JSON(appErr.Code, envelope{Success: false, Message: appErr.Message})
}
