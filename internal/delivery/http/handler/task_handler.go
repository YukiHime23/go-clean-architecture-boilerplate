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

type taskUsecase interface {
	Create(userID uint, input *domain.CreateTaskInput) (*domain.Task, error)
	GetByID(id, userID uint) (*domain.Task, error)
	GetAllByUserID(userID uint, page, limit int) ([]domain.Task, int64, error)
	Update(id, userID uint, input *domain.UpdateTaskInput) (*domain.Task, error)
	Delete(id, userID uint) error
}

type TaskHandler struct {
	taskUsecase taskUsecase
	log         *zap.Logger
}

func NewTaskHandler(uc taskUsecase, log *zap.Logger) *TaskHandler {
	return &TaskHandler{taskUsecase: uc, log: log}
}

func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	input := &domain.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.TaskStatus(req.Status),
		DueDate:     req.DueDate,
	}

	task, err := h.taskUsecase.Create(userID, input)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, "task created", mapToTaskResponse(task))
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	task, usecaseErr := h.taskUsecase.GetByID(id, userID)
	if usecaseErr != nil {
		response.HandleError(c, usecaseErr)
		return
	}
	response.OK(c, "task retrieved", mapToTaskResponse(task))
}

func (h *TaskHandler) GetAll(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tasks, total, err := h.taskUsecase.GetAllByUserID(userID, page, limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	taskResponses := make([]dto.TaskResponse, len(tasks))
	for i, t := range tasks {
		taskResponses[i] = *mapToTaskResponse(&t)
	}

	meta := &response.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}
	response.OKWithMeta(c, "tasks retrieved", taskResponses, meta)
}

func (h *TaskHandler) Update(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	input := &domain.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	}
	if req.Status != nil {
		status := domain.TaskStatus(*req.Status)
		input.Status = &status
	}

	task, usecaseErr := h.taskUsecase.Update(id, userID, input)
	if usecaseErr != nil {
		response.HandleError(c, usecaseErr)
		return
	}
	response.OK(c, "task updated", mapToTaskResponse(task))
}

func (h *TaskHandler) Delete(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}
	if usecaseErr := h.taskUsecase.Delete(id, userID); usecaseErr != nil {
		response.HandleError(c, usecaseErr)
		return
	}
	response.NoContent(c)
}

func mapToTaskResponse(t *domain.Task) *dto.TaskResponse {
	return &dto.TaskResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
