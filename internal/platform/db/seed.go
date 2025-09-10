package db

import (
	"context"
	"time"

	"github.com/SokratisChaimanas/platform-go-challenge/ent"
	"github.com/SokratisChaimanas/platform-go-challenge/ent/asset"
	"github.com/google/uuid"
)

// seedDevOnce seeds fixed data when running in dev (or SEED=1).
// It runs only if the DB has no users or no assets.
func seedDevOnce(ctx context.Context, client *ent.Client, appEnv string) error {
	// Check if data already exists: if any users or assets, skip seeding.
	hasUsers, err := client.User.Query().Exist(ctx)
	if err != nil {
		return err
	}

	hasAssets, err := client.Asset.Query().Exist(ctx)
	if err != nil {
		return err
	}

	if hasUsers || hasAssets {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback() }()

	now := time.Now().UTC()

	// --- Seed Users (3 fixed UUIDs).
	userIDs := []uuid.UUID{
		uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		uuid.MustParse("33333333-3333-3333-3333-333333333333"),
	}
	for _, id := range userIDs {
		if _, err := tx.User.
			Create().
			SetID(id).
			SetCreatedAt(now).
			Save(ctx); err != nil {
			return err
		}
	}

	// --- Seed Assets (2 per type, fixed UUIDs).
	assetsToSeed := []struct {
		ID          uuid.UUID
		Type        asset.AssetType
		Description string
		Payload     map[string]any
	}{
		// Chart
		{
			ID:          uuid.MustParse("aaaaaaa1-0000-0000-0000-000000000001"),
			Type:        asset.AssetTypeChart,
			Description: "Daily active users - last 7 days",
			Payload: map[string]any{
				"title": "DAU",
				"x":     []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
				"y":     []int{120, 150, 160, 170, 180, 200, 220},
			},
		},
		{
			ID:          uuid.MustParse("aaaaaaa1-0000-0000-0000-000000000002"),
			Type:        asset.AssetTypeChart,
			Description: "Conversion rate by channel",
			Payload: map[string]any{
				"title":   "CR by Channel",
				"series":  []string{"Email", "Ads", "Organic"},
				"values":  []float64{2.1, 1.3, 3.2},
				"y_label": "%",
			},
		},
		// Insight
		{
			ID:          uuid.MustParse("bbbbbbb2-0000-0000-0000-000000000001"),
			Type:        asset.AssetTypeInsight,
			Description: "Insight: heavy social usage",
			Payload: map[string]any{
				"text": "40% of millennials spend more than 3 hours on social media daily",
			},
		},
		{
			ID:          uuid.MustParse("bbbbbbb2-0000-0000-0000-000000000002"),
			Type:        asset.AssetTypeInsight,
			Description: "Insight: cart abandonment",
			Payload: map[string]any{
				"text": "Cart abandonment decreased 8% after one-click checkout rollout",
			},
		},
		// Audience
		{
			ID:          uuid.MustParse("ccccccc3-0000-0000-0000-000000000001"),
			Type:        asset.AssetTypeAudience,
			Description: "Audience: M 24-35, >3h/day social, >1 purchase last month",
			Payload: map[string]any{
				"gender":          "Male",
				"age_group":       "24-35",
				"country":         "US",
				"social_hours":    ">3",
				"purchases_month": ">1",
			},
		},
		{
			ID:          uuid.MustParse("ccccccc3-0000-0000-0000-000000000002"),
			Type:        asset.AssetTypeAudience,
			Description: "Audience: F 24-35, >3h/day social, >1 purchase last month",
			Payload: map[string]any{
				"gender":          "Female",
				"age_group":       "24-35",
				"country":         "UK",
				"social_hours":    ">3",
				"purchases_month": ">1",
			},
		},
	}

	for _, a := range assetsToSeed {
		if _, err := tx.Asset.
			Create().
			SetID(a.ID).
			SetAssetType(a.Type).
			SetDescription(a.Description).
			SetPayload(a.Payload).
			SetCreatedAt(now).
			Save(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
