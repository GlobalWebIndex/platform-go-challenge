package redis

import (
	"context"
	"time"

	"go-favorites-app/internal/core/ports"

	"github.com/redis/go-redis/v9"
)

type Adapter struct {
	client *redis.Client
}

func NewAdapter(addr string) *Adapter {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Adapter{client: rdb}
}

// Ensure Adapter implements ports.Cache
var _ ports.Cache = (*Adapter)(nil)

const (
	SetKey = "favorites:all"
	Prefix = "favorite:"
)

func (a *Adapter) AddToSet(ctx context.Context, id string, score float64) error {
	pipe := a.client.Pipeline()
	pipe.ZAdd(ctx, SetKey, redis.Z{Score: score, Member: id})
	// Refresh TTL if needed, but for "all" set usually we keep it or use logic
	_, err := pipe.Exec(ctx)
	return err
}

func (a *Adapter) Set(ctx context.Context, id string, data []byte) error {
	return a.client.Set(ctx, Prefix+id, data, 24*time.Hour).Err()
}

func (a *Adapter) GetBatch(ctx context.Context, ids []string) (map[string][]byte, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = Prefix + id
	}

	vals, err := a.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]byte)
	for i, val := range vals {
		if v, ok := val.(string); ok {
			result[ids[i]] = []byte(v)
		}
	}
	return result, nil
}

func (a *Adapter) GetIdsFromSet(ctx context.Context, start, stop int64) ([]string, error) {
	return a.client.ZRevRange(ctx, SetKey, start, stop).Result()
}

func (a *Adapter) Remove(ctx context.Context, id string) error {
	pipe := a.client.Pipeline()
	pipe.ZRem(ctx, SetKey, id)
	pipe.Del(ctx, Prefix+id)
	_, err := pipe.Exec(ctx)
	return err
}

func (a *Adapter) Invalidate(ctx context.Context, id string) error {
	return a.client.Del(ctx, Prefix+id).Err()
}
