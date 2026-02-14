// Package application defines the assets hexagon ports (read-only Catalog).
// Implementations live in assets/adapters (e.g. memory).
package application

import (
	"context"

	"github.com/gwi/platform-go-challenge/assets/domain"
)

// Catalog is the outbound port for reading catalog items by type and id (owner-scoped).
// Implemented by adapters (e.g. in-memory, persistence).
type Catalog interface {
	GetInsight(ctx context.Context, ownerID string, id string) (domain.Insight, bool)
	GetInsights(ctx context.Context, ownerID string, ids []string) map[string]domain.Insight
	GetChart(ctx context.Context, ownerID string, id string) (domain.Chart, bool)
	GetCharts(ctx context.Context, ownerID string, ids []string) map[string]domain.Chart
	GetAudience(ctx context.Context, ownerID string, id string) (domain.Audience, bool)
	GetAudiences(ctx context.Context, ownerID string, ids []string) map[string]domain.Audience
}
