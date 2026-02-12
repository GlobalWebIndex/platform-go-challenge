package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/giannis84/platform-go-challenge/internal/database"
	"github.com/giannis84/platform-go-challenge/internal/logging"
	"github.com/giannis84/platform-go-challenge/internal/models"
)

// testContext returns a context with a discarding logger for tests.
func testContext() context.Context {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return logging.NewContextWithLogger(context.Background(), logger)
}

var testCols = []string{"id", "user_id", "asset_type", "description", "data", "created_at", "updated_at"}

// setupTest creates a sqlmock-backed db and returns the mock + test context.
func setupTest(t *testing.T) (sqlmock.Sqlmock, context.Context) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	database.DB = db
	return mock, testContext()
}

func chartData(id string) []byte {
	data, _ := json.Marshal(&models.Chart{ID: id, Title: "T", XAxisTitle: "X", YAxisTitle: "Y"})
	return data
}

// assertError checks wantErr, optional *ValidationError, and error substring.
func assertError(t *testing.T, err error, wantErr, wantValErr bool, errSubstr string) {
	t.Helper()
	if wantErr && err == nil {
		t.Fatal("expected error, got nil")
	}
	if !wantErr && err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !wantErr {
		return
	}
	if wantValErr {
		var valErr *ValidationError
		if !errors.As(err, &valErr) {
			t.Fatalf("expected *ValidationError, got %T: %v", err, err)
		}
	}
	if errSubstr != "" && !strings.Contains(err.Error(), errSubstr) {
		t.Errorf("expected error to contain %q, got: %v", errSubstr, err)
	}
}

func assertValidation(t *testing.T, err error, wantErr bool, errSubstr string) {
	t.Helper()
	assertError(t, err, wantErr, false, errSubstr)
}

// --- Handler tests ---

