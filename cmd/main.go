package main

import (
	"log"

	"github.com/yuraaaam1/testTask/internal/config"
	"github.com/yuraaaam1/testTask/internal/repository"
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

	log.Println("Migrations applied successfully")
}
