package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	JWT    JWTConfig    `yaml:"jwt"`
}

type ServerConfig struct {
	Port                   int   `yaml:"port"`
	ShutdownTimeoutSeconds int   `yaml:"shutdown_timeout_seconds"`
	ReadTimeoutSeconds     int   `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds    int   `yaml:"write_timeout_seconds"`
	IdleTimeoutSeconds     int   `yaml:"idle_timeout_seconds"`
	MaxRequestBodyBytes    int64 `yaml:"max_request_body_bytes"`
	RateLimitPerSecond     int   `yaml:"rate_limit_per_second"`
	RateLimitBurst         int   `yaml:"rate_limit_burst"`
}

// JWTConfig: use JWT_SECRET from env in production (no default).
type JWTConfig struct {
	TokenExpirySeconds int `yaml:"token_expiry_seconds"`
}

func DefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Port:                   8080,
			ShutdownTimeoutSeconds: 10,
			ReadTimeoutSeconds:     15,
			WriteTimeoutSeconds:    15,
			IdleTimeoutSeconds:     60,
			MaxRequestBodyBytes:    1048576, // 1MB
			RateLimitPerSecond:     100,
			RateLimitBurst:         200,
		},
		JWT: JWTConfig{
			TokenExpirySeconds: 3600,
		},
	}
}

// Load merges config from path with defaults. Missing file returns default config.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}
	applyDefaults(&cfg)
	return cfg, nil
}

func applyDefaults(c *Config) {
	if c.Server.Port <= 0 {
		c.Server.Port = 8080
	}
	if c.Server.ShutdownTimeoutSeconds <= 0 {
		c.Server.ShutdownTimeoutSeconds = 10
	}
	if c.Server.ReadTimeoutSeconds <= 0 {
		c.Server.ReadTimeoutSeconds = 15
	}
	if c.Server.WriteTimeoutSeconds <= 0 {
		c.Server.WriteTimeoutSeconds = 15
	}
	if c.Server.IdleTimeoutSeconds <= 0 {
		c.Server.IdleTimeoutSeconds = 60
	}
	if c.Server.MaxRequestBodyBytes <= 0 {
		c.Server.MaxRequestBodyBytes = 1048576 // 1MB
	}
	if c.Server.RateLimitPerSecond <= 0 {
		c.Server.RateLimitPerSecond = 100
	}
	if c.Server.RateLimitBurst <= 0 {
		c.Server.RateLimitBurst = 200
	}
	if c.JWT.TokenExpirySeconds <= 0 {
		c.JWT.TokenExpirySeconds = 3600
	}
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

func (c *Config) ShutdownTimeout() time.Duration {
	return time.Duration(c.Server.ShutdownTimeoutSeconds) * time.Second
}

func (c *Config) TokenExpiry() time.Duration {
	return time.Duration(c.JWT.TokenExpirySeconds) * time.Second
}

func (c *Config) ReadTimeout() time.Duration {
	return time.Duration(c.Server.ReadTimeoutSeconds) * time.Second
}

func (c *Config) WriteTimeout() time.Duration {
	return time.Duration(c.Server.WriteTimeoutSeconds) * time.Second
}

func (c *Config) IdleTimeout() time.Duration {
	return time.Duration(c.Server.IdleTimeoutSeconds) * time.Second
}
