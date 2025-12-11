// Package config provides configuration management for the application.
package config

// Config holds the application configuration.
type Config struct {
	// HTTPAddr is the address the HTTP server listens on (e.g., ":8080").
	HTTPAddr string
}
