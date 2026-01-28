package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"platform-go-challenge/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserFavourites_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "type", "title", "data", "description"}).
		AddRow("a1", models.AssetChart, ptrString("Sales"), json.RawMessage(`{}`), nil)

	mock.ExpectQuery("SELECT").
		WithArgs("u1").
		WillReturnRows(rows)

	favs, err := GetUserFavourites(context.Background(), db, "u1")
	if err != nil {
		t.Fatalf("GetUserFavourites error: %v", err)
	}
	if len(favs) != 1 {
		t.Fatalf("expected 1 favourite, got %d", len(favs))
	}
	if favs[0].AssetID != "a1" {
		t.Fatalf("unexpected favourite: %+v", favs[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserFavourites_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "type", "title", "data", "description"})

	mock.ExpectQuery("SELECT").
		WithArgs("u1").
		WillReturnRows(rows)

	favs, err := GetUserFavourites(context.Background(), db, "u1")
	if err != nil {
		t.Fatalf("GetUserFavourites error: %v", err)
	}
	if len(favs) != 0 {
		t.Fatalf("expected empty list, got %d", len(favs))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserFavourites_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").
		WithArgs("u1").
		WillReturnError(sql.ErrConnDone)

	favs, err := GetUserFavourites(context.Background(), db, "u1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if favs != nil {
		t.Fatalf("expected nil, got %+v", favs)
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

	mock.ExpectExec("INSERT INTO favourites").
		WithArgs("u1", "a1", nil).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = AddFavourite(context.Background(), db, "u1", "a1", nil)
	if err != nil {
		t.Fatalf("AddFavourite error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAddFavourite_WithDescription(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	desc := ptrString("My favourite")
	mock.ExpectExec("INSERT INTO favourites").
		WithArgs("u1", "a1", desc).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = AddFavourite(context.Background(), db, "u1", "a1", desc)
	if err != nil {
		t.Fatalf("AddFavourite error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAddFavourite_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO favourites").
		WithArgs("u1", "a1", nil).
		WillReturnError(sql.ErrConnDone)

	err = AddFavourite(context.Background(), db, "u1", "a1", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
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

	mock.ExpectExec("DELETE FROM favourites").
		WithArgs("u1", "a1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = RemoveFavourite(context.Background(), db, "u1", "a1")
	if err != nil {
		t.Fatalf("RemoveFavourite error: %v", err)
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

	mock.ExpectExec("DELETE FROM favourites").
		WithArgs("u1", "missing").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = RemoveFavourite(context.Background(), db, "u1", "missing")
	if err == nil {
		t.Fatal("expected error for missing favourite")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRemoveFavourite_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM favourites").
		WithArgs("u1", "a1").
		WillReturnError(sql.ErrConnDone)

	err = RemoveFavourite(context.Background(), db, "u1", "a1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateFavouriteDescription_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	desc := ptrString("Updated description")
	mock.ExpectExec("UPDATE favourites").
		WithArgs("u1", "a1", desc).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = UpdateFavouriteDescription(context.Background(), db, "u1", "a1", desc)
	if err != nil {
		t.Fatalf("UpdateFavouriteDescription error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateFavouriteDescription_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	desc := ptrString("Updated")
	mock.ExpectExec("UPDATE favourites").
		WithArgs("u1", "missing", desc).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = UpdateFavouriteDescription(context.Background(), db, "u1", "missing", desc)
	if err == nil {
		t.Fatal("expected error for missing favourite")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateFavouriteDescription_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	desc := ptrString("Updated")
	mock.ExpectExec("UPDATE favourites").
		WithArgs("u1", "a1", desc).
		WillReturnError(sql.ErrConnDone)

	err = UpdateFavouriteDescription(context.Background(), db, "u1", "a1", desc)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
