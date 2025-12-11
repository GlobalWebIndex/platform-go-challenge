package config

import (
	"os"
)

const (
	defaultHTTPAddr = ":8080"     // Default HTTP server address
	envHTTPAddr     = "HTTP_ADDR" // Environment variable name for HTTP address
)

// Load reads configuration from the environment, applying defaults where needed.
func Load() Config {
	addr := defaultHTTPAddr
	if v := os.Getenv(envHTTPAddr); v != "" {
		addr = v
	}

	return Config{
		HTTPAddr: addr,
	}
}
