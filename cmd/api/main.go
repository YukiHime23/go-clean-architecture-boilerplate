package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-clean-architecture-boilerplate/config"
	deliveryHttp "go-clean-architecture-boilerplate/internal/delivery/http"
	"go-clean-architecture-boilerplate/internal/delivery/http/handler"
	"go-clean-architecture-boilerplate/internal/delivery/http/middleware"
	"go-clean-architecture-boilerplate/internal/repository"
	"go-clean-architecture-boilerplate/internal/usecase"
	"go-clean-architecture-boilerplate/pkg/database"
	"go-clean-architecture-boilerplate/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// --- 1. Load .env file (must be first) ---
	// godotenv.Load is a no-op when .env doesn't exist (Docker/k8s use injected env).
	if err := godotenv.Load(); err != nil {
		// Not fatal — silently continue; env vars may come from the shell or container.
		_ = err
	}

	// --- 2. Load & validate typed config ---
	cfg, err := config.Load()
	if err != nil {
		// Use a temporary std-logger here because Zap hasn't been initialized yet.
		_, _ = os.Stderr.WriteString("FATAL: " + err.Error() + "\n")
		os.Exit(1)
	}

	// --- 3. Set Gin mode from config (before anything uses Gin) ---
	gin.SetMode(cfg.App.GinMode)

	// --- 4. Initialize structured logger ---
	log, err := logger.New(cfg.App)
	if err != nil {
		_, _ = os.Stderr.WriteString("FATAL: failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer log.Sync() //nolint:errcheck

	log.Info("configuration loaded",
		zap.String("gin_mode", cfg.App.GinMode),
		zap.String("port", cfg.App.Port),
		zap.String("db_host", cfg.Database.Host),
	)

	// --- 5. Connect to database ---
	db, err := database.Open(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	log.Info("database connection established")

	// =========================================================================
	// Dependency Injection (manual wiring — innermost layers first)
	// =========================================================================

	// Repositories
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Usecases  (JWT config injected — no os.Getenv inside business logic)
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWT, log)
	userUC := usecase.NewUserUsecase(userRepo, log)
	taskUC := usecase.NewTaskUsecase(taskRepo, log)

	// Handlers
	authHandler := handler.NewAuthHandler(authUC, log)
	userHandler := handler.NewUserHandler(userUC, log)
	taskHandler := handler.NewTaskHandler(taskUC, log)

	// =========================================================================
	// Router  (JWT secret injected — no os.Getenv inside delivery layer)
	// =========================================================================
	router := deliveryHttp.NewRouter(deliveryHttp.RouterConfig{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		TaskHandler: taskHandler,
		JWTSecret:   cfg.JWT.Secret,
	})

	// Attach global Zap request/response logger.
	router.Use(middleware.ZapLogger(log))

	// =========================================================================
	// HTTP Server with Graceful Shutdown
	// =========================================================================
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server gracefully…")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", zap.Error(err))
	}
	log.Info("server stopped")
}
