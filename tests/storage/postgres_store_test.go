
package storage_test

import (
	"assetsApp/internal/models"
	"assetsApp/internal/storage"
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pool  *pgxpool.Pool
	store *storage.PostgresStore
)

func TestMain(m *testing.M) {
	// setup
	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:13"),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute),
		),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	pool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	store = storage.NewPostgresStore(pool)

	// Create tables
	if err := createTables(); err != nil {
		log.Fatalf("failed to create tables: %s", err)
	}

	// run tests
	exitCode := m.Run()

	// teardown
	pool.Close()
	if err := pgContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	os.Exit(exitCode)
}

func createTables() error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		CREATE TABLE users (
			id UUID PRIMARY KEY,
			name VARCHAR(255)
		);
		CREATE TABLE assets (
			asset_id VARCHAR(255) PRIMARY KEY,
			title VARCHAR(255),
			description TEXT,
			asset_type VARCHAR(50),
			user_id UUID REFERENCES users(id)
		);
		CREATE TABLE charts (
			id VARCHAR(255) PRIMARY KEY REFERENCES assets(asset_id),
			title VARCHAR(255),
			description TEXT,
			x_axis_title VARCHAR(255),
			y_axis_title VARCHAR(255)
		);
		CREATE TABLE chart_data (
			chart_id VARCHAR(255) REFERENCES charts(id),
			datapoint_code VARCHAR(255),
			value FLOAT,
			PRIMARY KEY (chart_id, datapoint_code)
		);
		CREATE TABLE insights (
			id VARCHAR(255) PRIMARY KEY REFERENCES assets(asset_id),
			description TEXT
		);
		CREATE TABLE audiences (
			id VARCHAR(255) PRIMARY KEY REFERENCES assets(asset_id),
			gender VARCHAR(50),
			country VARCHAR(255),
			age_group VARCHAR(50),
			social_hours INT,
			purchases INT,
			description TEXT
		);
		CREATE TABLE favourites (
			user_id UUID REFERENCES users(id),
			asset_id VARCHAR(255) REFERENCES assets(asset_id),
			asset_type VARCHAR(50),
			PRIMARY KEY (user_id, asset_id)
		);
	`)
	return err
}

func cleanup() {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		DELETE FROM favourites;
		DELETE FROM chart_data;
		DELETE FROM charts;
		DELETE FROM insights;
		DELETE FROM audiences;
		DELETE FROM assets;
		DELETE FROM users;
	`)
	if err != nil {
		log.Fatalf("failed to cleanup database: %s", err)
	}
}

func TestPostgresStore_AddAndGetAsset(t *testing.T) {
	defer cleanup()

	userID := uuid.New()
	chart := &models.Chart{
		ID:          "chart1",
		Title:       "Test Chart",
		Description: "Test Description",
		XAxisTitle:  "X",
		YAxisTitle:  "Y",
		Data: []models.ChartData{
			{DatapointCode: "dp1", Value: 10.5},
		},
	}

	store.Add(userID, chart)

	assets := store.Get(userID)
	assert.Len(t, assets, 1)
	assert.Equal(t, chart.ID, assets[0].GetID())
}

func TestPostgresStore_RemoveAsset(t *testing.T) {
	defer cleanup()

	userID := uuid.New()
	chart := &models.Chart{ID: "chart1"}

	store.Add(userID, chart)
	assert.Len(t, store.Get(userID), 1)

	removed := store.Remove(userID, "chart1")
	assert.True(t, removed)
	assert.Len(t, store.Get(userID), 0)
}

func TestPostgresStore_EditDescription(t *testing.T) {
	defer cleanup()

	userID := uuid.New()
	chart := &models.Chart{ID: "chart1", Description: "Old Description"}

	store.Add(userID, chart)

	edited := store.EditDescription(userID, "chart1", "New Description")
	assert.True(t, edited)

	assets := store.Get(userID)
	assert.Len(t, assets, 1)

	retrievedChart, ok := assets[0].(*models.Chart)
	assert.True(t, ok)
	assert.Equal(t, "New Description", retrievedChart.Description)
}

func TestPostgresStore_AddAndGetFavourites(t *testing.T) {
	defer cleanup()

	userID := uuid.New()
	chart := &models.Chart{
		ID:          "chart1",
		Title:       "Test Chart",
		Description: "Test Description",
	}
	store.Add(userID, chart)

	added := store.AddFavourite(userID, "chart1", "chart")
	assert.True(t, added)

	favourites := store.GetFavourites(userID)
	assert.Len(t, favourites, 1)
	assert.Equal(t, "chart1", favourites[0].Asset.GetID())
}

func TestPostgresStore_RemoveFavourite(t *testing.T) {
	defer cleanup()

	userID := uuid.New()
	chart := &models.Chart{ID: "chart1"}
	store.Add(userID, chart)
	store.AddFavourite(userID, "chart1", "chart")
	assert.Len(t, store.GetFavourites(userID), 1)

	removed := store.RemoveFavourite(userID, "chart1")
	assert.True(t, removed)
	assert.Len(t, store.GetFavourites(userID), 0)
}
