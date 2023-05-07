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

// The delete action in the gorm is done by selecting the requested description and then
// soft delete it, so we have two possible error sources: the select and the update
func TestInsightRepository_DeleteInsightWithSelectError(t *testing.T) {
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
		name                              string
		args                              args
		mockSqlSelectInsightQueryExpected string
		mockSqlErrorReturned              error
		expectedError                     error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlSelectInsightQueryExpected: `SELECT * FROM "insights" WHERE id = $1 AND "insights"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:              errors.New("random error"),
			expectedError:                     errors.New("random error"),
		},
		{
			name: "description not found",
			args: args{
				uid: 666,
			},
			mockSqlSelectInsightQueryExpected: `SELECT * FROM "insights" WHERE id = $1 AND "insights"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:              gorm.ErrRecordNotFound,
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

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectInsightQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)

			actual := repo.DeleteInsight(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestInsightRepository_DeleteInsightWithDeleteError(t *testing.T) {
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
		name                              string
		args                              args
		mockSqlSelectInsightQueryExpected string
		mockSqlDeleteInsightQueryExpected string
		mockInsightReturned               *Insight
		mockSqlErrorReturned              error
		expectedError                     error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlSelectInsightQueryExpected: `SELECT * FROM "insights" WHERE id = $1 AND "insights"."deleted_at" IS NULL LIMIT 1`,
			mockInsightReturned: &Insight{
				Model:       gorm.Model{ID: 666},
				Text:        "mockText",
				Description: "mockDescription",
			},
			// The query is UPDATE as the deletion is SOFT in the gorm action I've set
			mockSqlDeleteInsightQueryExpected: `UPDATE "insights" SET "deleted_at"=$1 WHERE id = $2 AND "insights"."deleted_at" IS NULL`,
			mockSqlErrorReturned:              errors.New("random error"),
			expectedError:                     errors.New("random error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InsightRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectInsightQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "text", "description"},
					).AddRow(
						tt.mockInsightReturned.ID, tt.mockInsightReturned.Text,
						tt.mockInsightReturned.Description,
					),
				)
			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlDeleteInsightQueryExpected)).
				WithArgs(sqlmock.AnyArg(), tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)
			mockDb.ExpectRollback()

			actual := repo.DeleteInsight(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
