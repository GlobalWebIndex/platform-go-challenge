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

func TestAudienceRepository_GetAudience(t *testing.T) {
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
		name                         string
		args                         args
		mockSqlAudienceQueryExpected string
		mockAudienceReturned         *Audience
		expected                     *domain.Audience
	}{
		{
			name: "valid",
			args: args{
				uid: 666,
			},
			mockSqlAudienceQueryExpected: `SELECT * FROM "audiences" WHERE id = $1 AND "audiences"."deleted_at" IS NULL LIMIT 1`,
			mockAudienceReturned: &Audience{
				Model:              gorm.Model{ID: 666},
				Gender:             "mockGender",
				BirthCountry:       "mockBirthCountry",
				AgeGroup:           "mockAgeGroup",
				HoursSpentOnSocial: 1,
				PurchasesLastMonth: 2,
				Description:        "mockDescription",
			},
			expected: &domain.Audience{
				Gender:             "mockGender",
				BirthCountry:       "mockBirthCountry",
				AgeGroup:           "mockAgeGroup",
				HoursSpentDaily:    1,
				PurchasesLastMonth: 2,
				Description:        "mockDescription",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AudienceRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlAudienceQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "gender", "birth_country", "age_group", "hours_spent_on_social", "purchases_last_month", "description"},
					).AddRow(
						tt.mockAudienceReturned.ID, tt.mockAudienceReturned.Gender, tt.mockAudienceReturned.BirthCountry, tt.mockAudienceReturned.AgeGroup,
						tt.mockAudienceReturned.HoursSpentOnSocial, tt.mockAudienceReturned.PurchasesLastMonth, tt.mockAudienceReturned.Description,
					),
				)

			actual, err := repo.GetAudience(context.Background(), tt.args.uid)
			if err != nil {
				t.Errorf("GetAudience() error = %v", err)
				return
			}

			assert.Equal(t, tt.expected, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestAudienceRepository_GetAudienceWithError(t *testing.T) {
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
		name                         string
		args                         args
		mockSqlAudienceQueryExpected string
		mockSqlErrorReturned         error
		expectedError                error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlAudienceQueryExpected: `SELECT * FROM "audiences" WHERE id = $1 AND "audiences"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:         errors.New("random error"),
			expectedError:                errors.New("random error"),
		},
		{
			name: "description not found",
			args: args{
				uid: 666,
			},
			mockSqlAudienceQueryExpected: `SELECT * FROM "audiences" WHERE id = $1 AND "audiences"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:         gorm.ErrRecordNotFound,
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

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlAudienceQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)

			_, actual := repo.GetAudience(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
