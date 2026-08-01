package main

import (
	"log"

	"cy11dsphk/server/internal/config"
	"cy11dsphk/server/internal/database"
	"cy11dsphk/server/internal/router"
)

func main() {
	cfg := config.Load()

	if err := database.Init(cfg.Database); err != nil {
		log.Fatalf("init database failed: %v", err)
	}

	if err := database.MigrateAndSeed(cfg); err != nil {
		log.Fatalf("migrate database failed: %v", err)
	}

	r := router.New(cfg)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
