package repositories

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestUserRepository_addFavouriteChart(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		userId  uint32
		assetId uint32
	}
	tests := []struct {
		name                              string
		args                              args
		mockSqlUserQueryExpected          string
		mockUserReturned                  *User
		mockSqlUserChartJoinQueryExpected string
		mockSqlJoinChartQueryExpected     string
		mockJoinedChartReturned           *Chart
		mockSqlNewChartQueryExpected      string
		mockNewChartReturned              *Chart
		mockSqlSaveUserQueryExpected      string
		mockSqlAddNewAssetQueryExpected   string
	}{
		{
			name: "valid",
			args: args{
				userId:  666,
				assetId: 777,
			},
			mockSqlUserQueryExpected:          `SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL LIMIT 1`,
			mockSqlUserChartJoinQueryExpected: `SELECT * FROM "users_charts" WHERE "users_charts"."user_id" = $1`,
			mockSqlJoinChartQueryExpected:     `SELECT * FROM "charts" WHERE "charts"."id" = $1 AND "charts"."deleted_at" IS NULL`,
			mockUserReturned: &User{
				Model:    gorm.Model{ID: 666},
				Username: "mockUsername",
				Password: "mockPassword",
				FavouriteCharts: []Chart{
					{
						Model: gorm.Model{ID: 667},
						Title: "title1",
					},
				},
			},
			mockJoinedChartReturned: &Chart{
				Model: gorm.Model{ID: 667},
				Title: "mockTitle",
			},
			mockSqlNewChartQueryExpected: `SELECT * FROM "charts" WHERE id = $1 AND "charts"."deleted_at" IS NULL LIMIT 1`,
			mockNewChartReturned: &Chart{
				Model: gorm.Model{ID: 777},
				Title: "mockTitle",
			},
			mockSqlSaveUserQueryExpected:    `UPDATE "users" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"username"=$4,"password"=$5 WHERE "users"."deleted_at" IS NULL AND "id" = $6`,
			mockSqlAddNewAssetQueryExpected: `INSERT INTO "users_charts" ("user_id","chart_id") VALUES ($1,$2),($3,$4) ON CONFLICT DO NOTHING`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlUserQueryExpected)).
				WithArgs(tt.args.userId).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "username", "password"},
					).AddRow(
						tt.mockUserReturned.ID, tt.mockUserReturned.Username, tt.mockUserReturned.Password,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlUserChartJoinQueryExpected)).
				WithArgs(tt.args.userId).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"user_id", "chart_id"},
					).AddRow(
						tt.mockUserReturned.ID, tt.mockJoinedChartReturned.ID,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlJoinChartQueryExpected)).
				WithArgs(tt.mockJoinedChartReturned.ID).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "title"},
					).AddRow(
						tt.mockJoinedChartReturned.ID, tt.mockJoinedChartReturned.Title,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlNewChartQueryExpected)).
				WithArgs(tt.args.assetId).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "title"},
					).AddRow(
						tt.mockNewChartReturned.ID, tt.mockNewChartReturned.Title,
					),
				)

			mockDb.ExpectBegin()

			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlSaveUserQueryExpected)).
				WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					tt.mockUserReturned.Username, tt.mockUserReturned.Password,
					tt.args.userId,
				).
				WillReturnResult(sqlmock.NewResult(int64(tt.mockUserReturned.ID), 1))

			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlAddNewAssetQueryExpected)).
				WithArgs(
					tt.mockUserReturned.ID, tt.mockUserReturned.FavouriteCharts[0].ID,
					uint(tt.args.userId), uint(tt.args.assetId),
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mockDb.ExpectCommit()

			err := repo.addFavouriteChart(context.Background(), tt.args.userId, tt.args.assetId)

			assert.NoError(t, err)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestUserRepository_editFavouriteChartDescription(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		userId      uint32
		assetId     uint32
		description string
	}
	tests := []struct {
		name                          string
		args                          args
		mockSqlJoinChartQueryExpected string
		mockJoinedChartReturned       *Chart
		mockSqlSaveChartQueryExpected string
	}{
		{
			name: "valid",
			args: args{
				userId:      666,
				assetId:     777,
				description: "new description",
			},
			mockJoinedChartReturned: &Chart{
				Model: gorm.Model{ID: 667},
				Title: "mockTitle",
			},
			mockSqlSaveChartQueryExpected: `UPDATE "charts" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"title"=$4,"x_title"=$5,"y_title"=$6,"data"=$7,"description"=$8 WHERE "charts"."deleted_at" IS NULL AND "id" = $9`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlJoinChartQueryExpected)).
				WithArgs(tt.args.userId, tt.args.assetId).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "title"},
					).AddRow(
						tt.mockJoinedChartReturned.ID, tt.mockJoinedChartReturned.Title,
					),
				)

			mockDb.ExpectBegin()

			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlSaveChartQueryExpected)).
				WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					tt.mockJoinedChartReturned.Title, tt.mockJoinedChartReturned.XTitle,
					tt.mockJoinedChartReturned.YTitle, tt.mockJoinedChartReturned.Data,
					tt.args.description, tt.mockJoinedChartReturned.ID,
				).
				WillReturnResult(sqlmock.NewResult(int64(tt.mockJoinedChartReturned.ID), 1))

			mockDb.ExpectCommit()

			err := repo.editFavouriteChartDescription(
				context.Background(),
				tt.args.userId,
				tt.args.assetId,
				tt.args.description,
			)

			assert.NoError(t, err)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestUserRepository_removeFavouriteChart(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		userId  uint32
		assetId uint32
	}
	tests := []struct {
		name                            string
		args                            args
		mockSelectUserQueryExpected     string
		mockSelectedUserReturned        *User
		mockSqlDeleteChartQueryExpected string
	}{
		{
			name: "valid",
			args: args{
				userId:  666,
				assetId: 777,
			},
			mockSelectUserQueryExpected: `SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`,
			mockSelectedUserReturned: &User{
				Model:    gorm.Model{ID: 666},
				Username: "username",
				Password: "password",
			},
			mockSqlDeleteChartQueryExpected: `DELETE FROM "users_charts" WHERE "users_charts"."user_id" = $1 AND "users_charts"."chart_id" = $2`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSelectUserQueryExpected)).
				WithArgs(tt.args.userId).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "username", "password"},
					).AddRow(
						tt.mockSelectedUserReturned.ID,
						tt.mockSelectedUserReturned.Username,
						tt.mockSelectedUserReturned.Password,
					),
				)

			mockDb.ExpectBegin()

			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlDeleteChartQueryExpected)).
				WithArgs(tt.args.userId, tt.args.assetId).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mockDb.ExpectCommit()

			err := repo.removeFavouriteChart(
				context.Background(),
				tt.args.userId,
				tt.args.assetId,
			)

			assert.NoError(t, err)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
