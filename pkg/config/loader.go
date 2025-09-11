package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// LoadConfig initializes Viper, reads the configuration file and environment variables,
// and validates the essential fields. Returns a Config pointer and an error.
func LoadConfig() (*Config, error) {
	// Optional: allow overriding config path via environment variable
	if path := viper.GetString("CONFIG_PATH"); path != "" {
		viper.AddConfigPath(path)
	}
	viper.AddConfigPath(".")      // fallback to current directory
	viper.SetConfigName("config") // config.json
	viper.SetConfigType("json")   // JSON type

	// Set defaults
	viper.SetDefault("server_addr", "0.0.0.0:8080")
	viper.SetDefault("shutdown_timeout", 30*time.Second)
	viper.SetDefault("read_header_timeout", 5*time.Second)
	viper.SetDefault("read_timeout", 10*time.Second)
	viper.SetDefault("write_timeout", 10*time.Second)
	viper.SetDefault("idle_timeout", 60*time.Second)

	// Bind environment variables
	viper.AutomaticEnv()
	viper.BindEnv("jwt_secret", "JWT_SECRET")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Non-file-not-found error
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// File not found is okay, fallback to env vars and defaults
	}

	// Unmarshal into Config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate essential fields
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable or config value is required")
	}

	return &cfg, nil
}
