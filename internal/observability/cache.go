package observability

import (
	"context"

	"go-favorites-app/internal/core/ports"
)

// InstrumentedCache is a decorator to intercept cache calls and record metrics.
type InstrumentedCache struct {
	inner ports.Cache
}

// NewInstrumentedCache creates a new instrumented cache wrapper.
func NewInstrumentedCache(inner ports.Cache) *InstrumentedCache {
	return &InstrumentedCache{inner: inner}
}

func (c *InstrumentedCache) AddToSet(ctx context.Context, id string, score float64) error {
	return c.inner.AddToSet(ctx, id, score)
}
func (c *InstrumentedCache) Set(ctx context.Context, id string, data []byte) error {
	return c.inner.Set(ctx, id, data)
}
func (c *InstrumentedCache) Remove(ctx context.Context, id string) error {
	return c.inner.Remove(ctx, id)
}
func (c *InstrumentedCache) Invalidate(ctx context.Context, id string) error {
	return c.inner.Invalidate(ctx, id)
}
func (c *InstrumentedCache) GetBatch(ctx context.Context, ids []string) (map[string][]byte, error) {
	res, err := c.inner.GetBatch(ctx, ids)
	if err == nil {
		hits := float64(len(res))
		misses := float64(len(ids) - len(res))
		cacheHits.Add(hits)
		cacheMisses.Add(misses)
	}
	return res, err
}
func (c *InstrumentedCache) GetIdsFromSet(ctx context.Context, start, stop int64) ([]string, error) {
	return c.inner.GetIdsFromSet(ctx, start, stop)
}
