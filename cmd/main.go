package main

import (
	"log"

	"github.com/yuraaaam1/testTask/internal/config"
	"github.com/yuraaaam1/testTask/internal/handler"
	"github.com/yuraaaam1/testTask/internal/repository"
	"github.com/yuraaaam1/testTask/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	repo := repository.NewSubscriptionRepository(db)
	svc := service.NewSubscriptionService(repo)
	sh := handler.NewSubscriptionHandler(svc)

	router := handler.NewRouter(sh)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
