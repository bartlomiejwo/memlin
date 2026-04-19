package repositories

import (
	"context"

	"backend/internal/db"
	sqlcdb "backend/internal/db/sqlc"
	"backend/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type WordRepository struct {
	store  *db.Store
	logger *zap.Logger
}

func NewWordRepository(store *db.Store, logger *zap.Logger) *WordRepository {
	return &WordRepository{
		store:  store,
		logger: logger,
	}
}

func (r *WordRepository) GetWords(ctx context.Context, limit, offset int) ([]models.Word, error) {
	words, err := r.store.ListWords(ctx, sqlcdb.ListWordsParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		r.logger.Debug("Failed to query words", zap.Error(err))
		return nil, err
	}

	// Convert from sqlc model to our application model
	result := make([]models.Word, len(words))
	for i, w := range words {
		result[i] = models.Word{
			ID:            int(w.ID),
			Word:          w.Word,
			Language:      w.Language,
			Pronunciation: w.Pronunciation.String,
			Category:      w.Category.String,
			Level:         w.Level.String,
			Popularity:    w.Popularity.Float64,
		}
	}

	return result, nil
}

func (r *WordRepository) GetWord(ctx context.Context, id int) (*models.Word, error) {
	word, err := r.store.GetWord(ctx, sqlcdb.GetWordParams{ID: int32(id)})
	if err != nil {
		r.logger.Debug("Failed to get word", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	return &models.Word{
		ID:            int(word.ID),
		Word:          word.Word,
		Language:      word.Language,
		Pronunciation: word.Pronunciation.String,
		Category:      word.Category.String,
		Level:         word.Level.String,
		Popularity:    word.Popularity.Float64,
	}, nil
}

func (r *WordRepository) CreateWord(ctx context.Context, params models.WordCreateParams) (*models.Word, error) {
	// Use a transaction to ensure data consistency
	var result *models.Word

	err := r.store.WithTx(ctx, func(q *sqlcdb.Queries) error {
		// First, make sure the language exists
		_, err := q.GetLanguageByCode(ctx, sqlcdb.GetLanguageByCodeParams{Code: params.LanguageCode})
		if err != nil {
			r.logger.Debug("Language not found", zap.String("code", params.LanguageCode), zap.Error(err))
			return err
		}

		// Then create the word
		word, err := q.CreateWord(ctx, sqlcdb.CreateWordParams{
			Code:          params.LanguageCode,
			Word:          params.Word,
			Pronunciation: pgtype.Text{String: params.Pronunciation, Valid: true},
			Category:      pgtype.Text{String: params.Category, Valid: true},
			Level:         pgtype.Text{String: params.Level, Valid: true},
			Popularity:    pgtype.Float8{Float64: params.Popularity, Valid: true},
		})
		if err != nil {
			return err
		}

		// Convert to our application model
		result = &models.Word{
			ID:            int(word.ID),
			Word:          word.Word,
			Language:      params.LanguageCode,
			Pronunciation: word.Pronunciation.String,
			Category:      word.Category.String,
			Level:         word.Level.String,
			Popularity:    word.Popularity.Float64,
		}
		return nil
	})

	if err != nil {
		r.logger.Debug("Failed to create word", zap.Error(err))
		return nil, err
	}

	return result, nil
}
