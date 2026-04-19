package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ResponseWriter wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Implementing WriteHeader to capture status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// SecurityLoggingMiddleware logs basic info about all incoming requests
func SecurityLoggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log basic request info immediately
		logger.Info("Incoming request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
		)

		// Continue to next middleware/handler
		next.ServeHTTP(w, r)
	})
}

// DetailedLoggingMiddleware logs comprehensive info about successful requests
func DetailedLoggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create custom response writer to capture the status code
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200 if WriteHeader is never called
		}

		// Process the request
		next.ServeHTTP(wrappedWriter, r)

		// Log detailed info after completion
		duration := time.Since(start)
		logger.Info("Request completed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("status", wrappedWriter.statusCode),
			zap.Duration("duration", duration),
			zap.String("user_agent", r.UserAgent()),
			zap.String("referer", r.Referer()),
		)
	})
}
