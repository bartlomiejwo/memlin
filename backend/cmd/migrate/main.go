package main

import (
	"fmt"
	"os"

	"backend/internal/config"
	"backend/internal/db"

	"go.uber.org/zap"
)

func main() {
	// Set up logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Connect to DB
	store, err := db.Connect(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to DB", zap.Error(err))
	}
	defer store.Close()

	// Ensure a command is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <up|down>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		if err := db.RunMigrations(store.GetPool()); err != nil {
			logger.Fatal("Migration failed", zap.Error(err))
		}
		logger.Info("✅ Migrations applied successfully!")
	case "down":
		if err := db.RunMigrationDown(store.GetPool()); err != nil {
			logger.Fatal("Rollback failed", zap.Error(err))
		}
		logger.Info("⏪ Rollback completed!")
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Usage: migrate <up|down>")
		os.Exit(1)
	}
}
