package services

import (
	"context"

	"backend/internal/models"
	"backend/internal/repositories"

	"go.uber.org/zap"
)

type WordService struct {
	Repo   *repositories.WordRepository
	Logger *zap.Logger
}

func NewWordService(repo *repositories.WordRepository, logger *zap.Logger) *WordService {
	if repo == nil {
		panic("WordRepository is not initialized!")
	}

	return &WordService{
		Repo:   repo,
		Logger: logger,
	}
}

func (s *WordService) GetWords(ctx context.Context, limit, offset int) ([]models.Word, error) {
	words, err := s.Repo.GetWords(ctx, limit, offset)
	if err != nil {
		s.Logger.Debug("Failed to get words from repository", zap.Error(err))
		return nil, err
	}
	return words, nil
}
