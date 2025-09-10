package config

import (
	"os"
)

// Config holds the minimal app settings.
type Config struct {
	AppEnv   string
	HTTPAddr string

	DBHost    string
	DBPort    string
	DBName    string
	DBUser    string
	DBPass    string
	DBSSLMode string

	LogLevel string
	LogPath  string
}

// LoadFromEnv builds a Config by reading environment variables.
// If a variable is missing, it falls back to a safe default.
func LoadFromEnv() *Config {
	return &Config{
		AppEnv:   getEnvOrFallback("APP_ENV", "dev"),
		HTTPAddr: getEnvOrFallback("HTTP_ADDR", ":8080"),

		DBHost:    getEnvOrFallback("DB_HOST", "postgres"),
		DBPort:    getEnvOrFallback("DB_PORT", "5432"),
		DBName:    getEnvOrFallback("DB_NAME", "favs"),
		DBUser:    getEnvOrFallback("DB_USER", "app"),
		DBPass:    getEnvOrFallback("DB_PASS", "app"),
		DBSSLMode: getEnvOrFallback("DB_SSLMODE", "disable"),

		LogLevel: getEnvOrFallback("LOG_LEVEL", "info"),
		LogPath:  getEnvOrFallback("LOG_PATH", "./logs/"),
	}
}

// getEnvOrFallback returns the environment variable value if set,
// otherwise it returns the provided fallback string.
func getEnvOrFallback(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