func TestAddFavourite(t *testing.T) {
	insertOK := func(m sqlmock.Sqlmock) {
		m.ExpectExec("INSERT INTO favourites").WillReturnResult(sqlmock.NewResult(0, 1))
	}

	tests := []struct {
		name       string
		userID     string
		asset      models.Asset
		setupMock  func(sqlmock.Sqlmock)
		wantErr    bool
		wantValErr bool
		errSubstr  string
	}{
		{name: "valid chart", userID: "user1", asset: &models.Chart{ID: "c1", Title: "Revenue", XAxisTitle: "Month", YAxisTitle: "USD"}, setupMock: insertOK},
		{name: "valid insight", userID: "user1", asset: &models.Insight{ID: "i1", Text: "40% of millennials spend 3h on social media"}, setupMock: insertOK},
		{name: "valid audience", userID: "user1", asset: &models.Audience{ID: "a1", Gender: []string{"Male"}, BirthCountry: []string{"Greece"}, AgeGroups: []string{"25-34"}, SocialMediaHoursDaily: "3-5", PurchasesLastMonth: 5}, setupMock: insertOK},
		{name: "valid audience minimal", userID: "user1", asset: &models.Audience{ID: "a2"}, setupMock: insertOK},
		{name: "chart missing title", userID: "user1", asset: &models.Chart{ID: "c1", XAxisTitle: "X", YAxisTitle: "Y"}, wantErr: true, wantValErr: true, errSubstr: "title is required"},
		{name: "insight missing text", userID: "user1", asset: &models.Insight{ID: "i1"}, wantErr: true, wantValErr: true, errSubstr: "text is required"},
		{name: "audience invalid gender", userID: "user1", asset: &models.Audience{ID: "a1", Gender: []string{"Other"}}, wantErr: true, wantValErr: true, errSubstr: "gender[0] has invalid value"},
		{name: "audience negative purchases", userID: "user1", asset: &models.Audience{ID: "a1", PurchasesLastMonth: -1}, wantErr: true, wantValErr: true, errSubstr: "purchases_last_month must not be negative"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, ctx := setupTest(t)
			if tt.setupMock != nil {
				tt.setupMock(mock)
			}
			err := AddFavourite(ctx, tt.userID, tt.asset, "")
			assertError(t, err, tt.wantErr, tt.wantValErr, tt.errSubstr)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestUpdateDescription(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		userID      string
		assetID     string
		description string
		setupMock   func(sqlmock.Sqlmock)
		wantErr     bool
		wantValErr  bool
		errSubstr   string
	}{
		{
			name: "valid update", userID: "user1", assetID: "c1", description: "Updated description",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
					WithArgs("user1", "c1").
					WillReturnRows(sqlmock.NewRows(testCols).AddRow("c1", "user1", "chart", "old", chartData("c1"), now, now))
				m.ExpectExec("UPDATE favourites").WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{name: "empty description", userID: "user1", assetID: "c1", description: "", wantErr: true, wantValErr: true, errSubstr: "description is required"},
		{name: "whitespace-only", userID: "user1", assetID: "c1", description: "   ", wantErr: true, wantValErr: true, errSubstr: "description is required"},
		{name: "too long", userID: "user1", assetID: "c1", description: strings.Repeat("d", 256), wantErr: true, wantValErr: true, errSubstr: "description exceeds maximum length"},
		{
			name: "not found", userID: "user1", assetID: "nonexistent", description: "Some description",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
					WithArgs("user1", "nonexistent").
					WillReturnRows(sqlmock.NewRows(testCols))
			},
			wantErr: true, errSubstr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := setupTest(t)
			if tt.setupMock != nil {
				tt.setupMock(mock)
			}
			err := UpdateDescription(tt.userID, tt.assetID, tt.description)
			assertError(t, err, tt.wantErr, tt.wantValErr, tt.errSubstr)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestRemoveFavourite(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		assetID   string
		setupMock func(sqlmock.Sqlmock)
		wantErr   bool
		errSubstr string
	}{
		{
			name: "remove existing", userID: "user1", assetID: "c1",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM favourites").WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "not found", userID: "user1", assetID: "nonexistent",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM favourites").WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true, errSubstr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := setupTest(t)
			tt.setupMock(mock)
			err := RemoveFavourite(tt.userID, tt.assetID)
			assertError(t, err, tt.wantErr, false, tt.errSubstr)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestGetUserFavourites(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		userID    string
		setupMock func(sqlmock.Sqlmock)
		wantCount int
	}{
		{
			name: "returns favourites", userID: "user1", wantCount: 2,
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").WithArgs("user1").WillReturnRows(
					sqlmock.NewRows(testCols).
						AddRow("a", "user1", "chart", "", chartData("a"), now, now).
						AddRow("b", "user1", "chart", "", chartData("b"), now, now))
			},
		},
		{
			name: "empty for unknown user", userID: "unknown", wantCount: 0,
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").WithArgs("unknown").
					WillReturnRows(sqlmock.NewRows(testCols))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := setupTest(t)
			tt.setupMock(mock)
			favourites, err := GetUserFavourites(tt.userID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(favourites) != tt.wantCount {
				t.Errorf("expected %d favourites, got %d", tt.wantCount, len(favourites))
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

// --- Validation tests ---

func TestValidateChart(t *testing.T) {
	tests := []struct {
		name      string
		chart     models.Chart
		wantErr   bool
		errSubstr string
	}{
		{name: "valid chart", chart: models.Chart{ID: "c1", Title: "Revenue", XAxisTitle: "Month", YAxisTitle: "USD"}},
		{name: "missing id", chart: models.Chart{Title: "Revenue", XAxisTitle: "Month", YAxisTitle: "USD"}, wantErr: true, errSubstr: "id is required"},
		{name: "missing title", chart: models.Chart{ID: "c1", XAxisTitle: "Month", YAxisTitle: "USD"}, wantErr: true, errSubstr: "title is required"},
		{name: "missing x_axis_title", chart: models.Chart{ID: "c1", Title: "Revenue", YAxisTitle: "USD"}, wantErr: true, errSubstr: "x_axis_title is required"},
		{name: "missing y_axis_title", chart: models.Chart{ID: "c1", Title: "Revenue", XAxisTitle: "Month"}, wantErr: true, errSubstr: "y_axis_title is required"},
		{name: "title too long", chart: models.Chart{ID: "c1", Title: strings.Repeat("a", 256), XAxisTitle: "Month", YAxisTitle: "USD"}, wantErr: true, errSubstr: "title exceeds maximum length"},
		{name: "all required fields missing", chart: models.Chart{}, wantErr: true, errSubstr: "id is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidation(t, validateChart(&tt.chart), tt.wantErr, tt.errSubstr)
		})
	}
}

func TestValidateInsight(t *testing.T) {
	tests := []struct {
		name      string
		insight   models.Insight
		wantErr   bool
		errSubstr string
	}{
		{name: "valid insight", insight: models.Insight{ID: "i1", Text: "40% of millennials spend 3h on social media"}},
		{name: "missing id", insight: models.Insight{Text: "some text"}, wantErr: true, errSubstr: "id is required"},
		{name: "missing text", insight: models.Insight{ID: "i1"}, wantErr: true, errSubstr: "text is required"},
		{name: "text too long", insight: models.Insight{ID: "i1", Text: strings.Repeat("x", 256)}, wantErr: true, errSubstr: "text exceeds maximum length"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidation(t, validateInsight(&tt.insight), tt.wantErr, tt.errSubstr)
		})
	}
}

func TestValidateAudience(t *testing.T) {
	tests := []struct {
		name      string
		audience  models.Audience
		wantErr   bool
		errSubstr string
	}{
		{name: "valid audience with all fields", audience: models.Audience{ID: "a1", Gender: []string{"Male"}, BirthCountry: []string{"Greece"}, AgeGroups: []string{"25-34"}, SocialMediaHoursDaily: "3-5", PurchasesLastMonth: 5}},
		{name: "valid audience with only required fields", audience: models.Audience{ID: "a2"}},
		{name: "valid audience with partial optional fields", audience: models.Audience{ID: "a3", AgeGroups: []string{"18-24", "25-34"}}},
		{name: "missing id", audience: models.Audience{Gender: []string{"Male"}}, wantErr: true, errSubstr: "id is required"},
		{name: "invalid gender value", audience: models.Audience{ID: "a1", Gender: []string{"Other"}}, wantErr: true, errSubstr: "gender[0] has invalid value"},
		{name: "invalid age group value", audience: models.Audience{ID: "a1", AgeGroups: []string{"10-17"}}, wantErr: true, errSubstr: "age_groups[0] has invalid value"},
		{name: "invalid social media hours", audience: models.Audience{ID: "a1", SocialMediaHoursDaily: "10+"}, wantErr: true, errSubstr: "social_media_hours_daily has invalid value"},
		{name: "negative purchases", audience: models.Audience{ID: "a1", PurchasesLastMonth: -1}, wantErr: true, errSubstr: "purchases_last_month must not be negative"},
		{name: "empty birth country entry", audience: models.Audience{ID: "a1", BirthCountry: []string{""}}, wantErr: true, errSubstr: "birth_country[0] is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidation(t, validateAudience(&tt.audience), tt.wantErr, tt.errSubstr)
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
		errSubstr   string
	}{
		{name: "valid description", description: "My favourite chart"},
		{name: "empty description", description: "", wantErr: true, errSubstr: "description is required"},
		{name: "whitespace-only description", description: "   ", wantErr: true, errSubstr: "description is required"},
		{name: "description too long", description: strings.Repeat("d", 256), wantErr: true, errSubstr: "description exceeds maximum length"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidation(t, validateDescription(tt.description), tt.wantErr, tt.errSubstr)
		})
	}
}

func TestValidateAssetID(t *testing.T) {
	tests := []struct {
		name      string
		assetID   string
		wantErr   bool
		errSubstr string
	}{
		{name: "valid asset id", assetID: "asset123"},
		{name: "empty asset id", assetID: "", wantErr: true, errSubstr: "asset_id is required"},
		{name: "whitespace-only asset id", assetID: "   ", wantErr: true, errSubstr: "asset_id is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidation(t, ValidateAssetID(tt.assetID), tt.wantErr, tt.errSubstr)
		})
	}
}

func TestValidationErrorCollectsAllFieldErrors(t *testing.T) {
	err := validateChart(&models.Chart{})
	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(valErr.Errors) != 4 {
		t.Errorf("expected 4 validation errors, got %d: %v", len(valErr.Errors), valErr.Errors)
	}
}
