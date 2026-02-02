package application

import (
	"context"

	"github.com/manos/favourites/assets/domain"
)

// AssetRepository defines the minimal interface the application needs from a repository.
type AssetRepository interface {
	GetInsight(ctx context.Context, userId string, id string) (domain.InsightAsset, bool)
	GetAudience(ctx context.Context, userId string, id string) (domain.AudienceAsset, bool)
	GetChart(ctx context.Context, userId string, id string) (domain.ChartAsset, bool)
	GetInsights(ctx context.Context, userId string, ids []string) map[string]domain.InsightAsset
	GetAudiences(ctx context.Context, userId string, ids []string) map[string]domain.AudienceAsset
	GetCharts(ctx context.Context, userId string, ids []string) map[string]domain.ChartAsset
}
