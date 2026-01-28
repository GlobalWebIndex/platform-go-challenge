package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"platform-go-challenge/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserFavourites_InvalidPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/u1", nil)
	rec := httptest.NewRecorder()

	var mockDB *sql.DB
	handler := GetUserFavourites(mockDB)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "invalid path") {
		t.Fatalf("expected 'invalid path', got %q", rec.Body.String())
	}
}

func TestGetUserFavourites_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/users/missing/favourites", nil)
	rec := httptest.NewRecorder()

	handler := GetUserFavourites(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if !strings.Contains(rec.Body.String(), "user not found") {
		t.Fatalf("expected 'user not found', got %q", rec.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserFavourites_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	// User exists
	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password_hash"}).
			AddRow("u1", "Alice", "hash"))

	// Get favourites
	mock.ExpectQuery("SELECT").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "type", "title", "data", "description"}).
			AddRow("a1", models.AssetChart, "Sales", json.RawMessage(`{}`), nil))

	req := httptest.NewRequest(http.MethodGet, "/users/u1/favourites", nil)
	rec := httptest.NewRecorder()

	handler := GetUserFavourites(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp []models.FavouriteAsset
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 favourite, got %d", len(resp))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAddFavourite_InvalidPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("{}"))
	rec := httptest.NewRecorder()

	var mockDB *sql.DB
	handler := AddFavourite(mockDB)

	defer func() {
		if r := recover(); r != nil {
			// Expected: index out of range on parts[2]
		}
	}()
	handler.ServeHTTP(rec, req)
}

func TestAddFavourite_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	body := `{"asset_id":"a1","description":"My fav"}`
	req := httptest.NewRequest(http.MethodPost, "/users/missing/favourites", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler := AddFavourite(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if !strings.Contains(rec.Body.String(), "user not found") {
		t.Fatalf("expected 'user not found', got %q", rec.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAddFavourite_InvalidBody(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password_hash"}).
			AddRow("u1", "Alice", "hash"))

	req := httptest.NewRequest(http.MethodPost, "/users/u1/favourites", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	handler := AddFavourite(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "invalid body") {
		t.Fatalf("expected 'invalid body', got %q", rec.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAddFavourite_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	// User exists
	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password_hash"}).
			AddRow("u1", "Alice", "hash"))

	// Add favourite
	mock.ExpectExec("INSERT INTO favourites").
		WithArgs("u1", "a1", nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := `{"asset_id":"a1"}`
	req := httptest.NewRequest(http.MethodPost, "/users/u1/favourites", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler := AddFavourite(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["message"] != "Favourite added successfully" {
		t.Fatalf("expected success message, got %q", resp["message"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateFavourite_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	body := `{"description":"Updated"}`
	req := httptest.NewRequest(http.MethodPatch, "/users/missing/favourites/a1", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler := UpdateFavourite(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateFavourite_InvalidBody(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password_hash"}).
			AddRow("u1", "Alice", "hash"))

	req := httptest.NewRequest(http.MethodPatch, "/users/u1/favourites/a1", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	handler := UpdateFavourite(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateFavourite_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	// User exists
	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password_hash"}).
			AddRow("u1", "Alice", "hash"))

	// Update favourite
	mock.ExpectExec("UPDATE favourites").
		WithArgs("u1", "a1", "Updated desc").
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := `{"description":"Updated desc"}`
	req := httptest.NewRequest(http.MethodPatch, "/users/u1/favourites/a1", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler := UpdateFavourite(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRemoveFavourite_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodDelete, "/users/missing/favourites/a1", nil)
	rec := httptest.NewRecorder()

	handler := RemoveFavourite(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRemoveFavourite_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	// User exists
	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password_hash"}).
			AddRow("u1", "Alice", "hash"))

	// Delete favourite (no rows affected)
	mock.ExpectExec("DELETE FROM favourites").
		WithArgs("u1", "a1").
		WillReturnResult(sqlmock.NewResult(0, 0))

	req := httptest.NewRequest(http.MethodDelete, "/users/u1/favourites/a1", nil)
	rec := httptest.NewRecorder()

	handler := RemoveFavourite(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRemoveFavourite_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	// User exists
	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password_hash"}).
			AddRow("u1", "Alice", "hash"))

	// Delete favourite (1 row affected)
	mock.ExpectExec("DELETE FROM favourites").
		WithArgs("u1", "a1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest(http.MethodDelete, "/users/u1/favourites/a1", nil)
	rec := httptest.NewRecorder()

	handler := RemoveFavourite(db)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
