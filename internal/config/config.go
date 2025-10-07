package config

import (
	"os"
	"time"
)

type Config struct {
	Port        string
	DatabaseURL string
	ReqTimeout  time.Duration
}

func FromEnv() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dsn := os.Getenv("DATABASE_URL")
	return Config{
		Port:        ":" + port,
		DatabaseURL: dsn,
		ReqTimeout:  5 * time.Second,
	}
}
