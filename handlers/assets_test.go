package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"platform-go-challenge/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateAsset_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/assets", nil)
	rec := httptest.NewRecorder()

	var mockDB *sql.DB
	handler := CreateAsset(mockDB)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestCreateAsset_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/assets", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	var mockDB *sql.DB
	handler := CreateAsset(mockDB)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "Invalid request body") {
		t.Fatalf("expected invalid body message, got %q", rec.Body.String())
	}
}

func TestCreateAsset_MissingTitle(t *testing.T) {
	body := `{"type":"chart","description":"d","data":{}}`
	req := httptest.NewRequest(http.MethodPost, "/assets", strings.NewReader(body))
	rec := httptest.NewRecorder()

	var mockDB *sql.DB
	handler := CreateAsset(mockDB)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "Title is required") {
		t.Fatalf("expected title required message, got %q", rec.Body.String())
	}
}

func TestCreateAsset_MissingType(t *testing.T) {
	body := `{"title":"Asset","description":"d","data":{}}`
	req := httptest.NewRequest(http.MethodPost, "/assets", strings.NewReader(body))
	rec := httptest.NewRecorder()

	var mockDB *sql.DB
	handler := CreateAsset(mockDB)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "Type is required") {
		t.Fatalf("expected type required message, got %q", rec.Body.String())
	}
}

func TestCreateAsset_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	body := `{"type":"chart","title":"Sales","description":"Monthly","data":{"points":[1,2]}}`
	req := httptest.NewRequest(http.MethodPost, "/assets", strings.NewReader(body))
	rec := httptest.NewRecorder()

	mock.ExpectQuery("INSERT INTO assets").
		WithArgs(models.AssetChart, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("a1"))

	handler := CreateAsset(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["id"] != "a1" {
		t.Fatalf("expected id a1, got %v", resp["id"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetAsset_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, type, title, data, created_at").
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/assets/missing", nil)
	rec := httptest.NewRecorder()

	handler := GetAsset(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestGetAsset_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	createdAt := time.Now()
	mock.ExpectQuery("SELECT id, type, title, data, created_at").
		WithArgs("a1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "type", "title", "data", "created_at"}).
			AddRow("a1", models.AssetChart, "Sales", json.RawMessage(`{"x":1}`), createdAt))

	req := httptest.NewRequest(http.MethodGet, "/assets/a1", nil)
	rec := httptest.NewRecorder()

	handler := GetAsset(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp models.Asset
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.ID != "a1" || resp.Type != models.AssetChart {
		t.Fatalf("unexpected asset response: %+v", resp)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetAllAssets_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "type", "title", "data", "created_at"}).
		AddRow("a1", models.AssetChart, "Sales", json.RawMessage(`{"x":1}`), time.Now()).
		AddRow("a2", models.AssetInsight, "Insight", json.RawMessage(`{"text":"hi"}`), time.Now())

	mock.ExpectQuery("SELECT id, type, title, data, created_at").WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/assets", nil)
	rec := httptest.NewRecorder()

	handler := GetAllAssets(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []models.Asset
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(resp))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
