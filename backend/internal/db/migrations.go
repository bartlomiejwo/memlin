package db

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// RunMigrations applies all pending migrations
func RunMigrations(db *pgxpool.Pool) error {
	// Convert pgxpool.Pool to *sql.DB
	sqlDB := stdlib.OpenDBFromPool(db)

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Get migration path from env or default
	migrationPath := getMigrationPath()

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to load migrations from %s: %w", migrationPath, err)
	}

	// Get the current version before running migrations
	version, dirty, _ := m.Version()
	log.Printf("Current DB version: %d (dirty: %v)", version, dirty)

	// Run migrations
	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println("✅ No migration changes needed.")
		return nil
	} else if err != nil {
		version, _, _ := m.Version()
		return fmt.Errorf("migration failed at version %d: %w", version, err)
	}

	log.Println("✅ Migrations applied successfully!")
	return nil
}

// getMigrationPath gets the migration directory from env or uses a default
func getMigrationPath() string {
	path := os.Getenv("MIGRATIONS_PATH")
	if path == "" {
		path = "file:///app/internal/db/migrations" // Default inside Docker container
	}
	return path
}

// RunMigrationDown rolls back the last migration
func RunMigrationDown(db *pgxpool.Pool) error {
	// Convert pgxpool.Pool to *sql.DB
	sqlDB := stdlib.OpenDBFromPool(db)

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(getMigrationPath(), "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	err = m.Steps(-1) // Rollback the last migration
	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	log.Println("⏪ Successfully rolled back last migration!")
	return nil
}
