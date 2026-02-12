package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/giannis84/platform-go-challenge/internal/auth"
	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "config.yaml"

// Config holds the application configuration.
type Config struct {
	APIPort    string `yaml:"api_port"`
	HealthPort string `yaml:"health_port"`

	// HTTP server timeouts (optional, defaults apply in server.go)
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`

	// JWT signing secret (env var only for testing). When empty, only unsigned tokens
	// (alg=none) are accepted if AllowUnsignedTokens is true.
	// Normally in production it should be fetched from a secrets provider like Vault, 
	// and not set via config file or env var.
	JWTSecret string `yaml:"-"`

	// AllowUnsignedTokens permits unsigned JWT tokens (alg=none) when true.
	// This should ONLY be enabled for local development and testing.
	// Requires explicit opt-in via ALLOW_UNSIGNED_TOKENS=true env var.
	AllowUnsignedTokens bool `yaml:"-"`

	// Database configuration (env vars only — secrets must not live in config.yaml)
	DBHost     string `yaml:"-"`
	DBPort     string `yaml:"-"`
	DBUser     string `yaml:"-"`
	DBPassword string `yaml:"-"`
	DBName     string `yaml:"-"`

	// Rate limiting configuration
	RateLimitRequests int           `yaml:"rate_limit_requests"` // Max requests per window (0 = disabled)
	RateLimitWindow   time.Duration `yaml:"rate_limit_window"`   // Time window for rate limiting
}

// Load reads configuration with the following precedence (highest wins):
//  1. Environment variables (API_PORT, HEALTH_PORT)
//  2. YAML config file (path from CONFIG_PATH env var, or "config.yaml")
//
// Database settings are loaded exclusively from environment variables.
func Load() (*Config, error) {
	cfg := &Config{}

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = defaultConfigPath
	}

	data, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parsing config file %s: %w", path, err)
		}
	}

	if v := os.Getenv("API_PORT"); v != "" {
		cfg.APIPort = v
	}
	if v := os.Getenv("HEALTH_PORT"); v != "" {
		cfg.HealthPort = v
	}

	if cfg.APIPort == "" {
		return nil, fmt.Errorf("api_port is required (set via config file or API_PORT env var)")
	}
	if cfg.HealthPort == "" {
		return nil, fmt.Errorf("health_port is required (set via config file or HEALTH_PORT env var)")
	}

	// Database configuration from environment variables
	cfg.DBHost = os.Getenv("POSTGRES_HOST")
	cfg.DBPort = os.Getenv("POSTGRES_PORT")
	cfg.DBUser = os.Getenv("POSTGRES_USER")
	cfg.DBPassword = os.Getenv("POSTGRES_PASSWORD")
	cfg.DBName = os.Getenv("POSTGRES_DB")

	// JWT secret (optional — when empty AND AllowUnsignedTokens is true, unsigned tokens are accepted)
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	// Allow unsigned tokens (explicit opt-in for dev/test only)
	cfg.AllowUnsignedTokens = os.Getenv("ALLOW_UNSIGNED_TOKENS") == "true"

	// HTTP server timeouts (optional — defaults apply in server.go if zero)
	if v := os.Getenv("READ_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.ReadTimeout = d
		}
	}
	if v := os.Getenv("WRITE_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.WriteTimeout = d
		}
	}
	if v := os.Getenv("IDLE_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.IdleTimeout = d
		}
	}

	if cfg.DBHost == "" {
		return nil, fmt.Errorf("POSTGRES_HOST env var is required")
	}
	if cfg.DBPort == "" {
		return nil, fmt.Errorf("POSTGRES_PORT env var is required")
	}
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("POSTGRES_USER env var is required")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD env var is required")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("POSTGRES_DB env var is required")
	}

	// Rate limiting configuration (env vars override config file)
	if v := os.Getenv("RATE_LIMIT_REQUESTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.RateLimitRequests = n
		}
	}
	if v := os.Getenv("RATE_LIMIT_WINDOW"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.RateLimitWindow = d
		}
	}

	// Apply rate limiting defaults if partially configured
	if cfg.RateLimitRequests > 0 && cfg.RateLimitWindow == 0 {
		cfg.RateLimitWindow = time.Minute // Default window: 1 minute
	}

	return cfg, nil
}

// PostgresConnString returns a PostgreSQL connection string.
func (c *Config) PostgresConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

// APIAddr returns the listen address for the API server.
func (c *Config) APIAddr() string {
	return ":" + c.APIPort
}

// HealthAddr returns the listen address for the health check server.
func (c *Config) HealthAddr() string {
	return ":" + c.HealthPort
}

// AuthConfig returns the JWT authentication configuration.
func (c *Config) AuthConfig() auth.AuthConfig {
	return auth.AuthConfig{
		Secret:              c.JWTSecret,
		AllowUnsignedTokens: c.AllowUnsignedTokens,
	}
}

// RateLimitConfig holds rate limiting settings.
type RateLimitConfig struct {
	Requests int           // Max requests per window (0 = disabled)
	Window   time.Duration // Time window for rate limiting
}

// RateLimitConfig returns the rate limiting configuration.
func (c *Config) RateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Requests: c.RateLimitRequests,
		Window:   c.RateLimitWindow,
	}
}
