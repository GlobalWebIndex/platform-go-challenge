package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"platform-go-challenge/models"
)

func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	user := models.User{
		ID:           "u1",
		Name:         "Alice",
		PasswordHash: "hashed_pwd",
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("u1", "Alice", "hashed_pwd").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("u1"))

	id, err := CreateUser(context.Background(), db, user)
	if err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}
	if id != "u1" {
		t.Fatalf("expected id u1, got %q", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreateUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	user := models.User{ID: "u1", Name: "Alice"}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("u1", "Alice", "").
		WillReturnError(sql.ErrConnDone)

	id, err := CreateUser(context.Background(), db, user)
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

func TestGetUserByID_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password_hash"}).
			AddRow("u1", "Alice", "hashed"))

	user, err := GetUserByID(context.Background(), db, "u1")
	if err != nil {
		t.Fatalf("GetUserByID error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != "u1" || user.Name != "Alice" {
		t.Fatalf("unexpected user: %+v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	user, err := GetUserByID(context.Background(), db, "missing")
	if err != nil {
		t.Fatalf("GetUserByID error: %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, password_hash").
		WithArgs("u1").
		WillReturnError(sql.ErrConnDone)

	user, err := GetUserByID(context.Background(), db, "u1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListUsers_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "password_hash"}).
		AddRow("u1", "Alice", "hash1").
		AddRow("u2", "Bob", "hash2")

	mock.ExpectQuery("SELECT id, name, password_hash").
		WillReturnRows(rows)

	users, err := ListUsers(context.Background(), db)
	if err != nil {
		t.Fatalf("ListUsers error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].ID != "u1" || users[1].ID != "u2" {
		t.Fatalf("unexpected users: %+v", users)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListUsers_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "password_hash"})

	mock.ExpectQuery("SELECT id, name, password_hash").
		WillReturnRows(rows)

	users, err := ListUsers(context.Background(), db)
	if err != nil {
		t.Fatalf("ListUsers error: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected empty list, got %d users", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListUsers_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, password_hash").
		WillReturnError(sql.ErrConnDone)

	users, err := ListUsers(context.Background(), db)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if users != nil {
		t.Fatalf("expected nil users, got %+v", users)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
