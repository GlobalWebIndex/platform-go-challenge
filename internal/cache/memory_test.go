package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheStats_HitRate(t *testing.T) {
	assert.Equal(t, 0.0, (CacheStats{}).HitRate())

	s := CacheStats{Hits: 3, Misses: 1}
	assert.InDelta(t, 0.75, s.HitRate(), 0.000001)
}

func TestInMemoryCache_SetGet_HitAndMiss(t *testing.T) {
	ctx := context.Background()
	c := NewInMemoryCache()

	// miss
	v, ok := c.Get(ctx, "missing")
	assert.False(t, ok)
	assert.Nil(t, v)

	stats := c.GetStats()
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses)

	// set
	err := c.Set(ctx, "k", "v", 1*time.Second)
	assert.NoError(t, err)

	stats = c.GetStats()
	assert.Equal(t, int64(1), stats.Sets)

	// hit
	v, ok = c.Get(ctx, "k")
	assert.True(t, ok)
	assert.Equal(t, "v", v)

	stats = c.GetStats()
	assert.Equal(t, int64(1), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses) // unchanged
	assert.Equal(t, 1, stats.Size)
}

func TestInMemoryCache_ExpiredEntryBehavesLikeMiss(t *testing.T) {
	ctx := context.Background()
	c := NewInMemoryCache()

	err := c.Set(ctx, "k", "v", 10*time.Millisecond)
	assert.NoError(t, err)

	// Wait for expiry
	time.Sleep(20 * time.Millisecond)

	v, ok := c.Get(ctx, "k")
	assert.False(t, ok)
	assert.Nil(t, v)

	stats := c.GetStats()
	// One miss should be recorded when expired is accessed
	assert.Equal(t, int64(1), stats.Misses)
	// entry should be removed
	assert.Equal(t, 0, stats.Size)
}

func TestInMemoryCache_Delete_UpdatesStats(t *testing.T) {
	ctx := context.Background()
	c := NewInMemoryCache()

	_ = c.Set(ctx, "k", "v", 1*time.Minute)

	err := c.Delete(ctx, "k")
	assert.NoError(t, err)

	stats := c.GetStats()
	assert.Equal(t, int64(1), stats.Deletes)
	assert.Equal(t, 0, stats.Size)

	_, ok := c.Get(ctx, "k")
	assert.False(t, ok)
}

func TestInMemoryCache_Clear_RemovesAll(t *testing.T) {
	ctx := context.Background()
	c := NewInMemoryCache()

	_ = c.Set(ctx, "a", 1, 1*time.Minute)
	_ = c.Set(ctx, "b", 2, 1*time.Minute)

	stats := c.GetStats()
	assert.Equal(t, 2, stats.Size)

	err := c.Clear(ctx)
	assert.NoError(t, err)

	stats = c.GetStats()
	assert.Equal(t, 0, stats.Size)

	_, ok := c.Get(ctx, "a")
	assert.False(t, ok)
	_, ok = c.Get(ctx, "b")
	assert.False(t, ok)
}
