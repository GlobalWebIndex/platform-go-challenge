package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// Config Redis Configuration.
type Config struct {
	Host     string
	Port     string
	Password string
	Db       string
}

// New returns a handle to a Redis client.
func New(ctx context.Context, cfg Config) (*redis.Client, error) {
	redisAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	redisDb, err := strconv.Atoi(cfg.Db)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: cfg.Password,
		DB:       redisDb,
	})
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return redisClient, err
}
