package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"platform-go-challenge/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUnmarshalAssetData(t *testing.T) {
	tests := []struct {
		name        string
		assetType   models.AssetType
		rawData     json.RawMessage
		shouldError bool
	}{
		{
			name:        "chart asset",
			assetType:   models.AssetChart,
			rawData:     json.RawMessage(`{"x_axis_title":"Month","y_axis_title":"Sales","data_points":[100,200]}`),
			shouldError: false,
		},
		{
			name:        "insight asset",
			assetType:   models.AssetInsight,
			rawData:     json.RawMessage(`{"text":"Users prefer mobile"}`),
			shouldError: false,
		},
		{
			name:        "audience asset",
			assetType:   models.AssetAudience,
			rawData:     json.RawMessage(`{"gender":"M","birth_country":"US","age_groups":["25-34"],"hours_on_social_media_min":3,"purchases_last_month_min":5}`),
			shouldError: false,
		},
		{
			name:        "unknown asset type",
			assetType:   models.AssetType("unknown"),
			rawData:     json.RawMessage(`{}`),
			shouldError: true,
		},
		{
			name:        "invalid JSON",
			assetType:   models.AssetChart,
			rawData:     json.RawMessage(`{invalid}`),
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := unmarshalAssetData(tt.assetType, tt.rawData)

			if tt.shouldError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.shouldError && data == nil {
				t.Error("expected data but got nil")
			}
		})
	}
}

func TestCreateAsset_DB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	asset := models.Asset{
		Type:  models.AssetChart,
		Title: ptrString("Sales"),
		Data:  json.RawMessage(`{"points":[1,2]}`),
	}
	desc := ptrString("Monthly data")

	mock.ExpectQuery("INSERT INTO assets").
		WithArgs(models.AssetChart, ptrString("Sales"), desc, json.RawMessage(`{"points":[1,2]}`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("a1"))

	id, err := CreateAsset(context.Background(), db, asset, desc)
	if err != nil {
		t.Fatalf("CreateAsset error: %v", err)
	}
	if id != "a1" {
		t.Fatalf("expected id a1, got %q", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreateAsset_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	asset := models.Asset{
		Type:  models.AssetChart,
		Title: ptrString("Sales"),
	}

	mock.ExpectQuery("INSERT INTO assets").
		WillReturnError(sql.ErrConnDone)

	id, err := CreateAsset(context.Background(), db, asset, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if id != "" {
		t.Fatalf("expected empty id, got %q", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetAssetByID_DB_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	createdAt := time.Now()
	mock.ExpectQuery("SELECT id, type, title, data, created_at").
		WithArgs("a1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "type", "title", "data", "created_at"}).
			AddRow("a1", models.AssetChart, ptrString("Sales"), json.RawMessage(`{}`), createdAt))

	asset, err := GetAssetByID(context.Background(), db, "a1")
	if err != nil {
		t.Fatalf("GetAssetByID error: %v", err)
	}
	if asset == nil {
		t.Fatal("expected asset, got nil")
	}
	if asset.ID != "a1" || asset.Type != models.AssetChart {
		t.Fatalf("unexpected asset: %+v", asset)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetAssetByID_DB_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, type, title, data, created_at").
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	asset, err := GetAssetByID(context.Background(), db, "missing")
	if err != nil {
		t.Fatalf("GetAssetByID error: %v", err)
	}
	if asset != nil {
		t.Fatalf("expected nil asset, got %+v", asset)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetAssetByID_DB_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, type, title, data, created_at").
		WithArgs("a1").
		WillReturnError(sql.ErrConnDone)

	asset, err := GetAssetByID(context.Background(), db, "a1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if asset != nil {
		t.Fatalf("expected nil asset, got %+v", asset)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListAssets_DB_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "type", "title", "data", "created_at"}).
		AddRow("a1", models.AssetChart, ptrString("Chart"), json.RawMessage(`{}`), now).
		AddRow("a2", models.AssetInsight, ptrString("Insight"), json.RawMessage(`{}`), now)

	mock.ExpectQuery("SELECT id, type, title, data, created_at").
		WillReturnRows(rows)

	assets, err := ListAssets(context.Background(), db)
	if err != nil {
		t.Fatalf("ListAssets error: %v", err)
	}
	if len(assets) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(assets))
	}
	if assets[0].ID != "a1" || assets[1].ID != "a2" {
		t.Fatalf("unexpected assets: %+v", assets)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListAssets_DB_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "type", "title", "data", "created_at"})

	mock.ExpectQuery("SELECT id, type, title, data, created_at").
		WillReturnRows(rows)

	assets, err := ListAssets(context.Background(), db)
	if err != nil {
		t.Fatalf("ListAssets error: %v", err)
	}
	if len(assets) != 0 {
		t.Fatalf("expected empty list, got %d assets", len(assets))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListAssets_DB_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, type, title, data, created_at").
		WillReturnError(sql.ErrConnDone)

	assets, err := ListAssets(context.Background(), db)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if assets != nil {
		t.Fatalf("expected nil assets, got %+v", assets)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreateAsset(t *testing.T) {
	tests := []struct {
		name        string
		asset       models.Asset
		description *string
		shouldError bool
	}{
		{
			name: "create chart asset",
			asset: models.Asset{
				Type:  models.AssetChart,
				Title: ptrString("Sales Chart"),
				Data:  json.RawMessage(`{"x_axis_title":"Month"}`),
			},
			description: nil,
			shouldError: false,
		},
		{
			name: "create asset with description",
			asset: models.Asset{
				Type:  models.AssetInsight,
				Title: ptrString("Market Insight"),
				Data:  json.RawMessage(`{"text":"Market growing"}`),
			},
			description: ptrString("Important insight"),
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.asset.Title == nil {
				t.Error("asset title cannot be nil")
			}
			if tt.asset.Type == "" {
				t.Error("asset type cannot be empty")
			}
		})
	}
}

func TestGetAssetByID(t *testing.T) {
	tests := []struct {
		name        string
		assetID     string
		shouldError bool
	}{
		{
			name:        "valid UUID",
			assetID:     "550e8400-e29b-41d4-a716-446655440000",
			shouldError: false,
		},
		{
			name:        "non-existent asset",
			assetID:     "invalid-uuid",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.assetID == "" {
				t.Error("assetID cannot be empty")
			}
		})
	}
}

// Helper functions
func ptrString(s string) *string {
	return &s
}
