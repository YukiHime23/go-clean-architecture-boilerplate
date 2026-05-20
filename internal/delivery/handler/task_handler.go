package handler

import (
	"net/http"
	"strconv"

	"go-clean-api/internal/delivery/middleware"
	"go-clean-api/internal/domain/entity"
	"go-clean-api/internal/usecase/task"
	"go-clean-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type taskService interface {
	Create(userID uint, input task.CreateInput) (*entity.Task, error)
	GetAll(userID uint) ([]*entity.Task, error)
	GetByID(userID, taskID uint) (*entity.Task, error)
	Update(userID, taskID uint, input task.UpdateInput) (*entity.Task, error)
	Delete(userID, taskID uint) error
}

type TaskHandler struct {
	svc taskService
}

func NewTaskHandler(svc taskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uint)
	var input task.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, newValidationError(err))
		return
	}
	t, err := h.svc.Create(userID, input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, t)
}

func (h *TaskHandler) GetAll(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uint)
	tasks, err := h.svc.GetAll(userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, tasks)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uint)
	taskID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	t, err := h.svc.GetByID(userID, taskID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, t)
}

func (h *TaskHandler) Update(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uint)
	taskID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var input task.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, newValidationError(err))
		return
	}
	t, err := h.svc.Update(userID, taskID, input)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, t)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uint)
	taskID, err := parseUintParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.Delete(userID, taskID); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

func parseUintParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 64)
	if err != nil {
		return 0, newValidationError(err)
	}
	return uint(val), nil
}
