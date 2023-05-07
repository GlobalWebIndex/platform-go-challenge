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

func TestInsightRepository_UpdateInsightDescription(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uid         uint32
		description string
	}
	tests := []struct {
		name                         string
		args                         args
		mockSqlInsightQueryExpected  string
		mockUpdatedInsightIdReturned int
	}{
		{
			name: "valid",
			args: args{
				uid:         666,
				description: "mockDescription",
			},
			mockSqlInsightQueryExpected:  `UPDATE "insights" SET "description"=$1,"updated_at"=$2 WHERE id = $3 AND "insights"."deleted_at" IS NULL`,
			mockUpdatedInsightIdReturned: 666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InsightRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlInsightQueryExpected)).
				WithArgs(tt.args.description, sqlmock.AnyArg(), tt.args.uid).
				WillReturnResult(sqlmock.NewResult(int64(tt.mockUpdatedInsightIdReturned), 1))
			mockDb.ExpectCommit()

			err := repo.UpdateInsightDescription(
				context.Background(),
				tt.args.uid,
				tt.args.description,
			)

			assert.NoError(t, err)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestInsightRepository_UpdateInsightDescriptionWithError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uid         uint32
		description string
	}
	tests := []struct {
		name                        string
		args                        args
		mockSqlInsightQueryExpected string
		mockSqlErrorReturned        error
		expectedError               error
	}{
		{
			name: "random error",
			args: args{
				uid:         666,
				description: "mockDescription",
			},
			mockSqlInsightQueryExpected: `UPDATE "insights" SET "description"=$1,"updated_at"=$2 WHERE id = $3 AND "insights"."deleted_at" IS NULL`,
			mockSqlErrorReturned:        errors.New("random error"),
			expectedError:               errors.New("random error"),
		},
		{
			name: "description not found",
			args: args{
				uid:         666,
				description: "mockDescription",
			},
			mockSqlInsightQueryExpected: `UPDATE "insights" SET "description"=$1,"updated_at"=$2 WHERE id = $3 AND "insights"."deleted_at" IS NULL`,
			mockSqlErrorReturned:        gorm.ErrRecordNotFound,
			expectedError: apierrors.UserNotFoundErrorWrapper{
				ReturnedStatusCode: http.StatusNotFound,
				OriginalError:      errors.New("userId 666 not found"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InsightRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlInsightQueryExpected)).
				WithArgs(tt.args.description, sqlmock.AnyArg(), tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)
			mockDb.ExpectRollback()

			actual := repo.UpdateInsightDescription(
				context.Background(),
				tt.args.uid,
				tt.args.description,
			)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
