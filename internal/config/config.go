package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL          string
	RedisAddr            string
	Port                 string
	AppEnv               string
	JWTSecret            string
	OtelExporterEndpoint string
}

// Load reads configuration from environment variables.
// It applies defaults for "local" environments but enforces strictness for others.
func Load() (Config, error) {
	cfg := Config{
		Port:                 os.Getenv("PORT"),
		AppEnv:               os.Getenv("APP_ENV"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		OtelExporterEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.JWTSecret == "" {
		if cfg.AppEnv == "local" {
			cfg.JWTSecret = "dev-secret-do-not-use-in-prod"
		} else {
			return Config{}, errors.New("JWT_SECRET is required")
		}
	}
	// Default to production safety if not explicitly set to local
	if cfg.AppEnv == "" {
		cfg.AppEnv = "production"
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	cfg.RedisAddr = os.Getenv("REDIS_ADDR")
	if cfg.RedisAddr == "" {
		return Config{}, errors.New("REDIS_ADDR is required")
	}

	return cfg, nil
}
