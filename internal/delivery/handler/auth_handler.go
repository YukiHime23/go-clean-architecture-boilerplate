package handler

import (
	"net/http"

	"go-clean-api/internal/usecase/auth"
	"go-clean-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type authService interface {
	Register(input auth.RegisterInput) (*auth.Output, error)
	Login(input auth.LoginInput) (*auth.Output, error)
}

type AuthHandler struct {
	svc authService
}

func NewAuthHandler(svc authService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input auth.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, newValidationError(err))
		return
	}
	out, err := h.svc.Register(input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, out)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input auth.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, newValidationError(err))
		return
	}
	out, err := h.svc.Login(input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, out)
}
