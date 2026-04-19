package repositories

import (
	"backend/internal/db"

	"go.uber.org/zap"
)

// Repositories struct holds all repository instances
type Repositories struct {
	WordRepository *WordRepository
	UserRepository *UserRepository
}

// InitRepositories initializes all repositories and returns a Repositories struct
func InitRepositories(store *db.Store, logger *zap.Logger) *Repositories {
	return &Repositories{
		WordRepository: NewWordRepository(store, logger),
		UserRepository: NewUserRepository(store, logger),
	}
}
