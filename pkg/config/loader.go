package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// LoadConfig initializes Viper and loads the configuration.
// It returns a pointer to the Config struct and an error.
func LoadConfig() (*Config, error) {
	// Set the file name and path for the config file.
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Set defaults for all configuration fields.
	viper.SetDefault("server_addr", "0.0.0.0:8080")
	viper.SetDefault("shutdown_timeout", 30*time.Second)
	viper.SetDefault("read_header_timeout", 5*time.Second)
	viper.SetDefault("read_timeout", 10*time.Second)
	viper.SetDefault("write_timeout", 10*time.Second)
	viper.SetDefault("idle_timeout", 60*time.Second)
	viper.AutomaticEnv()
	viper.BindEnv("jwt_secret", "JWT_SECRET")

	// Read the config file. If it doesn't exist, Viper will use the defaults and env vars.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// A non-file-not-found error means something went wrong.
			return nil, err
		}
		log.Println("config file not found, using environment variables and defaults")
	}

	// Unmarshal the loaded values into the Config struct.
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Basic validation for the JWT secret.
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable or config value is required.")
	}

	log.Println("Loaded config:", cfg)

	return &cfg, nil
}
