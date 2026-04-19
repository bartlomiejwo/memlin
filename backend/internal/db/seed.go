package db

import (
	"context"
	"fmt"
	"log"

	sqlcdb "backend/internal/db/sqlc"
)

func SeedData(ctx context.Context, store *Store) error {
	// Get a count of languages
	languages, err := store.ListLanguages(ctx)
	if err != nil {
		return fmt.Errorf("failed to check languages: %w", err)
	}

	// Seed languages if none exist
	if len(languages) == 0 {
		log.Println("Seeding languages...")

		// Using a transaction for seeding languages
		err := store.WithTx(ctx, func(q *sqlcdb.Queries) error {
			_, err := q.CreateLanguage(ctx, sqlcdb.CreateLanguageParams{Code: "en", Name: "English"})
			if err != nil {
				return err
			}

			_, err = q.CreateLanguage(ctx, sqlcdb.CreateLanguageParams{Code: "de", Name: "German"})
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("languages seeding failed: %w", err)
		}
	}

	// Check if we have words
	words, err := store.ListWords(ctx, sqlcdb.ListWordsParams{Limit: 10, Offset: 0})
	if err != nil {
		return fmt.Errorf("failed to check words: %w", err)
	}

	// Seed words if none exist
	if len(words) == 0 {
		log.Println("Seeding words...")

		// First get language IDs
		enLang, err := store.GetLanguageByCode(ctx, sqlcdb.GetLanguageByCodeParams{Code: "en"})
		if err != nil {
			return fmt.Errorf("failed to get English language ID: %w", err)
		}

		deLang, err := store.GetLanguageByCode(ctx, sqlcdb.GetLanguageByCodeParams{Code: "de"})
		if err != nil {
			return fmt.Errorf("failed to get German language ID: %w", err)
		}

		// Using raw query for words seeding as it's simpler for seed data
		// In a real app, you might want to create proper sqlc queries for this too
		_, err = store.pool.Exec(ctx, `
            INSERT INTO words (language_id, word, pronunciation, category, level, popularity) VALUES
            ($1, 'doctor', '/ˈdɒk.tər/', 'Professions & Jobs', 'B2', 0.9),
            ($2, 'der Arzt', '/ˈaʁt͡st/', 'Professions & Jobs', 'B2', 0.9),
            ($2, 'die Ärztin', '/ˈɛʁt͡s.tɪn/', 'Professions & Jobs', 'B2', 0.9)
            ON CONFLICT DO NOTHING;
        `, enLang.ID, deLang.ID)
		if err != nil {
			return fmt.Errorf("words seeding failed: %w", err)
		}
	}

	log.Println("✅ Initial data seeded successfully")
	return nil
}
