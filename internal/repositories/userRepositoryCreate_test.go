package repositories

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestUserRepository_CreateUser(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		user *domain.User
	}
	tests := []struct {
		name                       string
		args                       args
		mockSqlUserQueryExpected   string
		mockInsertedUserIdReturned int
		expectedUserUid            uint
	}{
		{
			name: "valid",
			args: args{
				user: &domain.User{
					Username: "mockUsername",
					Password: "mockPassword",
				},
			},
			mockSqlUserQueryExpected:   `INSERT INTO "users" ("created_at","updated_at","deleted_at","username","password") VALUES ($1,$2,$3,$4,$5)`,
			mockInsertedUserIdReturned: 666,
			expectedUserUid:            666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlUserQueryExpected)).
				WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					tt.args.user.Username, tt.args.user.Password,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.mockInsertedUserIdReturned))
			mockDb.ExpectCommit()

			actual, err := repo.CreateUser(context.Background(), tt.args.user)
			if err != nil {
				t.Errorf("CreateUser() error = %v", err)
			}

			assert.Equal(t, tt.expectedUserUid, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestUserRepository_CreateUserWithError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		user *domain.User
	}
	tests := []struct {
		name                     string
		args                     args
		mockSqlUserQueryExpected string
		expectedErrorMessage     string
	}{
		{
			name: "random error",
			args: args{
				user: &domain.User{
					Username: "mockUsername",
					Password: "mockPassword",
				},
			},
			mockSqlUserQueryExpected: `INSERT INTO "users" ("created_at","updated_at","deleted_at","username","password") VALUES ($1,$2,$3,$4,$5)`,
			expectedErrorMessage:     "random error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlUserQueryExpected)).
				WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					tt.args.user.Username, tt.args.user.Password,
				).
				WillReturnError(errors.New(tt.expectedErrorMessage))
			mockDb.ExpectRollback()

			_, err := repo.CreateUser(context.Background(), tt.args.user)
			actualErrorMessage := err.Error()

			assert.Equal(t, tt.expectedErrorMessage, actualErrorMessage)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
