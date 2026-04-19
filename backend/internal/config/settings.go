package config

import (
	"os"
	"strconv"
)

// AppSettings holds global settings for the application
type AppSettings struct {
	DefaultLimit  int
	DefaultOffset int
}

// Global instance of settings
var Settings = loadSettings()

// loadSettings initializes settings from environment variables or defaults
func loadSettings() AppSettings {
	limit := getEnvInt("DEFAULT_LIMIT", 10)
	offset := getEnvInt("DEFAULT_OFFSET", 0)

	return AppSettings{
		DefaultLimit:  limit,
		DefaultOffset: offset,
	}
}

// Helper function to read integer env vars with a default fallback
func getEnvInt(key string, defaultValue int) int {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return num
}
