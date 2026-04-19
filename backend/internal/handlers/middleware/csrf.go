package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

// CSRFProtectionMiddleware adds CSRF protection
func CSRFProtectionMiddleware(secretKey string, next http.Handler) http.Handler {
	csrfMiddleware := csrf.Protect([]byte(secretKey))
	return csrfMiddleware(next)
}
