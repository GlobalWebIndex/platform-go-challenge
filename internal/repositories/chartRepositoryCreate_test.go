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

func TestChartRepository_CreateChart(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		chart *domain.Chart
	}
	tests := []struct {
		name                        string
		args                        args
		mockSqlChartQueryExpected   string
		mockInsertedChartIdReturned int
		expectedChartUid            uint
	}{
		{
			name: "valid",
			args: args{
				chart: &domain.Chart{
					Title:       "mockTitle",
					XAxisTitle:  "mockXTitle",
					YAxisTitle:  "mockYTitle",
					Data:        "mockData",
					Description: "mockDescription",
				},
			},
			mockSqlChartQueryExpected:   `INSERT INTO "charts" ("created_at","updated_at","deleted_at","title","x_title","y_title","data","description") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			mockInsertedChartIdReturned: 666,
			expectedChartUid:            666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChartRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlChartQueryExpected)).
				WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					tt.args.chart.Title, tt.args.chart.XAxisTitle, tt.args.chart.YAxisTitle,
					tt.args.chart.Data, tt.args.chart.Description,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.mockInsertedChartIdReturned))
			mockDb.ExpectCommit()

			actual, err := repo.CreateChart(context.Background(), tt.args.chart)
			if err != nil {
				t.Errorf("CreateChart() error = %v", err)
			}

			assert.Equal(t, tt.expectedChartUid, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestChartRepository_CreateChartWithError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		chart *domain.Chart
	}
	tests := []struct {
		name                      string
		args                      args
		mockSqlChartQueryExpected string
		expectedErrorMessage      string
	}{
		{
			name: "random error",
			args: args{
				chart: &domain.Chart{
					Title:       "mockTitle",
					XAxisTitle:  "mockXTitle",
					YAxisTitle:  "mockYTitle",
					Data:        "mockData",
					Description: "mockDescription",
				},
			},
			mockSqlChartQueryExpected: `INSERT INTO "charts" ("created_at","updated_at","deleted_at","title","x_title","y_title","data","description") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			expectedErrorMessage:      "random error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChartRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlChartQueryExpected)).
				WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					tt.args.chart.Title, tt.args.chart.XAxisTitle, tt.args.chart.YAxisTitle,
					tt.args.chart.Data, tt.args.chart.Description,
				).
				WillReturnError(errors.New(tt.expectedErrorMessage))
			mockDb.ExpectRollback()

			_, err := repo.CreateChart(context.Background(), tt.args.chart)
			actualErrorMessage := err.Error()

			assert.Equal(t, tt.expectedErrorMessage, actualErrorMessage)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
