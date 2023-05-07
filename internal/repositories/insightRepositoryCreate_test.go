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

func TestInsightRepository_CreateInsight(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		insight *domain.Insight
	}
	tests := []struct {
		name                          string
		args                          args
		mockSqlInsightQueryExpected   string
		mockInsertedInsightIdReturned int
		expectedInsightUid            uint
	}{
		{
			name: "valid",
			args: args{
				insight: &domain.Insight{
					Text:        "mockText",
					Description: "mockDescription",
				},
			},
			mockSqlInsightQueryExpected:   `INSERT INTO "insights" ("created_at","updated_at","deleted_at","text","description") VALUES ($1,$2,$3,$4,$5)`,
			mockInsertedInsightIdReturned: 666,
			expectedInsightUid:            666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InsightRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlInsightQueryExpected)).
				WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					tt.args.insight.Text, tt.args.insight.Description,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.mockInsertedInsightIdReturned))
			mockDb.ExpectCommit()

			actual, err := repo.CreateInsight(context.Background(), tt.args.insight)
			if err != nil {
				t.Errorf("CreateInsight() error = %v", err)
			}

			assert.Equal(t, tt.expectedInsightUid, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestInsightRepository_CreateInsightWithError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		insight *domain.Insight
	}
	tests := []struct {
		name                        string
		args                        args
		mockSqlInsightQueryExpected string
		expectedErrorMessage        string
	}{
		{
			name: "random error",
			args: args{
				insight: &domain.Insight{
					Text:        "mockText",
					Description: "mockDescription",
				},
			},
			mockSqlInsightQueryExpected: `INSERT INTO "insights" ("created_at","updated_at","deleted_at","text","description") VALUES ($1,$2,$3,$4,$5)`,
			expectedErrorMessage:        "random error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InsightRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlInsightQueryExpected)).
				WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					tt.args.insight.Text, tt.args.insight.Description,
				).
				WillReturnError(errors.New(tt.expectedErrorMessage))
			mockDb.ExpectRollback()

			_, err := repo.CreateInsight(context.Background(), tt.args.insight)
			actualErrorMessage := err.Error()

			assert.Equal(t, tt.expectedErrorMessage, actualErrorMessage)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
