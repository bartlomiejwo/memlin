package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"

	"go.uber.org/zap"
)

type UserService struct {
	repo   *repositories.UserRepository
	logger *zap.Logger
}

func NewUserService(repo *repositories.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserService) CreateOrUpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	return s.repo.CreateOrUpdateUser(ctx, user)
}

func (s *UserService) GetPermissionsByUserID(ctx context.Context, userID int) ([]models.Permission, error) {
	return s.repo.GetUserPermissions(ctx, userID)
}
