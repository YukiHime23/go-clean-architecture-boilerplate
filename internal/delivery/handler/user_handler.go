package handler

import (
	"net/http"

	"go-clean-api/internal/delivery/middleware"
	"go-clean-api/internal/domain/entity"
	"go-clean-api/internal/usecase/user"
	"go-clean-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type userService interface {
	GetMe(id uint) (*entity.User, error)
	UpdateMe(id uint, input user.UpdateInput) (*entity.User, error)
	GetAll() ([]*entity.User, error)
}

type UserHandler struct {
	svc userService
}

func NewUserHandler(svc userService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uint)
	u, err := h.svc.GetMe(userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, u)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uint)
	var input user.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, newValidationError(err))
		return
	}
	u, err := h.svc.UpdateMe(userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, u)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.svc.GetAll()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, users)
}
