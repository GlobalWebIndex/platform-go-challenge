package repositories

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"testing"
)

func TestChartRepository_GetChart(t *testing.T) {
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
		name                      string
		args                      args
		mockSqlChartQueryExpected string
		mockChartReturned         *Chart
		expected                  *domain.Chart
	}{
		{
			name: "valid",
			args: args{
				uid: 666,
			},
			mockSqlChartQueryExpected: `SELECT * FROM "charts" WHERE id = $1 AND "charts"."deleted_at" IS NULL LIMIT 1`,
			mockChartReturned: &Chart{
				Model:       gorm.Model{ID: 666},
				Title:       "mockTitle",
				XTitle:      "mockXTitle",
				YTitle:      "mockYTitle",
				Data:        "mockData",
				Description: "mockDescription",
			},
			expected: &domain.Chart{
				Title:       "mockTitle",
				XAxisTitle:  "mockXTitle",
				YAxisTitle:  "mockYTitle",
				Data:        "mockData",
				Description: "mockDescription",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChartRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlChartQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "title", "x_title", "y_title", "data", "description"},
					).AddRow(
						tt.mockChartReturned.ID, tt.mockChartReturned.Title, tt.mockChartReturned.XTitle, tt.mockChartReturned.YTitle,
						tt.mockChartReturned.Data, tt.mockChartReturned.Description,
					),
				)

			actual, err := repo.GetChart(context.Background(), tt.args.uid)
			if err != nil {
				t.Errorf("GetChart() error = %v", err)
				return
			}

			assert.Equal(t, tt.expected, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestChartRepository_GetChartWithError(t *testing.T) {
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
		name                      string
		args                      args
		mockSqlChartQueryExpected string
		mockSqlErrorReturned      error
		expectedError             error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlChartQueryExpected: `SELECT * FROM "charts" WHERE id = $1 AND "charts"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:      errors.New("random error"),
			expectedError:             errors.New("random error"),
		},
		{
			name: "description not found",
			args: args{
				uid: 666,
			},
			mockSqlChartQueryExpected: `SELECT * FROM "charts" WHERE id = $1 AND "charts"."deleted_at" IS NULL LIMIT 1`,
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

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlChartQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)

			_, actual := repo.GetChart(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
