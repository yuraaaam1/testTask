// @title Subscriptions API
// @version 1.0
// @description REST API для агрегации данных об онлайн подписках пользователей
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"log"

	"github.com/yuraaaam1/testTask/internal/config"
	"github.com/yuraaaam1/testTask/internal/handler"
	"github.com/yuraaaam1/testTask/internal/logger"
	"github.com/yuraaaam1/testTask/internal/repository"
	"github.com/yuraaaam1/testTask/internal/service"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Init(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Log.Sync()

	cfg := config.Load()

	db, err := repository.NewDB(cfg)
	if err != nil {
		logger.Log.Fatal("failed to connect to db", zap.Error(err))
	}
	defer db.Close()

	if err := repository.RunMigrations(db); err != nil {
		logger.Log.Fatal("failed to run migrations", zap.Error(err))
	}

	repo := repository.NewSubscriptionRepository(db)
	svc := service.NewSubscriptionService(repo)
	sh := handler.NewSubscriptionHandler(svc)

	router := handler.NewRouter(sh)

	logger.Log.Info("Server starting", zap.String("port", cfg.ServerPort))
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		logger.Log.Fatal("failed to start server", zap.Error(err))
	}
}
