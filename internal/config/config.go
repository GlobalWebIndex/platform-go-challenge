package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	LogLevel     string
	JWTSecret    string
}

func Load() *Config {
	return &Config{
		Port:         getEnvInt("PORT", 8080),
		ReadTimeout:  getEnvDuration("READ_TIMEOUT", 15*time.Second),
		WriteTimeout: getEnvDuration("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:  getEnvDuration("IDLE_TIMEOUT", 60*time.Second),
		LogLevel:     getEnvString("LOG_LEVEL", "info"),
		JWTSecret:    getEnvString("JWT_SECRET", "your-secret-key"),
	}
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
