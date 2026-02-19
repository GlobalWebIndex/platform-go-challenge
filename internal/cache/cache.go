package cache

import (
	"context"
	"time"
)

// Cache defines the interface for caching operations
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) (interface{}, bool)

	// Set stores a value in cache with TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// Clear removes all values from cache
	Clear(ctx context.Context) error

	// GetStats returns cache statistics
	GetStats() CacheStats
}

// CacheStats holds cache performance metrics
type CacheStats struct {
	Hits    int64
	Misses  int64
	Sets    int64
	Deletes int64
	Size    int
}

// HitRate calculates the cache hit rate
func (s CacheStats) HitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total)
}
