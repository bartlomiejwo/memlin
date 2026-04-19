package middleware

import (
	"net/http"
	"slices"
)

// RoleMiddleware checks if the user has one of the required roles based on JWT claims
func RoleMiddleware(requiredRoles []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user claims from the context (already set by the JWT middleware)
		authClaims, ok := GetAuthClaims(r)

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the user's role matches any of the required roles
		if !hasRole(authClaims.Role, requiredRoles) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Continue processing if role matches
		next.ServeHTTP(w, r)
	})
}

// hasRole checks if the user's role matches any of the required roles
func hasRole(userRole string, requiredRoles []string) bool {
	// Check if the user has any of the required roles
	return slices.Contains(requiredRoles, userRole)
}
