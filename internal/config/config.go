package config

import (
	"os"
	"time"
)

const (
	// Default values for database configuration
	defaultDBNetwork  = "tcp"            // Network type for database connection
	defaultDBAddress  = "localhost:5432" // Address of the database server
	defaultDBUser     = "gwi"            // Default database username
	defaultDBPassword = "gwi"            // Default database password
	defaultDBName     = "gwi"            // Default database name
	defaultDBAppName  = "gwi"            // Default application name for database connection
	// Default server configuration
	defaultAddr = ":8080" // Default address for the server
)

// Config holds the configuration for the server and database.
type Config struct {
	Server   ServerConfig   // Server-specific configuration
	Database DatabaseConfig // Database-specific configuration
}

// ServerConfig holds the configuration for the server.
type ServerConfig struct {
	Addr            string        // Address the server listens on
	ShutdownTimeout time.Duration // Timeout duration for server shutdown
}

// DatabaseConfig holds the configuration for the database connection.
type DatabaseConfig struct {
	Network  string // Network type for database connection (e.g., tcp)
	Address  string // Address of the database server
	User     string // Username for database authentication
	Password string // Password for database authentication
	Name     string // Name of the database
	AppName  string // Application name for database connection
}

// Load initializes and returns a Config struct populated with values from environment variables.
// If an environment variable is not set, it falls back to the default value.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Addr:            getEnv("ADDR", defaultAddr), // Load server address from environment or use default
			ShutdownTimeout: 10 * time.Second,            // Set server shutdown timeout
		},
		Database: DatabaseConfig{
			Network:  getEnv("DB_NETWORK", defaultDBNetwork),   // Load database network type from environment or use default
			Address:  getEnv("DB_ADDRESS", defaultDBAddress),   // Load database address from environment or use default
			User:     getEnv("DB_USER", defaultDBUser),         // Load database username from environment or use default
			Password: getEnv("DB_PASSWORD", defaultDBPassword), // Load database password from environment or use default
			Name:     getEnv("DB_NAME", defaultDBName),         // Load database name from environment or use default
			AppName:  getEnv("DB_APP_NAME", defaultDBAppName),  // Load application name for database from environment or use default
		},
	}
}

// getEnv retrieves the value of the environment variable specified by key.
// If the environment variable is not set, it returns the provided defaultValue.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
