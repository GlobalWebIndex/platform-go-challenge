package database

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/giannis84/platform-go-challenge/internal/models"
	"github.com/lib/pq"
)

var testCols = []string{"id", "user_id", "asset_type", "description", "data", "created_at", "updated_at"}

func setupTestDB(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	DB = db
	return mock
}

func testChartJSON(id string) []byte {
	data, _ := json.Marshal(&models.Chart{ID: id, Title: "T", XAxisTitle: "X", YAxisTitle: "Y"})
	return data
}

// --- GetUserFavouritesFromDB ---

func TestGetUserFavouritesFromDB(t *testing.T) {
	now := time.Now()

	t.Run("returns favourites", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
			WithArgs("user1").
			WillReturnRows(sqlmock.NewRows(testCols).
				AddRow("c1", "user1", "chart", "desc", testChartJSON("c1"), now, now))

		favs, err := GetUserFavouritesFromDB("user1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(favs) != 1 {
			t.Fatalf("expected 1, got %d", len(favs))
		}
		if favs[0].ID != "c1" || favs[0].UserID != "user1" {
			t.Errorf("unexpected favourite: %+v", favs[0])
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("returns empty for unknown user", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
			WithArgs("unknown").
			WillReturnRows(sqlmock.NewRows(testCols))

		favs, err := GetUserFavouritesFromDB("unknown")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(favs) != 0 {
			t.Errorf("expected empty, got %d", len(favs))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("returns error on query failure", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
			WillReturnError(fmt.Errorf("connection failed"))

		_, err := GetUserFavouritesFromDB("user1")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// --- GetFavouriteFromDB ---

func TestGetFavouriteFromDB(t *testing.T) {
	now := time.Now()

	t.Run("returns favourite", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
			WithArgs("user1", "c1").
			WillReturnRows(sqlmock.NewRows(testCols).
				AddRow("c1", "user1", "chart", "desc", testChartJSON("c1"), now, now))

		fav, err := GetFavouriteFromDB("user1", "c1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fav.ID != "c1" || fav.Description != "desc" {
			t.Errorf("unexpected favourite: %+v", fav)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("returns ErrNotFound", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
			WithArgs("user1", "missing").
			WillReturnRows(sqlmock.NewRows(testCols))

		_, err := GetFavouriteFromDB("user1", "missing")
		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// --- AddFavouriteInDB ---

func TestAddFavouriteInDB(t *testing.T) {
	now := time.Now()
	fav := &models.FavouriteAsset{
		ID: "c1", UserID: "user1", AssetType: "chart",
		Description: "desc", CreatedAt: now, UpdatedAt: now,
		Data: &models.Chart{ID: "c1", Title: "T", XAxisTitle: "X", YAxisTitle: "Y"},
	}

	t.Run("inserts successfully", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectExec("INSERT INTO favourites").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := AddFavouriteInDB(context.Background(), fav)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("returns ErrAlreadyExists on unique violation", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectExec("INSERT INTO favourites").
			WillReturnError(&pq.Error{Code: "23505"})

		err := AddFavouriteInDB(context.Background(), fav)
		if err != ErrAlreadyExists {
			t.Errorf("expected ErrAlreadyExists, got: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("returns error on insert failure", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectExec("INSERT INTO favourites").
			WillReturnError(fmt.Errorf("connection failed"))

		err := AddFavouriteInDB(context.Background(), fav)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// --- UpdateFavouriteInDB ---

func TestUpdateFavouriteInDB(t *testing.T) {
	now := time.Now()
	fav := &models.FavouriteAsset{
		ID: "c1", UserID: "user1", AssetType: "chart",
		Description: "new desc", UpdatedAt: now,
		Data: &models.Chart{ID: "c1", Title: "T", XAxisTitle: "X", YAxisTitle: "Y"},
	}

	t.Run("updates successfully", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectExec("UPDATE favourites").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := UpdateFavouriteInDB(fav)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("returns ErrNotFound when no rows affected", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectExec("UPDATE favourites").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := UpdateFavouriteInDB(fav)
		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// --- DeleteFavouriteFromDB ---

func TestDeleteFavouriteFromDB(t *testing.T) {
	t.Run("deletes successfully", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectExec("DELETE FROM favourites").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := DeleteFavouriteFromDB("user1", "c1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("returns ErrNotFound when no rows affected", func(t *testing.T) {
		mock := setupTestDB(t)
		mock.ExpectExec("DELETE FROM favourites").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := DeleteFavouriteFromDB("user1", "missing")
		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}
