package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds the application's main configuration settings
type Config struct {
	ServerAddr           string
	ProdEnv              bool
	DBURL                string
	CORSAllowedOrigins   []string
	CSRFSecretKey        string
	JWTSecretKey         string
	JWTRefreshSecretKey  string
	JWTExpiryMinutes     int
	JWTRefreshExpiryDays int
	GoogleClientID       string
	GoogleClientSecret   string
	IPRateLimit          int
	IPBurst              int
	UserRateLimit        int
	UserBurst            int
	UseUserRate          bool
}

// Load reads environment variables and loads the configuration
func Load() (*Config, error) {
	// Parse numeric and boolean values with error handling
	prodEnv, err := strconv.ParseBool(getEnv("PROD_ENV", "true"))
	if err != nil {
		return nil, fmt.Errorf("invalid PROD_ENV: %w", err)
	}
	jwtExpiryMinutes, err := strconv.Atoi(getEnv("JWT_EXPIRY_MINUTES", "15"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_TIME: %w", err)
	}
	jwtRefreshExpiryDays, err := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_DAYS", "7"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRY_DAYS: %w", err)
	}
	ipRateLimit, err := strconv.Atoi(getEnv("IP_RATE_LIMIT", "300"))
	if err != nil {
		return nil, fmt.Errorf("invalid IP_RATE_LIMIT: %w", err)
	}
	ipBurst, err := strconv.Atoi(getEnv("IP_BURST", "30"))
	if err != nil {
		return nil, fmt.Errorf("invalid IP_BURST: %w", err)
	}
	userRateLimit, err := strconv.Atoi(getEnv("USER_RATE_LIMIT", "60"))
	if err != nil {
		return nil, fmt.Errorf("invalid USER_RATE_LIMIT: %w", err)
	}
	userBurst, err := strconv.Atoi(getEnv("USER_BURST", "15"))
	if err != nil {
		return nil, fmt.Errorf("invalid USER_BURST: %w", err)
	}
	useUserRate, err := strconv.ParseBool(getEnv("USE_USER_RATE", "true"))
	if err != nil {
		return nil, fmt.Errorf("invalid USE_USER_RATE: %w", err)
	}

	// Split the CORS_ALLOWED_ORIGIN into a slice
	corsOrigins := strings.Split(getEnv("CORS_ALLOWED_ORIGIN", "*"), ",")

	return &Config{
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		ProdEnv:    prodEnv,
		DBURL: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "password"),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_NAME", "vocab_db"),
		),
		CORSAllowedOrigins:   corsOrigins,
		CSRFSecretKey:        getEnv("CSRF_SECRET_KEY", ""),
		JWTSecretKey:         getEnv("JWT_SECRET_KEY", ""),
		JWTRefreshSecretKey:  getEnv("JWT_REFRESH_SECRET_KEY", ""),
		JWTExpiryMinutes:     jwtExpiryMinutes,
		JWTRefreshExpiryDays: jwtRefreshExpiryDays,
		GoogleClientID:       getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:   getEnv("GOOGLE_CLIENT_SECRET", ""),
		IPRateLimit:          ipRateLimit,
		IPBurst:              ipBurst,
		UserRateLimit:        userRateLimit,
		UserBurst:            userBurst,
		UseUserRate:          useUserRate,
	}, nil
}

// getEnv fetches an environment variable or returns a default value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
