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

func TestChartRepository_DeleteChart(t *testing.T) {
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
		name                            string
		args                            args
		mockSqlSelectChartQueryExpected string
		mockSqlDeleteChartQueryExpected string
		mockChartReturned               *Chart
		mockDeletedChartIdReturned      int
	}{
		{
			name: "valid",
			args: args{
				uid: 666,
			},
			mockSqlSelectChartQueryExpected: `SELECT * FROM "charts" WHERE id = $1 AND "charts"."deleted_at" IS NULL LIMIT 1`,
			mockChartReturned: &Chart{
				Model:       gorm.Model{ID: 666},
				Title:       "mockTitle",
				XTitle:      "mockXTitle",
				YTitle:      "mockYTitle",
				Data:        "mockData",
				Description: "mockDescription",
			},
			// The query is UPDATE as the deletion is SOFT in the gorm action I've set
			mockSqlDeleteChartQueryExpected: `UPDATE "charts" SET "deleted_at"=$1 WHERE id = $2 AND "charts"."deleted_at" IS NULL`,
			mockDeletedChartIdReturned:      666,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChartRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectChartQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "title", "x_title", "y_title", "data", "description"},
					).AddRow(
						tt.mockChartReturned.ID, tt.mockChartReturned.Title, tt.mockChartReturned.XTitle, tt.mockChartReturned.YTitle,
						tt.mockChartReturned.Data, tt.mockChartReturned.Description,
					),
				)
			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlDeleteChartQueryExpected)).
				WithArgs(sqlmock.AnyArg(), tt.args.uid).
				WillReturnResult(sqlmock.NewResult(int64(tt.mockDeletedChartIdReturned), 1))
			mockDb.ExpectCommit()

			err := repo.DeleteChart(context.Background(), tt.args.uid)

			assert.NoError(t, err)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

// The delete action in the gorm is done by selecting the requested description and then
// soft delete it, so we have two possible error sources: the select and the update
func TestChartRepository_DeleteChartWithSelectError(t *testing.T) {
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
		name                            string
		args                            args
		mockSqlSelectChartQueryExpected string
		mockSqlErrorReturned            error
		expectedError                   error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlSelectChartQueryExpected: `SELECT * FROM "charts" WHERE id = $1 AND "charts"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:            errors.New("random error"),
			expectedError:                   errors.New("random error"),
		},
		{
			name: "description not found",
			args: args{
				uid: 666,
			},
			mockSqlSelectChartQueryExpected: `SELECT * FROM "charts" WHERE id = $1 AND "charts"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:            gorm.ErrRecordNotFound,
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

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectChartQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)

			actual := repo.DeleteChart(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestChartRepository_DeleteChartWithDeleteError(t *testing.T) {
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
		name                            string
		args                            args
		mockSqlSelectChartQueryExpected string
		mockSqlDeleteChartQueryExpected string
		mockChartReturned               *Chart
		mockSqlErrorReturned            error
		expectedError                   error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlSelectChartQueryExpected: `SELECT * FROM "charts" WHERE id = $1 AND "charts"."deleted_at" IS NULL LIMIT 1`,
			mockChartReturned: &Chart{
				Model:       gorm.Model{ID: 666},
				Title:       "mockTitle",
				XTitle:      "mockXTitle",
				YTitle:      "mockYTitle",
				Data:        "mockData",
				Description: "mockDescription",
			},
			// The query is UPDATE as the deletion is SOFT in the gorm action I've set
			mockSqlDeleteChartQueryExpected: `UPDATE "charts" SET "deleted_at"=$1 WHERE id = $2 AND "charts"."deleted_at" IS NULL`,
			mockSqlErrorReturned:            errors.New("random error"),
			expectedError:                   errors.New("random error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChartRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectChartQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "title", "x_title", "y_title", "data", "description"},
					).AddRow(
						tt.mockChartReturned.ID, tt.mockChartReturned.Title, tt.mockChartReturned.XTitle, tt.mockChartReturned.YTitle,
						tt.mockChartReturned.Data, tt.mockChartReturned.Description,
					),
				)
			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlDeleteChartQueryExpected)).
				WithArgs(sqlmock.AnyArg(), tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)
			mockDb.ExpectRollback()

			actual := repo.DeleteChart(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
