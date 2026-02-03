package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	originalDBURL := os.Getenv("DATABASE_URL")
	originalRedisAddr := os.Getenv("REDIS_ADDR")
	originalPort := os.Getenv("PORT")
	originalAppEnv := os.Getenv("APP_ENV")
	originalJWTSecret := os.Getenv("JWT_SECRET")

	defer func() {
		os.Setenv("DATABASE_URL", originalDBURL)
		os.Setenv("REDIS_ADDR", originalRedisAddr)
		os.Setenv("PORT", originalPort)
		os.Setenv("APP_ENV", originalAppEnv)
		os.Setenv("JWT_SECRET", originalJWTSecret)
	}()

	t.Run("success with all values set", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://localhost:5432/test")
		os.Setenv("REDIS_ADDR", "localhost:6379")
		os.Setenv("PORT", "9000")
		os.Setenv("APP_ENV", "test")
		os.Setenv("JWT_SECRET", "super-secret")

		cfg, err := Load()
		assert.NoError(t, err)
		assert.Equal(t, "postgres://localhost:5432/test", cfg.DatabaseURL)
		assert.Equal(t, "localhost:6379", cfg.RedisAddr)
		assert.Equal(t, "9000", cfg.Port)
		assert.Equal(t, "test", cfg.AppEnv)
		assert.Equal(t, "super-secret", cfg.JWTSecret)
	})

	t.Run("default values for Port and AppEnv", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://localhost:5432/test")
		os.Setenv("REDIS_ADDR", "localhost:6379")
		os.Setenv("JWT_SECRET", "super-secret")
		os.Unsetenv("PORT")
		os.Unsetenv("APP_ENV")

		cfg, err := Load()
		assert.NoError(t, err)
		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "production", cfg.AppEnv)
	})

	t.Run("missing DATABASE_URL", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")
		os.Setenv("REDIS_ADDR", "localhost:6379")
		os.Setenv("JWT_SECRET", "super-secret")

		_, err := Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DATABASE_URL is required")
	})

	t.Run("missing REDIS_ADDR", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://localhost:5432/test")
		os.Unsetenv("REDIS_ADDR")
		os.Setenv("JWT_SECRET", "super-secret")

		_, err := Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "REDIS_ADDR is required")
	})

	t.Run("missing JWT_SECRET", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://localhost:5432/test")
		os.Setenv("REDIS_ADDR", "localhost:6379")
		os.Unsetenv("JWT_SECRET")

		_, err := Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "JWT_SECRET is required")
	})
}
