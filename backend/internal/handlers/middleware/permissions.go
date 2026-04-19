// middleware/permissions.go
package middleware

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
)

// PermissionMiddleware checks if the user has the required permission
func PermissionMiddleware(requiredPermission string, svc *services.UserService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user claims from the context (set by JWT middleware)
		authClaims, ok := GetAuthClaims(r)

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Fetch user permissions from DB using the user_id
		permissions, err := svc.GetPermissionsByUserID(r.Context(), authClaims.UserID)
		if err != nil {
			http.Error(w, "Failed to fetch permissions", http.StatusInternalServerError)
			return
		}

		// Check if user has the required permission
		if !hasPermission(permissions, requiredPermission) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Continue processing if permission is valid
		next.ServeHTTP(w, r)
	})
}

// hasPermission checks if the user has the required permission
func hasPermission(permissions []models.Permission, requiredPermission string) bool {
	// Create a map of permission codename for quick lookup
	permissionMap := make(map[string]struct{}, len(permissions))
	for _, p := range permissions {
		permissionMap[p.Codename] = struct{}{}
	}

	// Check if user has the required permission
	_, exists := permissionMap[requiredPermission]
	return exists
}
