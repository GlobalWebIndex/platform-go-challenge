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

func TestUserRepository_GetUser(t *testing.T) {
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
		name                                 string
		args                                 args
		mockSqlUserQueryExpected             string
		mockSqlUserChartJoinQueryExpected    string
		mockSqlChartQueryExpected            string
		mockSqlUserInsightJoinQueryExpected  string
		mockSqlInsightQueryExpected          string
		mockSqlUserAudienceJoinQueryExpected string
		mockSqlAudienceQueryExpected         string
		mockUserReturned                     *User
		expected                             *domain.User
	}{
		{
			name: "valid",
			args: args{
				uid: 666,
			},
			mockSqlUserQueryExpected:             `SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL LIMIT 1`,
			mockSqlUserChartJoinQueryExpected:    `SELECT * FROM "users_charts" WHERE "users_charts"."user_id" = $1`,
			mockSqlChartQueryExpected:            `SELECT * FROM "charts" WHERE "charts"."id" = $1 AND "charts"."deleted_at" IS NULL`,
			mockSqlUserInsightJoinQueryExpected:  `SELECT * FROM "users_insights" WHERE "users_insights"."user_id" = $1`,
			mockSqlInsightQueryExpected:          `SELECT * FROM "insights" WHERE "insights"."id" = $1 AND "insights"."deleted_at" IS NULL`,
			mockSqlUserAudienceJoinQueryExpected: `SELECT * FROM "users_audiences" WHERE "users_audiences"."user_id" = $1`,
			mockSqlAudienceQueryExpected:         `SELECT * FROM "audiences" WHERE "audiences"."id" = $1 AND "audiences"."deleted_at" IS NULL`,
			mockUserReturned: &User{
				Model:    gorm.Model{ID: 666},
				Username: "mockUsername",
				Password: "mockPassword",
				FavouriteCharts: []Chart{
					{
						Model: gorm.Model{ID: 1},
						Title: "title1",
					},
				},
				FavouriteInsights: []Insight{
					{
						Model: gorm.Model{ID: 1},
						Text:  "text1",
					},
				},
				FavouriteAudiences: []Audience{
					{
						Model:  gorm.Model{ID: 1},
						Gender: "gender1",
					},
				},
			},
			expected: &domain.User{
				Username: "mockUsername",
				Password: "mockPassword",
				FavouriteAssets: []domain.Asset{
					domain.Chart{
						Title:     "title1",
						Data:      "",
						AssetType: "chart",
					},
					domain.Audience{
						Gender:    "gender1",
						AssetType: "audience",
					},
					domain.Insight{
						Text:      "text1",
						AssetType: "insight",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.MatchExpectationsInOrder(false)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlUserQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "username", "password"},
					).AddRow(
						tt.mockUserReturned.ID, tt.mockUserReturned.Username, tt.mockUserReturned.Password,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlUserChartJoinQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"user_id", "chart_id"},
					).AddRow(
						tt.mockUserReturned.ID, tt.mockUserReturned.FavouriteCharts[0].ID,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlChartQueryExpected)).
				WithArgs(tt.mockUserReturned.FavouriteCharts[0].ID).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "title"},
					).AddRow(
						tt.mockUserReturned.FavouriteCharts[0].ID, tt.mockUserReturned.FavouriteCharts[0].Title,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlUserAudienceJoinQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"user_id", "audience_id"},
					).AddRow(
						tt.mockUserReturned.ID, tt.mockUserReturned.FavouriteAudiences[0].ID,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlAudienceQueryExpected)).
				WithArgs(tt.mockUserReturned.FavouriteAudiences[0].ID).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "gender"},
					).AddRow(
						tt.mockUserReturned.FavouriteAudiences[0].ID, tt.mockUserReturned.FavouriteAudiences[0].Gender,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlUserInsightJoinQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"user_id", "insight_id"},
					).AddRow(
						tt.mockUserReturned.ID, tt.mockUserReturned.FavouriteInsights[0].ID,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlInsightQueryExpected)).
				WithArgs(tt.mockUserReturned.FavouriteInsights[0].ID).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "text"},
					).AddRow(
						tt.mockUserReturned.FavouriteInsights[0].ID, tt.mockUserReturned.FavouriteInsights[0].Text,
					),
				)

			actual, err := repo.GetUser(context.Background(), tt.args.uid)
			if err != nil {
				t.Errorf("GetUser() error = %v", err)
				return
			}

			assert.Equal(t, tt.expected.Username, actual.Username)
			assert.Equal(t, tt.expected.Password, actual.Password)
			assert.ElementsMatch(t, tt.expected.FavouriteAssets, actual.FavouriteAssets)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestUserRepository_GetUserWithError(t *testing.T) {
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
		name                     string
		args                     args
		mockSqlUserQueryExpected string
		mockSqlErrorReturned     error
		expectedError            error
	}{
		{
			name: "random error",
			args: args{
				uid: 666,
			},
			mockSqlUserQueryExpected: `SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:     errors.New("random error"),
			expectedError:            errors.New("random error"),
		},
		{
			name: "description not found",
			args: args{
				uid: 666,
			},
			mockSqlUserQueryExpected: `SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL LIMIT 1`,
			mockSqlErrorReturned:     gorm.ErrRecordNotFound,
			expectedError: apierrors.UserNotFoundErrorWrapper{
				ReturnedStatusCode: http.StatusNotFound,
				OriginalError:      errors.New("userId 666 not found"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlUserQueryExpected)).
				WithArgs(tt.args.uid).
				WillReturnError(tt.mockSqlErrorReturned)

			_, actual := repo.GetUser(context.Background(), tt.args.uid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
