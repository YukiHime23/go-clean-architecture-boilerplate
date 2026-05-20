package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-clean-api/config"
	"go-clean-api/internal/delivery/handler"
	"go-clean-api/internal/delivery/router"
	mysqlrepo "go-clean-api/internal/repository/mysql"
	authUsecase "go-clean-api/internal/usecase/auth"
	taskUsecase "go-clean-api/internal/usecase/task"
	userUsecase "go-clean-api/internal/usecase/user"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := mysqlrepo.NewDB(cfg.DB.DSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	userRepo := mysqlrepo.NewUserRepository(db)
	taskRepo := mysqlrepo.NewTaskRepository(db)

	authUC := authUsecase.New(userRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	userUC := userUsecase.New(userRepo)
	taskUC := taskUsecase.New(taskRepo)

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	router.Setup(r, router.Handlers{
		Auth: handler.NewAuthHandler(authUC),
		User: handler.NewUserHandler(userUC),
		Task: handler.NewTaskHandler(taskUC),
	}, cfg.JWT.Secret)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
