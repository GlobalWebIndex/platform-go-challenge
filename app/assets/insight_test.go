package assets

import (
	"context"
	"database/sql"
	"gwi-platform/models"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsight_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	insight := &Insight{
		BaseAsset: models.BaseAsset{AssetType: models.InsightAsset},
		Text:      "Test Insight",
	}

	mock.ExpectPrepare("INSERT INTO tinsights")
	mock.ExpectExec("INSERT INTO tinsights").
		WithArgs(insight.Text).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = insight.Add(ctx, db)
	if err != nil {
		t.Errorf("error was not expected while adding insight: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if insight.AssetID != 1 {
		t.Errorf("expected AssetID to be 1, got %d", insight.AssetID)
	}
}

func TestInsight_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	insight := &Insight{}

	mock.ExpectExec("DELETE FROM tinsights").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = insight.Delete(ctx, db, 1)
	if err != nil {
		t.Errorf("error was not expected while deleting insight: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsight_Modify(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	insight := &Insight{
		BaseAsset: models.BaseAsset{AssetID: 1, AssetType: models.InsightAsset},
		Text:      "Updated Insight",
	}

	mock.ExpectExec("UPDATE tinsights").
		WithArgs(insight.Text, insight.AssetID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = insight.Modify(ctx, db)
	if err != nil {
		t.Errorf("error was not expected while modifying insight: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsight_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"insight_id", "text"}).
		AddRow(1, "Test Insight")

	mock.ExpectQuery("SELECT (.+) FROM tinsights").
		WithArgs(1).
		WillReturnRows(rows)

	insight := &Insight{}
	result, err := insight.Get(ctx, db, 1)
	if err != nil {
		t.Errorf("error was not expected while getting insight: %s", err)
	}

	expectedInsight := models.Insight{
		BaseAsset: models.BaseAsset{AssetID: 1, AssetType: models.InsightAsset},
		Text:      "Test Insight",
	}

	if !reflect.DeepEqual(result, expectedInsight) {
		t.Errorf("expected %+v, got %+v", expectedInsight, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsight_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()

	mock.ExpectQuery("SELECT (.+) FROM tinsights").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	insight := &Insight{}
	_, err = insight.Get(ctx, db, 1)
	if err == nil {
		t.Error("expected error, got nil")
	} else if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
