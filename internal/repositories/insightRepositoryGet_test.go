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

func TestInsightRepository_GetInsight(t *testing.T) {
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
		name                        string
		args                        args
		mockSqlInsightQueryExpected string
		mockInsightReturned         *Insight
		expected                    *domain.Insight
	}{
		{
			name: "valid",
			args: args{
				uid: 666,
			},
			mockSqlInsightQueryExpected: `SELECT * FROM "insights" WHERE id = $1 AND "insights"."deleted_at" IS NULL LIMIT 1`,
			mockInsightReturned: &Insight{
				Model:       gorm.Model{ID: 666},
				Text:        "mockText",
				Description: "mockDescription",
			},
			expected: &domain.Insight{
				Text:        "mockText",
				Description: "mockDescription",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InsightRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlInsightQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "text", "description"},
					).AddRow(
						tt.mockInsightReturned.ID, tt.mockInsightReturned.Text,
						tt.mockInsightReturned.Description,
					),
				)

			actual, err := repo.GetInsight(context.Background(), tt.args.uid)
			if err != nil {
				t.Errorf("GetInsight() error = %v", err)
				return
			}

			assert.Equal(t, tt.expected, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestInsightRepository_GetInsightWithError(t *testing.T) {
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
		name                        string
		args                        args
		mockSqlInsightQueryExpected string
		mockSqlErrorReturned        error
		expectedError               error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlInsightQueryExpected: `SELECT * FROM "insights" WHERE id = $1 AND "insights"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:        errors.New("random error"),
			expectedError:               errors.New("random error"),
		},
		{
			name: "description not found",
			args: args{
				uid: 666,
			},
			mockSqlInsightQueryExpected: `SELECT * FROM "insights" WHERE id = $1 AND "insights"."deleted_at" IS NULL LIMIT 1`,
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

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlInsightQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)

			_, actual := repo.GetInsight(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
