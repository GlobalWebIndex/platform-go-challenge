package assets

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"gwi-platform/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAudience_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	audience := &Audience{
		Gender:                  "Male",
		BirthCountry:            "USA",
		AgeGroup:                "25-34",
		DailyHoursOnSocialMedia: 2.5,
		PurchasesLastMonth:      3,
	}

	mock.ExpectPrepare("INSERT INTO taudiences")
	mock.ExpectExec("INSERT INTO taudiences").
		WithArgs(audience.Gender, audience.BirthCountry, audience.AgeGroup, audience.DailyHoursOnSocialMedia, audience.PurchasesLastMonth).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = audience.Add(ctx, db)
	if err != nil {
		t.Errorf("error was not expected while adding audience: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if audience.AssetID != 1 {
		t.Errorf("expected AssetID to be 1, got %d", audience.AssetID)
	}
}

func TestAudience_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	audience := &Audience{}

	mock.ExpectExec("DELETE FROM taudiences").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = audience.Delete(ctx, db, 1)
	if err != nil {
		t.Errorf("error was not expected while deleting audience: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAudience_Modify(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	audience := &Audience{
		BaseAsset:               models.BaseAsset{AssetID: 1, AssetType: models.AudienceAsset},
		Gender:                  "Female",
		BirthCountry:            "Canada",
		AgeGroup:                "35-44",
		DailyHoursOnSocialMedia: 3.0,
		PurchasesLastMonth:      5,
	}

	mock.ExpectExec("UPDATE taudiences").
		WithArgs(audience.Gender, audience.BirthCountry, audience.AgeGroup, audience.DailyHoursOnSocialMedia, audience.PurchasesLastMonth, audience.AssetID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = audience.Modify(ctx, db)
	if err != nil {
		t.Errorf("error was not expected while modifying audience: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAudience_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"audience_id", "gender", "birth_country", "age_group", "daily_hours_on_social_media", "purchases_last_month"}).
		AddRow(1, "Male", "USA", "25-34", 2.5, 3)

	mock.ExpectQuery("SELECT (.+) FROM taudiences").
		WithArgs(1).
		WillReturnRows(rows)

	audience := &Audience{}
	result, err := audience.Get(ctx, db, 1)
	if err != nil {
		t.Errorf("error was not expected while getting audience: %s", err)
	}

	expectedAudience := models.Audience{
		BaseAsset:               models.BaseAsset{AssetID: 1, AssetType: models.AudienceAsset},
		Gender:                  "Male",
		BirthCountry:            "USA",
		AgeGroup:                "25-34",
		DailyHoursOnSocialMedia: 2.5,
		PurchasesLastMonth:      3,
	}

	if !reflect.DeepEqual(result, expectedAudience) {
		t.Errorf("expected %+v, got %+v", expectedAudience, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAudience_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()

	mock.ExpectQuery("SELECT (.+) FROM taudiences").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	audience := &Audience{}
	_, err = audience.Get(ctx, db, 1)
	if err == nil {
		t.Error("expected error, got nil")
	} else if err.Error() != "audience with id 1 not found" {
		t.Errorf("expected 'not found' error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
