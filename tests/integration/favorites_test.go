package integration

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	tc_redis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	adapter_redis "go-favorites-app/internal/adapter/cache/redis"
	repo "go-favorites-app/internal/adapter/storage/postgres"
	"go-favorites-app/internal/core/domain/favorites"
	"go-favorites-app/internal/core/service"
)

// NoOpEnricher mock
type NoOpEnricher struct{}

func (e *NoOpEnricher) Enrich(ctx context.Context, asset favorites.Asset) error {
	return nil
}

func TestFavoritesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// 1. Start Postgres Container
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres: %v", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate postgres: %v", err)
		}
	}()

	pgConnStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get pg connection string: %v", err)
	}

	// 2. Start Redis Container
	// Use the alias tc_redis to avoid conflict with local redis package
	redisContainer, err := tc_redis.Run(ctx,
		"redis:7-alpine",
	)
	if err != nil {
		t.Fatalf("failed to start redis: %v", err)
	}
	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate redis: %v", err)
		}
	}()

	redisConnStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get redis connection string: %v", err)
	}

	// 3. Setup Dependencies
	dbPool, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}
	defer dbPool.Close()

	// Init Schema
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS favorites (
		id UUID PRIMARY KEY,
		type VARCHAR(50) NOT NULL,
		asset_data JSONB NOT NULL,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_favorites_asset_data ON favorites USING GIN (asset_data);
	CREATE INDEX IF NOT EXISTS idx_favorites_type ON favorites (type);
	`
	if _, err := dbPool.Exec(ctx, schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create a test user
	testUserID := uuid.NewString()
	_, err = dbPool.Exec(ctx, "INSERT INTO users (id, email, password_hash) VALUES ($1, 'test@example.com', 'hash')", testUserID)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Init Adapters
	repository := repo.NewRepository(dbPool)

	// Clean up redis connection string for the adapter
	// The adapter uses redis.NewClient(&redis.Options{Addr: addr}).
	// If redisConnStr starts with "redis://", we might need to strip it depending on go-redis version logic,
	// but usually parsing is cleaner. For simplicity in this test environment:
	redisUrl := redisConnStr
	if len(redisUrl) > 8 && redisUrl[:8] == "redis://" {
		redisUrl = redisUrl[8:]
	}

	cache := adapter_redis.NewAdapter(redisUrl)
	svc := service.NewService(repository, cache, &NoOpEnricher{}, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// 4. Seed Data
	totalAssets := 1000
	t.Logf("Seeding %d assets...", totalAssets)
	for i := 0; i < totalAssets; i++ {
		id := uuid.NewString()
		name := fmt.Sprintf("Asset %d", i)
		assetType := favorites.AssetTypeChart
		if i%3 == 1 {
			assetType = favorites.AssetTypeInsight
		} else if i%3 == 2 {
			assetType = favorites.AssetTypeAudience
		}

		// Create concrete type
		var asset favorites.Asset
		base := favorites.BaseAsset{ID: id, Name: name, Type: assetType, UserID: testUserID}

		switch assetType {
		case favorites.AssetTypeChart:
			asset = favorites.Chart{BaseAsset: base, XAxis: "time", YAxis: "value"}
		case favorites.AssetTypeInsight:
			asset = favorites.Insight{BaseAsset: base, Content: "insight text"}
		case favorites.AssetTypeAudience:
			asset = favorites.Audience{
				BaseAsset: base,
				Rules: favorites.AudienceRules{
					Gender:  "all",
					Country: "US",
					AgeMin:  18,
					AgeMax:  99,
				},
			}
		}

		if err := svc.Save(ctx, asset); err != nil {
			t.Fatalf("failed to save asset %d: %v", i, err)
		}
	}

	// 5. Verify FindAll with Iterator
	t.Log("Verifying FindAll iterator...")
	limit := 100
	offset := 0
	count := 0

	for count < totalAssets {
		iter, err := svc.FindAll(ctx, limit, offset)
		if err != nil {
			t.Fatalf("FindAll failed: %v", err)
		}

		pageCount := 0
		for _, err := range iter {
			if err != nil {
				t.Fatalf("Iterator error: %v", err)
			}
			pageCount++
			count++
		}

		if pageCount == 0 {
			break
		}
		offset += limit
	}

	if count != totalAssets {
		t.Errorf("Expected %d assets, got %d", totalAssets, count)
	} else {
		t.Logf("Successfully retrieved all %d assets via iterator/paging", count)
	}
}
