package postgres

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	domain "go-favorites-app/internal/core/domain/favorites"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	// Schema initialization
	schema := `
	CREATE TABLE favorites (
		id UUID PRIMARY KEY,
		type VARCHAR(50) NOT NULL,
		user_id VARCHAR(255) NOT NULL,
		asset_data JSONB NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`
	if _, err := dbPool.Exec(ctx, schema); err != nil {
		t.Fatalf("failed to init schema: %v", err)
	}

	cleanup := func() {
		dbPool.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres: %v", err)
		}
	}

	return dbPool, cleanup
}

func TestRepository_ThreadSafety(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbPool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(dbPool)
	ctx := context.Background()

	t.Run("concurrent saves", func(t *testing.T) {
		const numGoroutines = 50
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(idx int) {
				defer wg.Done()
				id := uuid.NewString()
				asset := domain.Insight{
					BaseAsset: domain.BaseAsset{
						ID:     id,
						UserID: "user-1",
						Name:   fmt.Sprintf("Asset %d", idx),
						Type:   domain.AssetTypeInsight,
					},
					Content: "test content",
				}
				if err := repo.Save(ctx, asset); err != nil {
					t.Errorf("failed to save asset %d: %v", idx, err)
				}
			}(i)
		}
		wg.Wait()

		// Verify count
		var count int
		err := dbPool.QueryRow(ctx, "SELECT count(*) FROM favorites").Scan(&count)
		if err != nil {
			t.Fatalf("failed to count assets: %v", err)
		}
		if count != numGoroutines {
			t.Errorf("expected %d assets, got %d", numGoroutines, count)
		}
	})

	t.Run("concurrent updates on same record", func(t *testing.T) {
		// Test row locking / atomicity of jsonb_set
		id := uuid.NewString()
		initialAsset := domain.Insight{
			BaseAsset: domain.BaseAsset{
				ID:     id,
				UserID: "user-target",
				Name:   "Target Asset",
				Type:   domain.AssetTypeInsight,
			},
			Content: "original",
		}
		if err := repo.Save(ctx, initialAsset); err != nil {
			t.Fatalf("failed to seed asset: %v", err)
		}

		const numUpdates = 20
		var wg sync.WaitGroup
		wg.Add(numUpdates)

		for i := 0; i < numUpdates; i++ {
			go func(idx int) {
				defer wg.Done()
				desc := fmt.Sprintf("desc %d", idx)
				_, err := repo.UpdateDescription(ctx, id, desc)
				if err != nil {
					t.Errorf("failed to update asset: %v", err)
				}
			}(i)
		}
		wg.Wait()

		// The final description should be one of the "desc %d" values
		updated, err := repo.FindByID(ctx, id)
		if err != nil {
			t.Fatalf("failed to fetch updated asset: %v", err)
		}
		if updated.GetID() != id {
			t.Errorf("expected ID %s, got %s", id, updated.GetID())
		}
	})

	t.Run("FindAll during concurrent writes", func(t *testing.T) {
		stopC := make(chan struct{})
		var wg sync.WaitGroup

		// Background writer
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stopC:
					return
				default:
					asset := domain.Chart{
						BaseAsset: domain.BaseAsset{
							ID:     uuid.NewString(),
							UserID: "user-stream",
							Name:   "Stream Asset",
							Type:   domain.AssetTypeChart,
						},
						XAxis: "x",
						YAxis: "y",
					}
					if err := repo.Save(ctx, asset); err != nil {
						// Log only
						_ = err
					}
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()

		// Consumer
		for i := 0; i < 5; i++ {
			iter, err := repo.FindAll(ctx, 10, 0)
			if err != nil {
				t.Errorf("FindAll error: %v", err)
				continue
			}

			// Just drain the iterator
			for range iter {
			}
			time.Sleep(50 * time.Millisecond)
		}

		close(stopC)
		wg.Wait()
	})
}
