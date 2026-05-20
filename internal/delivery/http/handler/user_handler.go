package handler

import (
	"math"
	"net/http"
	"strconv"

	"go-clean-architecture-boilerplate/internal/delivery/http/dto"
	"go-clean-architecture-boilerplate/internal/delivery/http/middleware"
	"go-clean-architecture-boilerplate/internal/delivery/http/response"
	"go-clean-architecture-boilerplate/internal/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type userUsecase interface {
	GetByID(id uint) (*domain.User, error)
	GetAll(page, limit int) ([]domain.User, int64, error)
	Update(id uint, input *domain.UpdateUserInput) (*domain.User, error)
	Delete(id uint) error
}

type UserHandler struct {
	userUsecase userUsecase
	log         *zap.Logger
}

func NewUserHandler(uc userUsecase, log *zap.Logger) *UserHandler {
	return &UserHandler{userUsecase: uc, log: log}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	user, err := h.userUsecase.GetByID(userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, "profile retrieved", mapToUserResponse(user))
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	user, usecaseErr := h.userUsecase.GetByID(id)
	if usecaseErr != nil {
		response.HandleError(c, usecaseErr)
		return
	}
	response.OK(c, "user retrieved", mapToUserResponse(user))
}

func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.userUsecase.GetAll(page, limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = *mapToUserResponse(&u)
	}

	meta := &response.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}
	response.OKWithMeta(c, "users retrieved", userResponses, meta)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	input := &domain.UpdateUserInput{
		Name:  req.Name,
		Email: req.Email,
	}

	user, err := h.userUsecase.Update(userID, input)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, "profile updated", mapToUserResponse(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	if usecaseErr := h.userUsecase.Delete(id); usecaseErr != nil {
		response.HandleError(c, usecaseErr)
		return
	}
	response.NoContent(c)
}

func mapToUserResponse(u *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func parseUintParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id parameter"})
		return 0, err
	}
	return uint(id), nil
}
