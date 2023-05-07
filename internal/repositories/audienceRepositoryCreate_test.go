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

func TestAudienceRepository_CreateAudience(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		audience *domain.Audience
	}
	tests := []struct {
		name                           string
		args                           args
		mockSqlAudienceQueryExpected   string
		mockInsertedAudienceIdReturned int
		expectedAudienceUid            uint
	}{
		{
			name: "valid",
			args: args{
				audience: &domain.Audience{
					Gender:             "mockGender",
					BirthCountry:       "mockBirthCountry",
					AgeGroup:           "mockAgeGroup",
					HoursSpentDaily:    1,
					PurchasesLastMonth: 2,
					Description:        "mockDescription",
				},
			},
			mockSqlAudienceQueryExpected:   `INSERT INTO "audiences" ("created_at","updated_at","deleted_at","gender","birth_country","age_group","hours_spent_on_social","purchases_last_month","description") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			mockInsertedAudienceIdReturned: 666,
			expectedAudienceUid:            666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AudienceRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlAudienceQueryExpected)).
				WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					tt.args.audience.Gender, tt.args.audience.BirthCountry, tt.args.audience.AgeGroup,
					tt.args.audience.HoursSpentDaily, tt.args.audience.PurchasesLastMonth,
					tt.args.audience.Description,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.mockInsertedAudienceIdReturned))
			mockDb.ExpectCommit()

			actual, err := repo.CreateAudience(context.Background(), tt.args.audience)
			if err != nil {
				t.Errorf("CreateAudience() error = %v", err)
			}

			assert.Equal(t, tt.expectedAudienceUid, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestAudienceRepository_CreateAudienceWithError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		audience *domain.Audience
	}
	tests := []struct {
		name                         string
		args                         args
		mockSqlAudienceQueryExpected string
		expectedErrorMessage         string
	}{
		{
			name: "random error",
			args: args{
				audience: &domain.Audience{
					Gender:             "mockGender",
					BirthCountry:       "mockBirthCountry",
					AgeGroup:           "mockAgeGroup",
					HoursSpentDaily:    1,
					PurchasesLastMonth: 2,
					Description:        "mockDescription",
				},
			},
			mockSqlAudienceQueryExpected: `INSERT INTO "audiences" ("created_at","updated_at","deleted_at","gender","birth_country","age_group","hours_spent_on_social","purchases_last_month","description") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			expectedErrorMessage:         "random error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AudienceRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlAudienceQueryExpected)).
				WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					tt.args.audience.Gender, tt.args.audience.BirthCountry, tt.args.audience.AgeGroup,
					tt.args.audience.HoursSpentDaily, tt.args.audience.PurchasesLastMonth,
					tt.args.audience.Description,
				).
				WillReturnError(errors.New(tt.expectedErrorMessage))
			mockDb.ExpectRollback()

			_, err := repo.CreateAudience(context.Background(), tt.args.audience)
			actualErrorMessage := err.Error()

			assert.Equal(t, tt.expectedErrorMessage, actualErrorMessage)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
