package utils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(secret string, userID int, role string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateJWTRefresh(secret string, userID int, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string, expiry time.Duration, prodEnv bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   prodEnv,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(expiry),
		Path:     "/", // optionally restrict to "/api" etc.
	})
}
