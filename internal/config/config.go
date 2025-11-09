package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration loaded from environment variables.
// Each field has sensible defaults to make local development frictionless.
type Config struct {
	// API settings
	Port         string // Port to bind the HTTP server on
	AppEnv       string // Environment mode (development, production)

	// Middleware & limits
	LogEnabled      bool          // Enable HTTP request logging
	RateLimitMillis int           // Minimum interval between requests (per user/IP)
	MaxBodyBytes    int64         // Maximum allowed request body size (bytes)
	APIKey          string        // Optional shared API key for simple auth (empty disables auth)

	// Timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Log level placeholder for future structured logging
	LogLevel string
}

// LoadConfig reads environment variables, applies defaults and returns a populated Config struct.
// It uses helper functions to handle type conversion and default values gracefully.
func LoadConfig() *Config {
	cfg := &Config{
		Port:            getEnv("APP_PORT", "8080"),
		AppEnv:          getEnv("APP_ENV", "development"),
		LogEnabled:      getEnvBool("ENABLE_HTTP_LOG", true),
		RateLimitMillis: getEnvInt("RATE_LIMIT_MS", 50),
		MaxBodyBytes:    getEnvInt64("MAX_BODY_BYTES", 1<<20), // 1MB default
		ReadTimeout:     getEnvDurationSec("READ_TIMEOUT", 5),
		WriteTimeout:    getEnvDurationSec("WRITE_TIMEOUT", 10),
		IdleTimeout:     getEnvDurationSec("IDLE_TIMEOUT", 60),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		APIKey:          getEnv("API_KEY", ""), // empty -> auth disabled
	}
	log.Printf("Config loaded: %+v", cfg)
	return cfg
}

// --- helpers ---

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		switch v {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Printf("invalid int for %s=%s, using default %d", key, v, def)
	}
	return def
}

func getEnvInt64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
		log.Printf("invalid int64 for %s=%s, using default %d", key, v, def)
	}
	return def
}

func getEnvDurationSec(key string, def int) time.Duration {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return time.Duration(n) * time.Second
		}
		log.Printf("invalid duration for %s=%s, using default %ds", key, v, def)
	}
	return time.Duration(def) * time.Second
}
