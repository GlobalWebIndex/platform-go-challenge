package config

import (
	"os"
)

// Config is the main configuration of app.
type Config struct {
	App   App
	Mongo Mongo
	Redis Redis
}

// App contains app configuration.
type App struct {
	WithCache string
	Port      string
}

// Mongo contains mongo configuration.
type Mongo struct {
	Database string
	User     string
	Pass     string
	Host     string
	Port     string
}

// Redis contains redis configuration.
type Redis struct {
	Host     string
	Port     string
	Password string
	Db       string
}

// New constructor
func New() *Config {
	cfg := &Config{}
	cfg.setAppConfig()
	cfg.setMongoConfig()
	cfg.setRedisConfig()

	return cfg
}

// SetAppConfig creates an App struct.
func (cfg *Config) setAppConfig() {
	cfg.App = App{
		WithCache: os.Getenv("APP_USE_CACHE"),
		Port:      os.Getenv("APP_PORT"),
	}
}

// SetMongoConfig creates a Mongo config struct.
func (cfg *Config) setMongoConfig() {
	cfg.Mongo = Mongo{
		Database: os.Getenv("MONGO_DATABASE"),
		User:     os.Getenv("MONGO_USER"),
		Pass:     os.Getenv("MONGO_PASSWORD"),
		Host:     os.Getenv("MONGO_HOST"),
		Port:     os.Getenv("MONGO_PORT"),
	}
}

// SetRedisConfig creates a Redis config struct.
func (cfg *Config) setRedisConfig() {
	cfg.Redis = Redis{
		Password: os.Getenv("REDIS_PASSWORD"),
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Db:       os.Getenv("REDIS_DB"),
	}
}
