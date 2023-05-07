package repositories

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"testing"
)

func TestUserRepository_DeleteUser(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uid uint32
	}
	tests := []struct {
		name                           string
		args                           args
		mockSqlSelectUserQueryExpected string
		mockSqlDeleteUserQueryExpected string
		mockUserReturned               *User
		mockDeletedUserIdReturned      int
	}{
		{
			name: "valid",
			args: args{
				uid: 666,
			},
			mockSqlSelectUserQueryExpected: `SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL LIMIT 1`,
			mockUserReturned: &User{
				Model:    gorm.Model{ID: 666},
				Username: "mockUsername",
				Password: "mockPassword",
				FavouriteCharts: []Chart{
					{
						Model: gorm.Model{ID: 1},
						Title: "title1",
					},
					{
						Model: gorm.Model{ID: 2},
						Title: "title2",
					},
				},
				FavouriteInsights: []Insight{
					{
						Model: gorm.Model{ID: 1},
						Text:  "text1",
					},
					{
						Model: gorm.Model{ID: 2},
						Text:  "text2",
					},
				},
				FavouriteAudiences: []Audience{
					{
						Model:  gorm.Model{ID: 1},
						Gender: "gender1",
					},
					{
						Model:  gorm.Model{ID: 2},
						Gender: "gender2",
					},
				},
			},
			// The query is UPDATE as the deletion is SOFT in the gorm action I've set
			mockSqlDeleteUserQueryExpected: `UPDATE "users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`,
			mockDeletedUserIdReturned:      666,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectUserQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "username", "password"},
					).AddRow(
						tt.mockUserReturned.ID, tt.mockUserReturned.Username, tt.mockUserReturned.Password,
					),
				)
			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlDeleteUserQueryExpected)).
				WithArgs(sqlmock.AnyArg(), tt.args.uid).
				WillReturnResult(sqlmock.NewResult(int64(tt.mockDeletedUserIdReturned), 1))
			mockDb.ExpectCommit()

			err := repo.DeleteUser(context.Background(), tt.args.uid)

			assert.NoError(t, err)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

// The delete action in the gorm is done by selecting the requested description and then
// soft delete it, so we have two possible error sources: the select and the update
func TestUserRepository_DeleteUserWithSelectError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uid uint32
	}
	tests := []struct {
		name                           string
		args                           args
		mockSqlSelectUserQueryExpected string
		mockSqlErrorReturned           error
		expectedError                  error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlSelectUserQueryExpected: `SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:           errors.New("random error"),
			expectedError:                  errors.New("random error"),
		},
		{
			name: "description not found",
			args: args{
				uid: 666,
			},
			mockSqlSelectUserQueryExpected: `SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:           gorm.ErrRecordNotFound,
			expectedError: apierrors.UserNotFoundErrorWrapper{
				ReturnedStatusCode: http.StatusNotFound,
				OriginalError:      errors.New("userId 666 not found"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectUserQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)

			actual := repo.DeleteUser(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestUserRepository_DeleteUserWithDeleteError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uid uint32
	}
	tests := []struct {
		name                           string
		args                           args
		mockSqlSelectUserQueryExpected string
		mockSqlDeleteUserQueryExpected string
		mockUserReturned               *User
		mockSqlErrorReturned           error
		expectedError                  error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlSelectUserQueryExpected: `SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL LIMIT 1`,
			mockUserReturned: &User{
				Model:    gorm.Model{ID: 666},
				Username: "mockUsername",
				Password: "mockPassword",
				FavouriteCharts: []Chart{
					{
						Model: gorm.Model{ID: 1},
						Title: "title1",
					},
					{
						Model: gorm.Model{ID: 2},
						Title: "title2",
					},
				},
				FavouriteInsights: []Insight{
					{
						Model: gorm.Model{ID: 1},
						Text:  "text1",
					},
					{
						Model: gorm.Model{ID: 2},
						Text:  "text2",
					},
				},
				FavouriteAudiences: []Audience{
					{
						Model:  gorm.Model{ID: 1},
						Gender: "gender1",
					},
					{
						Model:  gorm.Model{ID: 2},
						Gender: "gender2",
					},
				},
			},
			// The query is UPDATE as the deletion is SOFT in the gorm action I've set
			mockSqlDeleteUserQueryExpected: `UPDATE "users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`,
			mockSqlErrorReturned:           errors.New("random error"),
			expectedError:                  errors.New("random error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectUserQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "username", "password"},
					).AddRow(
						tt.mockUserReturned.ID, tt.mockUserReturned.Username, tt.mockUserReturned.Password,
					),
				)
			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlDeleteUserQueryExpected)).
				WithArgs(sqlmock.AnyArg(), tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)
			mockDb.ExpectRollback()

			actual := repo.DeleteUser(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
