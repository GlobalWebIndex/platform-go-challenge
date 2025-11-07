package integration_test

import (
	"assetsApp/internal/handlers"
	assetServices "assetsApp/internal/services/asset"
	favouriteServices "assetsApp/internal/services/favourite"
	"assetsApp/internal/storage"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pool   *pgxpool.Pool
	store  *storage.PostgresStore
	router *mux.Router
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

	// Set up router
	assetService := assetServices.NewAssetService(store)
	favouriteService := favouriteServices.NewFavouriteService(store)
	assetHandler := handlers.NewAssetHandler(assetService)
	favouriteHandler := handlers.NewFavouriteHandler(favouriteService)

	router = mux.NewRouter()
	router.HandleFunc("/users/{userId}/assets", assetHandler.GetAssets).Methods("GET")
	router.HandleFunc("/users/{userId}/assets", assetHandler.AddAsset).Methods("POST")
	router.HandleFunc("/users/{userId}/assets/{assetId}", assetHandler.RemoveAsset).Methods("DELETE")
	router.HandleFunc("/users/{userId}/favourites", favouriteHandler.GetFavourites).Methods("GET")
	router.HandleFunc("/users/{userId}/favourites/{assetId}", favouriteHandler.AddFavourite).Methods("POST")
	router.HandleFunc(
		"/users/{userId}/favourites/{assetId}",
		favouriteHandler.RemoveFavourite,
	).Methods("DELETE")

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

func TestIntegration_AssetAndFavouriteFlow(t *testing.T) {
	defer cleanup()

	userID := uuid.New()
	assetID := "chart1"

	// 1. Add an asset
	assetData := map[string]interface{}{
		"type":         "chart",
		"id":           assetID,
		"title":        "Test Chart",
		"description":  "Test Description",
		"x_axis_title": "X-Axis",
		"y_axis_title": "Y-Axis",
		"data":         []map[string]interface{}{},
	}
	assetPayload, _ := json.Marshal(assetData)
	req := httptest.NewRequest("POST", "/users/"+userID.String()+"/assets", bytes.NewBuffer(assetPayload))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Add asset to favourites
	favouritePayload := []byte(`{"asset_type": "chart"}`)
	req = httptest.NewRequest("POST", "/users/"+userID.String()+"/favourites/"+assetID, bytes.NewBuffer(favouritePayload))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 3. Get favourites and verify
	req = httptest.NewRequest("GET", "/users/"+userID.String()+"/favourites", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), assetID)

	// 4. Remove from favourites
	req = httptest.NewRequest("DELETE", "/users/"+userID.String()+"/favourites/"+assetID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. Get favourites and verify empty
	req = httptest.NewRequest("GET", "/users/"+userID.String()+"/favourites", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]\n", w.Body.String())

	// 6. Remove asset
	req = httptest.NewRequest("DELETE", "/users/"+userID.String()+"/assets/"+assetID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. Get assets and verify empty
	req = httptest.NewRequest("GET", "/users/"+userID.String()+"/assets", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "null\n", w.Body.String())
}
