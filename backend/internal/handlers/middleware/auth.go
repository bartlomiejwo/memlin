package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"backend/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey = contextKey("user")

type AuthClaims struct {
	UserID   int
	Role     string
	Language string
}

func JWTMiddleware(secretKey, refreshSecret string, accessExp, refreshExp time.Duration, prodEnv bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}

		// Parse token
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				// Attempt to refresh tokens (access & refresh)
				newAccessToken, newRefreshToken, err := getRefreshTokens(r, secretKey, refreshSecret, accessExp, refreshExp)
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				// Replace the request's Authorization header with the new access token
				r.Header.Set("Authorization", "Bearer "+newAccessToken)

				// Set the refresh token in the cookie
				utils.SetRefreshTokenCookie(w, newRefreshToken, refreshExp, prodEnv)

				tokenStr = newAccessToken

			case errors.Is(err, jwt.ErrTokenMalformed),
				errors.Is(err, jwt.ErrTokenSignatureInvalid):
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return

			default:
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		} else if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract user claims
		userIDFloat, ok1 := claims["user_id"].(float64)
		role, ok2 := claims["role"].(string)
		if !ok1 || !ok2 {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		authClaims := AuthClaims{
			UserID: int(userIDFloat),
			Role:   role,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, authClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Generates both access and refresh tokens
func getRefreshTokens(r *http.Request, secretKey, refreshSecret string, accessExp, refreshExp time.Duration) (string, string, error) {
	// Get refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		return "", "", err
	}
	refreshTokenStr := cookie.Value

	// Parse and validate refresh token
	refreshClaims := jwt.MapClaims{}
	refreshToken, err := jwt.ParseWithClaims(refreshTokenStr, refreshClaims, func(token *jwt.Token) (any, error) {
		return []byte(refreshSecret), nil
	})
	if err != nil || !refreshToken.Valid {
		return "", "", err
	}

	// Extract user ID and role
	userIDFloat, ok1 := refreshClaims["user_id"].(float64)
	role, ok2 := refreshClaims["role"].(string)
	if !ok1 || !ok2 {
		return "", "", errors.New("invalid refresh token claims")
	}
	userID := int(userIDFloat)

	// Generate new access token
	newAccessToken, err := utils.GenerateJWT(secretKey, userID, role, accessExp)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token
	newRefreshToken, err := utils.GenerateJWTRefresh(refreshSecret, userID, refreshExp)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// Retrieves authentication claims from the request context
func GetAuthClaims(r *http.Request) (AuthClaims, bool) {
	authClaims, ok := r.Context().Value(UserContextKey).(AuthClaims)
	return authClaims, ok
}
