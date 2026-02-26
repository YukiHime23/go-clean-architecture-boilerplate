package http

import (
	"go-clean-architecture-boilerplate/internal/delivery/http/handler"
	"go-clean-architecture-boilerplate/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

// RouterConfig holds all handlers and settings needed to register routes.
type RouterConfig struct {
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
	TaskHandler *handler.TaskHandler
	JWTSecret   string
}

// NewRouter sets up and returns a configured Gin engine.
func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		// --- Public routes ---
		auth := api.Group("/auth")
		{
			auth.POST("/register", cfg.AuthHandler.Register)
			auth.POST("/login", cfg.AuthHandler.Login)
		}

		// --- Protected routes (require valid JWT) ---
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			// User management
			users := protected.Group("/users")
			{
				users.GET("/me", cfg.UserHandler.GetMe)
				users.PUT("/me", cfg.UserHandler.UpdateMe)
				users.GET("", cfg.UserHandler.GetAll)
				users.GET("/:id", cfg.UserHandler.GetByID)
				users.DELETE("/:id", cfg.UserHandler.Delete)
			}

			// Task management (always scoped to the authenticated user)
			tasks := protected.Group("/tasks")
			{
				tasks.POST("", cfg.TaskHandler.Create)
				tasks.GET("", cfg.TaskHandler.GetAll)
				tasks.GET("/:id", cfg.TaskHandler.GetByID)
				tasks.PUT("/:id", cfg.TaskHandler.Update)
				tasks.DELETE("/:id", cfg.TaskHandler.Delete)
			}
		}
	}

	return r
}
