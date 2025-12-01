package config

import (
	"os"
	"time"
)

type ctxKey string

type Config struct {
	JWTSecret   []byte
	AdminAPIKey string
	ServerPort  string
	TokenExpiry time.Duration
	ContextKeys struct {
		UserID ctxKey
		Role   ctxKey
	}
}

func NewConfig() *Config {
	cfg := &Config{
		JWTSecret:   []byte(getEnvOrDefault("JWT_SECRET", "supersecretkey")),
		AdminAPIKey: getEnvOrDefault("ADMIN_API_KEY", "supersecret-admin-key"),
		ServerPort:  getEnvOrDefault("SERVER_PORT", ":8080"),
		TokenExpiry: time.Hour,
	}
	cfg.ContextKeys.UserID = "userID"
	cfg.ContextKeys.Role = "role"

	return cfg
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
