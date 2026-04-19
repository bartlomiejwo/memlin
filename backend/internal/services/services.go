package services

import (
	"backend/internal/repositories"

	"go.uber.org/zap"
)

// Services struct holds all service instances
type Services struct {
	WordService *WordService
	UserService *UserService
}

// InitServices initializes all services and returns a Services struct
func InitServices(repos *repositories.Repositories, logger *zap.Logger) *Services {
	return &Services{
		WordService: NewWordService(repos.WordRepository, logger),
		UserService: NewUserService(repos.UserRepository, logger),
	}
}
