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

func TestChartRepository_UpdateChartDescription(t *testing.T) {
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
		name                       string
		args                       args
		mockSqlChartQueryExpected  string
		mockUpdatedChartIdReturned int
	}{
		{
			name: "valid",
			args: args{
				uid:         666,
				description: "mockDescription",
			},
			mockSqlChartQueryExpected:  `UPDATE "charts" SET "description"=$1,"updated_at"=$2 WHERE id = $3 AND "charts"."deleted_at" IS NULL`,
			mockUpdatedChartIdReturned: 666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChartRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlChartQueryExpected)).
				WithArgs(tt.args.description, sqlmock.AnyArg(), tt.args.uid).
				WillReturnResult(sqlmock.NewResult(int64(tt.mockUpdatedChartIdReturned), 1))
			mockDb.ExpectCommit()

			err := repo.UpdateChartDescription(
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

func TestChartRepository_UpdateChartDescriptionWithError(t *testing.T) {
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
		name                      string
		args                      args
		mockSqlChartQueryExpected string
		mockSqlErrorReturned      error
		expectedError             error
	}{
		{
			name: "random error",
			args: args{
				uid:         666,
				description: "mockDescription",
			},
			mockSqlChartQueryExpected: `UPDATE "charts" SET "description"=$1,"updated_at"=$2 WHERE id = $3 AND "charts"."deleted_at" IS NULL`,
			mockSqlErrorReturned:      errors.New("random error"),
			expectedError:             errors.New("random error"),
		},
		{
			name: "description not found",
			args: args{
				uid:         666,
				description: "mockDescription",
			},
			mockSqlChartQueryExpected: `UPDATE "charts" SET "description"=$1,"updated_at"=$2 WHERE id = $3 AND "charts"."deleted_at" IS NULL`,
			mockSqlErrorReturned:      gorm.ErrRecordNotFound,
			expectedError: apierrors.UserNotFoundErrorWrapper{
				ReturnedStatusCode: http.StatusNotFound,
				OriginalError:      errors.New("userId 666 not found"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChartRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlChartQueryExpected)).
				WithArgs(tt.args.description, sqlmock.AnyArg(), tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)
			mockDb.ExpectRollback()

			actual := repo.UpdateChartDescription(
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
