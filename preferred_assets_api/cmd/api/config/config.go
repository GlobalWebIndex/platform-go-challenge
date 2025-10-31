package config

import (
	"os"
	"time"
)

type KeycloakConfig struct {
	Realm        string
	URL          string
	ExternalURL  string
	ClientID     string
	ClientSecret string
	Timeout      time.Duration
}

type Config struct {
	Keycloak KeycloakConfig
	Server   struct {
		Port string
	}
}

func Load() *Config {
	cfg := &Config{}

	// Keycloak configuration
	cfg.Keycloak.Realm = getEnv("KEYCLOAK_REALM", "your-realm")
	cfg.Keycloak.URL = getEnv("KEYCLOAK_URL", "http://localhost:8080")
	cfg.Keycloak.ExternalURL = getEnv("KEYCLOAK_EXTERNAL_URL", "http://localhost:8090")
	cfg.Keycloak.ClientID = getEnv("KEYCLOAK_CLIENT_ID", "your-client-id")
	cfg.Keycloak.ClientSecret = getEnv("KEYCLOAK_CLIENT_SECRET", "your-client-secret")
	cfg.Keycloak.Timeout = 10 * time.Second

	// Server configuration
	cfg.Server.Port = getEnv("SERVER_PORT", "8081")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
