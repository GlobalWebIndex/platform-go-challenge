package config_test

import (
	"os"
	"platform-go-challenge/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_New(t *testing.T) {
	os.Setenv("APP_USE_CACHE", "appusecache")
	os.Setenv("APP_PORT", "appport")
	os.Setenv("MONGO_DATABASE", "mongodatabase")
	os.Setenv("MONGO_USER", "mongouser")
	os.Setenv("MONGO_PASSWORD", "mongopass")
	os.Setenv("MONGO_HOST", "mongohost")
	os.Setenv("MONGO_PORT", "mongoport")
	os.Setenv("REDIS_PASSWORD", "redispass")
	os.Setenv("REDIS_HOST", "redishost")
	os.Setenv("REDIS_PORT", "redisport")
	os.Setenv("REDIS_DB", "redisdb")

	expCfg := &config.Config{
		App: config.App{
			WithCache: "appusecache",
			Port:      "appport",
		},
		Mongo: config.Mongo{
			Database: "mongodatabase",
			User:     "mongouser",
			Pass:     "mongopass",
			Host:     "mongohost",
			Port:     "mongoport",
		},
		Redis: config.Redis{
			Host:     "redishost",
			Port:     "redisport",
			Password: "redispass",
			Db:       "redisdb",
		},
	}

	cfg := config.New()
	assert.Equal(t, expCfg, cfg)
}
