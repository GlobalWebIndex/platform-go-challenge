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

func TestAudienceRepository_DeleteAudience(t *testing.T) {
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
		name                               string
		args                               args
		mockSqlSelectAudienceQueryExpected string
		mockSqlDeleteAudienceQueryExpected string
		mockAudienceReturned               *Audience
		mockDeletedAudienceIdReturned      int
	}{
		{
			name: "valid",
			args: args{
				uid: 666,
			},
			mockSqlSelectAudienceQueryExpected: `SELECT * FROM "audiences" WHERE id = $1 AND "audiences"."deleted_at" IS NULL LIMIT 1`,
			mockAudienceReturned: &Audience{
				Model:              gorm.Model{ID: 666},
				Gender:             "mockGender",
				BirthCountry:       "mockBirthCountry",
				AgeGroup:           "mockAgeGroup",
				HoursSpentOnSocial: 1,
				PurchasesLastMonth: 2,
				Description:        "mockDescription",
			},
			// The query is UPDATE as the deletion is SOFT in the gorm action I've set
			mockSqlDeleteAudienceQueryExpected: `UPDATE "audiences" SET "deleted_at"=$1 WHERE id = $2 AND "audiences"."deleted_at" IS NULL`,
			mockDeletedAudienceIdReturned:      666,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AudienceRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectAudienceQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "gender", "birth_country", "age_group", "hours_spent_on_social", "purchases_last_month", "description"},
					).AddRow(
						tt.mockAudienceReturned.ID, tt.mockAudienceReturned.Gender, tt.mockAudienceReturned.BirthCountry, tt.mockAudienceReturned.AgeGroup,
						tt.mockAudienceReturned.HoursSpentOnSocial, tt.mockAudienceReturned.PurchasesLastMonth, tt.mockAudienceReturned.Description,
					),
				)
			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlDeleteAudienceQueryExpected)).
				WithArgs(sqlmock.AnyArg(), tt.args.uid).
				WillReturnResult(sqlmock.NewResult(int64(tt.mockDeletedAudienceIdReturned), 1))
			mockDb.ExpectCommit()

			err := repo.DeleteAudience(context.Background(), tt.args.uid)

			assert.NoError(t, err)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

// The delete action in the gorm is done by selecting the requested description and then
// soft delete it, so we have two possible error sources: the select and the update
func TestAudienceRepository_DeleteAudienceWithSelectError(t *testing.T) {
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
		name                               string
		args                               args
		mockSqlSelectAudienceQueryExpected string
		mockSqlErrorReturned               error
		expectedError                      error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlSelectAudienceQueryExpected: `SELECT * FROM "audiences" WHERE id = $1 AND "audiences"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:               errors.New("random error"),
			expectedError:                      errors.New("random error"),
		},
		{
			name: "description not found",
			args: args{
				uid: 666,
			},
			mockSqlSelectAudienceQueryExpected: `SELECT * FROM "audiences" WHERE id = $1 AND "audiences"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:               gorm.ErrRecordNotFound,
			expectedError: apierrors.UserNotFoundErrorWrapper{
				ReturnedStatusCode: http.StatusNotFound,
				OriginalError:      errors.New("userId 666 not found"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AudienceRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectAudienceQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)

			actual := repo.DeleteAudience(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestAudienceRepository_DeleteAudienceWithDeleteError(t *testing.T) {
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
		name                               string
		args                               args
		mockSqlSelectAudienceQueryExpected string
		mockSqlDeleteAudienceQueryExpected string
		mockAudienceReturned               *Audience
		mockSqlErrorReturned               error
		expectedError                      error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlSelectAudienceQueryExpected: `SELECT * FROM "audiences" WHERE id = $1 AND "audiences"."deleted_at" IS NULL LIMIT 1`,
			mockAudienceReturned: &Audience{
				Model:              gorm.Model{ID: 666},
				Gender:             "mockGender",
				BirthCountry:       "mockBirthCountry",
				AgeGroup:           "mockAgeGroup",
				HoursSpentOnSocial: 1,
				PurchasesLastMonth: 2,
				Description:        "mockDescription",
			},
			// The query is UPDATE as the deletion is SOFT in the gorm action I've set
			mockSqlDeleteAudienceQueryExpected: `UPDATE "audiences" SET "deleted_at"=$1 WHERE id = $2 AND "audiences"."deleted_at" IS NULL`,
			mockSqlErrorReturned:               errors.New("random error"),
			expectedError:                      errors.New("random error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AudienceRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlSelectAudienceQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "gender", "birth_country", "age_group", "hours_spent_on_social", "purchases_last_month", "description"},
					).AddRow(
						tt.mockAudienceReturned.ID, tt.mockAudienceReturned.Gender, tt.mockAudienceReturned.BirthCountry, tt.mockAudienceReturned.AgeGroup,
						tt.mockAudienceReturned.HoursSpentOnSocial, tt.mockAudienceReturned.PurchasesLastMonth, tt.mockAudienceReturned.Description,
					),
				)
			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlDeleteAudienceQueryExpected)).
				WithArgs(sqlmock.AnyArg(), tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)
			mockDb.ExpectRollback()

			actual := repo.DeleteAudience(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
