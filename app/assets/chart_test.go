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

func TestChart_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	chart := &Chart{
		BaseAsset:  models.BaseAsset{AssetType: models.ChartAsset},
		Title:      "Test Chart",
		XAxisTitle: "X Axis",
		YAxisTitle: "Y Axis",
		DataPoints: &[]models.DataPoint{{X: 1, Y: 2}, {X: 3, Y: 4}},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tcharts").
		WithArgs(chart.Title, chart.XAxisTitle, chart.YAxisTitle).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare("INSERT INTO tchart_data_points")
	mock.ExpectExec("INSERT INTO tchart_data_points").
		WithArgs(1, 1.0, 2.0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO tchart_data_points").
		WithArgs(1, 3.0, 4.0).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err = chart.Add(ctx, db)
	if err != nil {
		t.Errorf("error was not expected while adding chart: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if chart.AssetID != 1 {
		t.Errorf("expected AssetID to be 1, got %d", chart.AssetID)
	}
}

func TestChart_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	chart := &Chart{}

	mock.ExpectExec("DELETE FROM tcharts").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = chart.Delete(ctx, db, 1)
	if err != nil {
		t.Errorf("error was not expected while deleting chart: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChart_Modify(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	chart := &Chart{
		BaseAsset:  models.BaseAsset{AssetID: 1, AssetType: models.ChartAsset},
		Title:      "Updated Chart",
		XAxisTitle: "Updated X",
		YAxisTitle: "Updated Y",
		DataPoints: &[]models.DataPoint{{X: 5, Y: 6}},
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tcharts").
		WithArgs(chart.Title, chart.XAxisTitle, chart.YAxisTitle, chart.AssetID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM tchart_data_points").
		WithArgs(chart.AssetID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectPrepare("INSERT INTO tchart_data_points")
	mock.ExpectExec("INSERT INTO tchart_data_points").
		WithArgs(chart.AssetID, 5.0, 6.0).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err = chart.Modify(ctx, db)
	if err != nil {
		t.Errorf("error was not expected while modifying chart: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChart_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"chart_id", "title", "x_axis_title", "y_axis_title"}).
		AddRow(1, "Test Chart", "X Axis", "Y Axis")

	mock.ExpectQuery("SELECT (.+) FROM tcharts").
		WithArgs(1).
		WillReturnRows(rows)

	dataPointRows := sqlmock.NewRows([]string{"x_value", "y_value"}).
		AddRow(1, 2).
		AddRow(3, 4)

	mock.ExpectQuery("SELECT (.+) FROM tchart_data_points").
		WithArgs(1).
		WillReturnRows(dataPointRows)

	chart := &Chart{}
	result, err := chart.Get(ctx, db, 1)
	if err != nil {
		t.Errorf("error was not expected while getting chart: %s", err)
	}

	expectedChart := models.Chart{
		BaseAsset:  models.BaseAsset{AssetID: 1, AssetType: models.ChartAsset},
		Title:      "Test Chart",
		XAxisTitle: "X Axis",
		YAxisTitle: "Y Axis",
		DataPoints: &[]models.DataPoint{{X: 1, Y: 2}, {X: 3, Y: 4}},
	}

	if !reflect.DeepEqual(result, expectedChart) {
		t.Errorf("expected %+v, got %+v", expectedChart, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChart_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()

	mock.ExpectQuery("SELECT (.+) FROM tcharts").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	chart := &Chart{}
	_, err = chart.Get(ctx, db, 1)
	if err == nil {
		t.Error("expected error, got nil")
	} else if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
