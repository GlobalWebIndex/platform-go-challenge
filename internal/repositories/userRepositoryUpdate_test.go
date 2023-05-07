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

func TestUserRepository_UpdateUser(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uid  uint32
		user *domain.User
	}
	tests := []struct {
		name                      string
		args                      args
		mockSqlUserQueryExpected  string
		mockUpdatedUserIdReturned int
	}{
		{
			name: "valid",
			args: args{
				uid: 666,
				user: &domain.User{
					Username: "mockUsername",
					Password: "mockPassword",
				},
			},
			mockSqlUserQueryExpected:  `UPDATE "users" SET "password"=$1,"username"=$2,"updated_at"=$3 WHERE id = $4 AND "users"."deleted_at" IS NULL`,
			mockUpdatedUserIdReturned: 666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlUserQueryExpected)).
				WithArgs(
					tt.args.user.Password, tt.args.user.Username,
					sqlmock.AnyArg(), tt.args.uid,
				).
				WillReturnResult(sqlmock.NewResult(int64(tt.mockUpdatedUserIdReturned), 1))
			mockDb.ExpectCommit()

			err := repo.UpdateUser(
				context.Background(),
				tt.args.uid,
				tt.args.user,
			)

			assert.NoError(t, err)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestUserRepository_UpdateUserWithError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uid  uint32
		user *domain.User
	}
	tests := []struct {
		name                     string
		args                     args
		mockSqlUserQueryExpected string
		mockSqlErrorReturned     error
		expectedError            error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
				user: &domain.User{
					Username: "mockUsername",
					Password: "mockPassword",
				},
			},
			mockSqlUserQueryExpected: `UPDATE "users" SET "password"=$1,"username"=$2,"updated_at"=$3 WHERE id = $4 AND "users"."deleted_at" IS NULL`,
			mockSqlErrorReturned:     errors.New("random error"),
			expectedError:            errors.New("random error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlUserQueryExpected)).
				WithArgs(
					tt.args.user.Password, tt.args.user.Username,
					sqlmock.AnyArg(), tt.args.uid,
				).
				WillReturnError(tt.mockSqlErrorReturned)
			mockDb.ExpectRollback()

			actual := repo.UpdateUser(
				context.Background(),
				tt.args.uid,
				tt.args.user,
			)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
