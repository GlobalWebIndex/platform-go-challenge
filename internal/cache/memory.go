package cache

import (
	"context"
	"sync"
	"time"
)

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

type InMemoryCache struct {
	mu      sync.RWMutex
	data    map[string]cacheEntry
	stats   CacheStats
	statsMu sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		data: make(map[string]cacheEntry),
	}

	go cache.cleanupExpired()

	return cache
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.RLock()
	entry, exists := c.data[key]
	c.mu.RUnlock()

	if !exists {
		c.recordMiss()
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		c.recordMiss()
		return nil, false
	}

	c.recordHit()
	return entry.value, true
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)

	c.mu.Lock()
	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: expiresAt,
	}
	c.mu.Unlock()

	c.recordSet()
	return nil
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()

	c.recordDelete()
	return nil
}

func (c *InMemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	c.data = make(map[string]cacheEntry)
	c.mu.Unlock()

	return nil
}

// GetStats returns cache statistics
func (c *InMemoryCache) GetStats() CacheStats {
	c.statsMu.RLock()
	defer c.statsMu.RUnlock()

	c.mu.RLock()
	size := len(c.data)
	c.mu.RUnlock()

	stats := c.stats
	stats.Size = size
	return stats
}

func (c *InMemoryCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.data {
			if now.After(entry.expiresAt) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *InMemoryCache) recordHit() {
	c.statsMu.Lock()
	c.stats.Hits++
	c.statsMu.Unlock()
}

func (c *InMemoryCache) recordMiss() {
	c.statsMu.Lock()
	c.stats.Misses++
	c.statsMu.Unlock()
}

func (c *InMemoryCache) recordSet() {
	c.statsMu.Lock()
	c.stats.Sets++
	c.statsMu.Unlock()
}

func (c *InMemoryCache) recordDelete() {
	c.statsMu.Lock()
	c.stats.Deletes++
	c.statsMu.Unlock()
}
