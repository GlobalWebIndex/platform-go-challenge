package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL string
	RedisAddr   string
}

func LoadConfig() *Config {
	// Load .env automatically
	_ = godotenv.Load()

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	redisAddr := os.Getenv("REDIS_ADDR")

	// Validate
	if user == "" || password == "" || dbName == "" || host == "" || port == "" {
		log.Fatal("Postgres environment variables are not set properly")
	}
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is not set")
	}

	// Build connection string dynamically
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName)

	return &Config{
		PostgresURL: postgresURL,
		RedisAddr:   redisAddr,
	}
}
