package assetcatalog

import (
	"context"
	"strings"
	"time"

	assetsapp "github.com/gwi/platform-go-challenge/assets/application"
	assetsdomain "github.com/gwi/platform-go-challenge/assets/domain"
	"github.com/gwi/platform-go-challenge/favourites/application"
	"github.com/gwi/platform-go-challenge/favourites/domain"
)

// Client adapts the assets Catalog to favourites' AssetCatalogClient.
// Resolves source_asset_id by kind (insight/chart/audience) with owner-scoped lookup.
type Client struct {
	catalog assetsapp.Catalog
}

var _ application.AssetCatalogClient = (*Client)(nil)

// NewClient returns an adapter that uses the assets catalog to resolve source_asset_id.
func NewClient(catalog assetsapp.Catalog) *Client {
	return &Client{catalog: catalog}
}

// GetAsset returns the catalog item as a favourites Asset, or ErrAssetNotFound when not found or not owned.
func (c *Client) GetAsset(userID string, sourceAssetID string, assetType string) (*domain.Asset, error) {
	ctx := context.Background()
	now := time.Now()

	switch assetType {
	case domain.AssetTypeInsight:
		a, ok := c.catalog.GetInsight(ctx, userID, sourceAssetID)
		if !ok {
			return nil, application.ErrAssetNotFound
		}
		out := insightToFavourite(a, now)
		return &out, nil
	case domain.AssetTypeChart:
		a, ok := c.catalog.GetChart(ctx, userID, sourceAssetID)
		if !ok {
			return nil, application.ErrAssetNotFound
		}
		out := chartToFavourite(a, now)
		return &out, nil
	case domain.AssetTypeAudience:
		a, ok := c.catalog.GetAudience(ctx, userID, sourceAssetID)
		if !ok {
			return nil, application.ErrAssetNotFound
		}
		out := audienceToFavourite(a, now)
		return &out, nil
	default:
		return nil, application.ErrAssetNotFound
	}
}

func insightToFavourite(a assetsdomain.Insight, now time.Time) domain.Asset {
	return domain.Asset{
		Type:          domain.AssetTypeInsight,
		Description:   a.Summary,
		SourceAssetID: a.ID,
		CreatedAt:     now,
		UpdatedAt:     now,
		Insight:       &domain.InsightData{Text: a.Body},
	}
}

func chartToFavourite(a assetsdomain.Chart, now time.Time) domain.Asset {
	return domain.Asset{
		Type:          domain.AssetTypeChart,
		Description:   a.Summary,
		SourceAssetID: a.ID,
		CreatedAt:     now,
		UpdatedAt:     now,
		Chart: &domain.ChartData{
			Title:      a.Name,
			XAxisTitle: a.AxisX,
			YAxisTitle: a.AxisY,
			Data:       a.Series,
		},
	}
}

func audienceToFavourite(a assetsdomain.Audience, now time.Time) domain.Asset {
	return domain.Asset{
		Type:          domain.AssetTypeAudience,
		Description:   a.Summary,
		SourceAssetID: a.ID,
		CreatedAt:     now,
		UpdatedAt:     now,
		Audience: &domain.AudienceData{
			SampleSize:        a.SampleSize,
			EstimatedReach:    a.Reach,
			PopulationPercent: a.SharePct,
			Gender:            joinStrings(a.Gender),
			BirthCountry:      joinStrings(a.BirthCountry),
			AgeGroups:         joinStrings(a.AgeGroups),
			HoursSocialMedia:   joinStrings(a.HoursSocialMedia),
			Purchases:         joinStrings(a.Purchases),
		},
	}
}

func joinStrings(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.Join(s, ", ")
}
