package middleware

import (
	"backend/internal/localization"
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// Define a custom type for the context key
const LanguageContextKey = contextKey("lang")

// DetectLanguage extracts the language from the JWT claims (if available) and the Accept-Language header
func DetectLanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := language.English.String() // Default language

		// Retrieve language from JWT claims (if authenticated)
		authClaims, ok := GetAuthClaims(r)
		if ok && authClaims.Language != "" {
			// If the language is found in the claims, prioritize it
			lang = authClaims.Language
		} else {
			// If not in claims, check the Accept-Language header
			acceptLang := r.Header.Get("Accept-Language")
			if acceptLang != "" {
				langs := strings.Split(acceptLang, ",")
				if len(langs) > 0 {
					lang = strings.TrimSpace(strings.Split(langs[0], ";")[0])
				}
			}
		}

		// Add language to request context using the custom type
		ctx := r.Context()
		ctx = context.WithValue(ctx, LanguageContextKey, lang)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetLanguage(r *http.Request) string {
	lang, ok := r.Context().Value(LanguageContextKey).(string)
	if !ok {
		return language.English.String()
	}
	return lang
}

// LocalizerContextKey is a custom key type to avoid context key collisions
const LocalizerContextKey = contextKey("localizer")

// LocalizerMiddleware injects the localizer into the request context
func LocalizerMiddleware(localizer *localization.Localizer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add localizer to the context
		ctx := context.WithValue(r.Context(), LocalizerContextKey, localizer)
		// Pass the modified request with the new context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetLocalizer extracts the localizer from the request context
func GetLocalizer(r *http.Request) *localization.Localizer {
	localizer, ok := r.Context().Value(LocalizerContextKey).(*localization.Localizer)
	if !ok {
		logger, _ := zap.NewProduction()
		defer logger.Sync() // Ensure logs are flushed properly
		return localization.NewLocalizer(logger)
	}
	return localizer
}
