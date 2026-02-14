package memory

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gwi/platform-go-challenge/assets/application"
	"github.com/gwi/platform-go-challenge/assets/domain"
)

var _ application.Catalog = (*InMemoryCatalog)(nil)

// InMemoryCatalog implements application.Catalog with pre-seeded data; access is owner-scoped.
// We only log on forbidden access (wrong owner); successful lookups and batch sizes are not logged to avoid noise in production.
type InMemoryCatalog struct {
	mu        sync.RWMutex
	insights  map[string]domain.Insight
	charts    map[string]domain.Chart
	audiences map[string]domain.Audience
}

// NewInMemoryCatalog returns a catalog pre-seeded with data from Seed().
func NewInMemoryCatalog() *InMemoryCatalog {
	insights, charts, audiences := Seed(time.Now())
	return &InMemoryCatalog{
		insights:  insights,
		charts:    charts,
		audiences: audiences,
	}
}

func (c *InMemoryCatalog) GetInsight(ctx context.Context, ownerID string, id string) (domain.Insight, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	a, ok := c.insights[id]
	if !ok {
		return domain.Insight{}, false
	}
	if a.OwnerID != ownerID {
		slog.Warn("catalog: insight forbidden, owner mismatch", "id", id, "requested_owner", ownerID, "item_owner", a.OwnerID)
		return domain.Insight{}, false
	}
	return a, true
}

func (c *InMemoryCatalog) GetInsights(ctx context.Context, ownerID string, ids []string) map[string]domain.Insight {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]domain.Insight, len(ids))
	for _, id := range ids {
		if a, ok := c.insights[id]; ok && a.OwnerID == ownerID {
			out[id] = a
		}
	}
	return out
}

func (c *InMemoryCatalog) GetChart(ctx context.Context, ownerID string, id string) (domain.Chart, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	a, ok := c.charts[id]
	if !ok {
		return domain.Chart{}, false
	}
	if a.OwnerID != ownerID {
		slog.Warn("catalog: chart forbidden, owner mismatch", "id", id, "requested_owner", ownerID, "item_owner", a.OwnerID)
		return domain.Chart{}, false
	}
	return a, true
}

func (c *InMemoryCatalog) GetCharts(ctx context.Context, ownerID string, ids []string) map[string]domain.Chart {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]domain.Chart, len(ids))
	for _, id := range ids {
		if a, ok := c.charts[id]; ok && a.OwnerID == ownerID {
			out[id] = a
		}
	}
	return out
}

func (c *InMemoryCatalog) GetAudience(ctx context.Context, ownerID string, id string) (domain.Audience, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	a, ok := c.audiences[id]
	if !ok {
		return domain.Audience{}, false
	}
	if a.OwnerID != ownerID {
		slog.Warn("catalog: audience forbidden, owner mismatch", "id", id, "requested_owner", ownerID, "item_owner", a.OwnerID)
		return domain.Audience{}, false
	}
	return a, true
}

func (c *InMemoryCatalog) GetAudiences(ctx context.Context, ownerID string, ids []string) map[string]domain.Audience {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]domain.Audience, len(ids))
	for _, id := range ids {
		if a, ok := c.audiences[id]; ok && a.OwnerID == ownerID {
			out[id] = a
		}
	}
	return out
}
