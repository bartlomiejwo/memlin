package main

import (
	"context"
	"net/http"

	"backend/internal/app"
	"backend/internal/config"
	"backend/internal/db"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Ensure secrets are set
	if cfg.CSRFSecretKey == "" || cfg.JWTSecretKey == "" {
		logger.Fatal("CSRF_SECRET_KEY and JWT_SECRET_KEY must be set")
	}

	// Ensure google secrets are set
	if cfg.GoogleClientID == "" || cfg.GoogleClientSecret == "" {
		logger.Fatal("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
	}

	// Connect to database and get the Store
	store, err := db.Connect(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to DB", zap.Error(err))
	}
	defer store.Close()

	// Run migrations using the pgxpool from the store
	if err := db.RunMigrations(store.GetPool()); err != nil {
		logger.Fatal("Migration failed", zap.Error(err))
	}

	// Seed initial data if needed
	ctx := context.Background()
	if err := db.SeedData(ctx, store); err != nil {
		logger.Fatal("Seeding failed", zap.Error(err))
	}

	// Initialize app with the store instead of dbPool
	appInstance := app.New(store, cfg, logger)

	mux := http.NewServeMux()
	appInstance.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: mux,
	}

	logger.Info("Backend API running", zap.String("addr", cfg.ServerAddr))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
