package router

import (
	"github.com/gin-gonic/gin"
	"go-clean-api/internal/delivery/handler"
	"go-clean-api/internal/delivery/middleware"
)

type Handlers struct {
	Auth *handler.AuthHandler
	User *handler.UserHandler
	Task *handler.TaskHandler
}

func Setup(r *gin.Engine, h Handlers, jwtSecret string) {
	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}

	authorized := api.Group("")
	authorized.Use(middleware.Auth(jwtSecret))
	{
		users := authorized.Group("/users")
		{
			users.GET("/me", h.User.GetMe)
			users.PUT("/me", h.User.UpdateMe)
			users.GET("", h.User.GetAll)
		}

		tasks := authorized.Group("/tasks")
		{
			tasks.POST("", h.Task.Create)
			tasks.GET("", h.Task.GetAll)
			tasks.GET("/:id", h.Task.GetByID)
			tasks.PUT("/:id", h.Task.Update)
			tasks.DELETE("/:id", h.Task.Delete)
		}
	}
}
