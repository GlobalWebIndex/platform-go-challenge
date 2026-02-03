package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	tc_redis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"go-favorites-app/internal/adapter/api/rest"
	adapter_redis "go-favorites-app/internal/adapter/cache/redis"
	repo "go-favorites-app/internal/adapter/storage/postgres"
	"go-favorites-app/internal/core/domain/favorites"
	"go-favorites-app/internal/core/service"
)

func TestHTTPIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// --- 1. Infrastructure Setup (Postgres + Redis) ---

	// Postgres
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

	// Redis
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

	// Connect to DB
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

	// --- 2. Application Wiring ---

	// Redis Adapter
	redisUrl := redisConnStr
	if len(redisUrl) > 8 && redisUrl[:8] == "redis://" {
		redisUrl = redisUrl[8:]
	}
	cache := adapter_redis.NewAdapter(redisUrl)

	// User Service
	userRepo := repo.NewUserRepository(dbPool)
	jwtSecret := "test-secret"
	authService := service.NewAuthService(userRepo, jwtSecret)

	// Favorite Service
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	favRepo := repo.NewRepository(dbPool)
	favService := service.NewService(favRepo, cache, &NoOpEnricher{}, logger)

	// Handlers
	authHandler := rest.NewAuthHandler(authService)
	favHandler := rest.NewHandler(favService, logger)

	// Router
	handler := rest.NewRouter(favHandler, authHandler, jwtSecret)
	server := httptest.NewServer(handler)
	defer server.Close()

	client := server.Client()

	// --- 3. Test Cases ---

	// Helper to authenticate
	authenticate := func(email, password string) string {
		// SignUp
		signUpBody := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)
		resp, err := client.Post(server.URL+"/signup", "application/json", bytes.NewBufferString(signUpBody))
		if err != nil {
			t.Fatalf("SignUp failed: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("SignUp failed status: %d body: %s", resp.StatusCode, body)
		}

		// Login
		loginBody := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)
		resp, err = client.Post(server.URL+"/login", "application/json", bytes.NewBufferString(loginBody))
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Login failed status: %d", resp.StatusCode)
		}

		var res map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}
		return res["token"]
	}

	// Helper to create asset
	createAsset := func(token string, name string) {
		assetBody := fmt.Sprintf(`{
			"type": "chart",
			"name": "%s",
			"x_axis": "time",
			"y_axis": "val"
		}`, name)

		req, _ := http.NewRequest("POST", server.URL+"/favorites", bytes.NewBufferString(assetBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Create failed status: %d body: %s", resp.StatusCode, body)
		}
	}

	// Helper to list assets
	listAssets := func(token string) []favorites.BaseAsset {
		req, _ := http.NewRequest("GET", server.URL+"/favorites", nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			return nil
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("List failed status: %d body: %s", resp.StatusCode, body)
		}

		var bases []favorites.BaseAsset
		dec := json.NewDecoder(resp.Body)
		for {
			var item struct {
				Name string `json:"name"`
			}
			if err := dec.Decode(&item); err != nil {
				if err == io.EOF {
					break
				}
				t.Fatalf("Failed to decode NDJSON line: %v", err)
			}
			bases = append(bases, favorites.BaseAsset{Name: item.Name})
		}
		return bases
	}

	t.Run("Multi-tenant Isolation", func(t *testing.T) {
		tokenA := authenticate("userA@example.com", "passA")
		tokenB := authenticate("userB@example.com", "passB")

		// User A creates 2 assets
		createAsset(tokenA, "Asset A1")
		createAsset(tokenA, "Asset A2")

		// User B creates 1 asset
		createAsset(tokenB, "Asset B1")

		// Verify A sees only A's assets
		assetsA := listAssets(tokenA)
		if len(assetsA) != 2 {
			t.Errorf("Expected User A to see 2 assets, saw %d", len(assetsA))
		}
		foundA1 := false
		foundA2 := false
		for _, a := range assetsA {
			if a.Name == "Asset B1" {
				t.Error("User A shouldn't see Asset B1")
			}
			if a.Name == "Asset A1" {
				foundA1 = true
			}
			if a.Name == "Asset A2" {
				foundA2 = true
			}
		}
		if !foundA1 || !foundA2 {
			t.Errorf("User A missing expected assets. Saw: %v", assetsA)
		}

		// Verify B sees only B's assets
		assetsB := listAssets(tokenB)
		if len(assetsB) != 1 {
			t.Errorf("Expected User B to see 1 asset, saw %d", len(assetsB))
		}
		if len(assetsB) > 0 && assetsB[0].Name != "Asset B1" {
			t.Errorf("Expected User B to see Asset B1, saw %s", assetsB[0].Name)
		}
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		// No token
		req, _ := http.NewRequest("GET", server.URL+"/favorites", nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected 401 without token, got %d", resp.StatusCode)
		}

		// Bad token
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected 401 with bad token, got %d", resp.StatusCode)
		}
	})
}
