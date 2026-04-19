package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/config"
	"backend/internal/services"

	"go.uber.org/zap"
)

func WordsHandler(svc *services.WordService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			svc.Logger.Warn("Invalid method on /api/words", zap.String("method", r.Method))
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context() // Get the context from the request

		// Pagination parameters
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = config.Settings.DefaultLimit
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = config.Settings.DefaultOffset
		}

		words, err := svc.GetWords(ctx, limit, offset)
		if err != nil {
			svc.Logger.Error("Failed to fetch words", zap.Error(err))
			http.Error(w, "Failed to fetch words", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(words); err != nil {
			svc.Logger.Error("Failed to encode words to JSON", zap.Error(err))
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
