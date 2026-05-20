package handler

import (
	"net/http"

	"go-clean-architecture-boilerplate/internal/delivery/http/dto"
	"go-clean-architecture-boilerplate/internal/delivery/http/response"
	"go-clean-architecture-boilerplate/internal/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type authUsecase interface {
	Register(input *domain.RegisterInput) (*domain.User, error)
	Login(input *domain.LoginInput) (string, error)
}

type AuthHandler struct {
	authUsecase authUsecase
	log         *zap.Logger
}

func NewAuthHandler(uc authUsecase, log *zap.Logger) *AuthHandler {
	return &AuthHandler{authUsecase: uc, log: log}
}

// Register godoc
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Map DTO to Domain Input
	input := &domain.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.authUsecase.Register(input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// Map Domain to Response DTO
	resp := &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response.Created(c, "user registered successfully", resp)
}

// Login godoc
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Map DTO to Domain Input
	input := &domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := h.authUsecase.Login(input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, "login successful", dto.LoginResponse{Token: token})
}
