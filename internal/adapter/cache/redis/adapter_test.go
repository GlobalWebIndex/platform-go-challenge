package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedisAdapter_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
	)
	if err != nil {
		t.Fatalf("failed to start redis: %v", err)
	}
	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate redis: %v", err)
		}
	}()

	endpoint, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// The redis-container module returns a URL like redis://localhost:port
	// but redis.NewClient expects just the host:port.
	// We need to strip the prefix if it exists.
	addr := endpoint
	if len(addr) > 8 && addr[:8] == "redis://" {
		addr = addr[8:]
	}

	adapter := NewAdapter(addr)
	defer adapter.client.Close()

	t.Run("Set and Get ids from set", func(t *testing.T) {
		id := "fav-1"
		err := adapter.AddToSet(ctx, id, 1.0)
		assert.NoError(t, err)

		ids, err := adapter.GetIdsFromSet(ctx, 0, -1)
		assert.NoError(t, err)
		assert.Contains(t, ids, id)
	})

	t.Run("Set and GetBatch", func(t *testing.T) {
		id := "fav-2"
		data := []byte(`{"id":"fav-2","name":"Redis Test"}`)

		err := adapter.Set(ctx, id, data)
		assert.NoError(t, err)

		batch, err := adapter.GetBatch(ctx, []string{id, "non-existent"})
		assert.NoError(t, err)
		assert.Equal(t, data, batch[id])
		assert.NotContains(t, batch, "non-existent")
	})

	t.Run("Remove", func(t *testing.T) {
		id := "fav-3"
		err := adapter.AddToSet(ctx, id, 3.0)
		assert.NoError(t, err)
		err = adapter.Set(ctx, id, []byte("data"))
		assert.NoError(t, err)

		err = adapter.Remove(ctx, id)
		assert.NoError(t, err)

		ids, _ := adapter.GetIdsFromSet(ctx, 0, -1)
		assert.NotContains(t, ids, id)

		batch, _ := adapter.GetBatch(ctx, []string{id})
		assert.Empty(t, batch)
	})
}
