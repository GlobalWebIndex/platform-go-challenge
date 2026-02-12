package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/giannis84/platform-go-challenge/internal/auth"
	"github.com/giannis84/platform-go-challenge/internal/config"
	"github.com/giannis84/platform-go-challenge/internal/database"
	"github.com/giannis84/platform-go-challenge/internal/logging"
	"github.com/giannis84/platform-go-challenge/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

var testCols = []string{"id", "user_id", "asset_type", "description", "data", "created_at", "updated_at"}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func testToken(sub string) string {
	claims := jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	s, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	return s
}

func addAuthHeader(req *http.Request, userID string) {
	req.Header.Set("Authorization", "Bearer "+testToken(userID))
}

func setupTestHandler(t *testing.T) (*chi.Mux, sqlmock.Sqlmock) {
	t.Helper()
	logger := testLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	database.DB = db

	router := chi.NewRouter()
	router.Use(logging.RequestLogger(logger))
	router.Group(RegisterFavouritesRoutes(auth.AuthConfig{
		Secret:              "",
		AllowUnsignedTokens: true,
	}, config.RateLimitConfig{}))

	return router, mock
}

func insightRequestBody() map[string]any {
	return map[string]any{
		"asset_type":  "insight",
		"description": "Social media usage insight",
		"asset_data": map[string]any{
			"id":   "insight1",
			"text": "40% of millennials spend more than 3 hours on social media daily",
		},
	}
}

func audienceRequestBody() map[string]any {
	return map[string]any{
		"asset_type":  "audience",
		"description": "Tech-savvy millennials",
		"asset_data": map[string]any{
			"id":                      "audience1",
			"gender":                  []string{"Male", "Female"},
			"birth_country":           []string{"US", "UK"},
			"age_groups":              []string{"25-34"},
			"social_media_hours_daily": "3-5",
			"purchases_last_month":    5,
		},
	}
}

func postFavourite(t *testing.T, router *chi.Mux, body map[string]any) *httptest.ResponseRecorder {
	t.Helper()
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/favourites", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	addAuthHeader(req, "user1")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func TestFavouritesRoutes_AddFavourite(t *testing.T) {
	router, mock := setupTestHandler(t)

	mock.ExpectExec("INSERT INTO favourites").
		WillReturnResult(sqlmock.NewResult(0, 1))

	rr := postFavourite(t, router, insightRequestBody())

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("expected status %d, got %d. Body: %s", http.StatusCreated, status, rr.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["message"] != "Favourite added successfully" {
		t.Errorf("unexpected response: %v", resp)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestFavouritesRoutes_GetUserFavourites(t *testing.T) {
	router, mock := setupTestHandler(t)
	now := time.Now()

	// Add favourite
	mock.ExpectExec("INSERT INTO favourites").
		WillReturnResult(sqlmock.NewResult(0, 1))
	rr := postFavourite(t, router, insightRequestBody())
	if rr.Code != http.StatusCreated {
		t.Fatalf("setup failed: status %d, body: %s", rr.Code, rr.Body.String())
	}

	// Get favourites
	insightData, _ := json.Marshal(models.Insight{
		ID:   "insight1",
		Text: "40% of millennials spend more than 3 hours on social media daily",
	})
	mock.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
		WithArgs("user1").
		WillReturnRows(sqlmock.NewRows(testCols).
			AddRow("insight1", "user1", "insight", "Social media usage insight", insightData, now, now))

	req := httptest.NewRequest("GET", "/api/v1/favourites", nil)
	req.Header.Set("Accept", "application/json")
	addAuthHeader(req, "user1")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, status, rr.Body.String())
	}

	var favourites []any
	json.Unmarshal(rr.Body.Bytes(), &favourites)
	if len(favourites) != 1 {
		t.Errorf("expected 1 favourite, got %d", len(favourites))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestFavouritesRoutes_UpdateDescription(t *testing.T) {
	router, mock := setupTestHandler(t)
	now := time.Now()

	// Add favourite
	mock.ExpectExec("INSERT INTO favourites").
		WillReturnResult(sqlmock.NewResult(0, 1))
	rr := postFavourite(t, router, audienceRequestBody())
	if rr.Code != http.StatusCreated {
		t.Fatalf("setup failed: status %d, body: %s", rr.Code, rr.Body.String())
	}

	// Update description: handler calls GetFavouriteFromDB then UpdateFavouriteInDB
	audienceData, _ := json.Marshal(models.Audience{
		ID: "audience1", Gender: []string{"Male", "Female"}, BirthCountry: []string{"US", "UK"},
		AgeGroups: []string{"25-34"}, SocialMediaHoursDaily: "3-5", PurchasesLastMonth: 5,
	})
	mock.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
		WithArgs("user1", "audience1").
		WillReturnRows(sqlmock.NewRows(testCols).
			AddRow("audience1", "user1", "audience", "Tech-savvy millennials", audienceData, now, now))
	mock.ExpectExec("UPDATE favourites").
		WillReturnResult(sqlmock.NewResult(0, 1))

	updateBody, _ := json.Marshal(map[string]string{"description": "Updated description for audience"})
	req := httptest.NewRequest("PATCH", "/api/v1/favourites/audience1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	addAuthHeader(req, "user1")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, status, rr.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["message"] != "Description updated successfully" {
		t.Errorf("unexpected response: %v", resp)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestFavouritesRoutes_RemoveFavourite(t *testing.T) {
	router, mock := setupTestHandler(t)

	// Add favourite
	mock.ExpectExec("INSERT INTO favourites").
		WillReturnResult(sqlmock.NewResult(0, 1))
	rr := postFavourite(t, router, insightRequestBody())
	if rr.Code != http.StatusCreated {
		t.Fatalf("setup failed: status %d, body: %s", rr.Code, rr.Body.String())
	}

	// Remove favourite
	mock.ExpectExec("DELETE FROM favourites").
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/api/v1/favourites/insight1", nil)
	req.Header.Set("Accept", "application/json")
	addAuthHeader(req, "user1")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, status, rr.Body.String())
	}

	// Verify removed by getting empty list
	mock.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
		WithArgs("user1").
		WillReturnRows(sqlmock.NewRows(testCols))

	req = httptest.NewRequest("GET", "/api/v1/favourites", nil)
	req.Header.Set("Accept", "application/json")
	addAuthHeader(req, "user1")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var favourites []any
	json.Unmarshal(rr.Body.Bytes(), &favourites)
	if len(favourites) != 0 {
		t.Errorf("expected 0 favourites after removal, got %d", len(favourites))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestFavouritesRoutes_AddDuplicateFavourite(t *testing.T) {
	router, mock := setupTestHandler(t)

	// First add succeeds
	mock.ExpectExec("INSERT INTO favourites").
		WillReturnResult(sqlmock.NewResult(0, 1))
	rr := postFavourite(t, router, insightRequestBody())
	if rr.Code != http.StatusCreated {
		t.Fatalf("first add failed: status %d, body: %s", rr.Code, rr.Body.String())
	}

	// Second add returns unique violation
	mock.ExpectExec("INSERT INTO favourites").
		WillReturnError(&pq.Error{Code: "23505"})
	rr = postFavourite(t, router, insightRequestBody())

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("expected status %d, got %d. Body: %s", http.StatusConflict, status, rr.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != "Favourite already exists" {
		t.Errorf("unexpected error message: %v", resp)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestFavouritesRoutes_WhitespaceAssetID(t *testing.T) {
	router, _ := setupTestHandler(t)

	tests := []struct {
		name     string
		method   string
		assetID  string
		wantCode int
	}{
		{name: "PATCH with whitespace assetID", method: "PATCH", assetID: "%20%20", wantCode: http.StatusBadRequest},
		{name: "DELETE with whitespace assetID", method: "DELETE", assetID: "%20", wantCode: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.method == "PATCH" {
				body, _ := json.Marshal(map[string]string{"description": "test"})
				req = httptest.NewRequest(tt.method, "/api/v1/favourites/"+tt.assetID, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, "/api/v1/favourites/"+tt.assetID, nil)
			}
			req.Header.Set("Accept", "application/json")
			addAuthHeader(req, "user1")
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("expected status %d, got %d. Body: %s", tt.wantCode, rr.Code, rr.Body.String())
			}
		})
	}
}

func TestFavouritesRoutes_AcceptHeaderMiddleware(t *testing.T) {
	router, mock := setupTestHandler(t)

	tests := []struct {
		name       string
		accept     string
		wantCode   int
		wantError  string
	}{
		{name: "missing Accept header", accept: "", wantCode: http.StatusNotAcceptable, wantError: "Accept header must include application/json"},
		{name: "wrong Accept header", accept: "text/html", wantCode: http.StatusNotAcceptable, wantError: "Accept header must include application/json"},
		{name: "Accept */* is allowed", accept: "*/*", wantCode: http.StatusOK, wantError: ""},
		{name: "Accept application/json is allowed", accept: "application/json", wantCode: http.StatusOK, wantError: ""},
		{name: "Accept with multiple types including json", accept: "text/html, application/json", wantCode: http.StatusOK, wantError: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantCode == http.StatusOK {
				mock.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows(testCols))
			}

			req := httptest.NewRequest("GET", "/api/v1/favourites", nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			addAuthHeader(req, "user1")
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("expected status %d, got %d. Body: %s", tt.wantCode, rr.Code, rr.Body.String())
			}
			if tt.wantError != "" {
				var resp map[string]string
				json.Unmarshal(rr.Body.Bytes(), &resp)
				if resp["error"] != tt.wantError {
					t.Errorf("expected error %q, got %q", tt.wantError, resp["error"])
				}
			}
		})
	}
}

func TestFavouritesRoutes_ContentTypeMiddleware(t *testing.T) {
	router, mock := setupTestHandler(t)

	tests := []struct {
		name        string
		method      string
		path        string
		contentType string
		wantCode    int
		wantError   string
	}{
		{name: "POST missing Content-Type", method: "POST", path: "/api/v1/favourites", contentType: "", wantCode: http.StatusUnsupportedMediaType, wantError: "Content-Type header must be application/json"},
		{name: "POST wrong Content-Type", method: "POST", path: "/api/v1/favourites", contentType: "text/plain", wantCode: http.StatusUnsupportedMediaType, wantError: "Content-Type header must be application/json"},
		{name: "PATCH missing Content-Type", method: "PATCH", path: "/api/v1/favourites/test1", contentType: "", wantCode: http.StatusUnsupportedMediaType, wantError: "Content-Type header must be application/json"},
		{name: "PATCH wrong Content-Type", method: "PATCH", path: "/api/v1/favourites/test1", contentType: "application/xml", wantCode: http.StatusUnsupportedMediaType, wantError: "Content-Type header must be application/json"},
		{name: "GET without Content-Type is allowed", method: "GET", path: "/api/v1/favourites", contentType: "", wantCode: http.StatusOK, wantError: ""},
		{name: "DELETE without Content-Type is allowed", method: "DELETE", path: "/api/v1/favourites/test1", contentType: "", wantCode: http.StatusNotFound, wantError: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantCode == http.StatusOK {
				mock.ExpectQuery("SELECT .+ FROM favourites WHERE user_id").
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows(testCols))
			}
			if tt.wantCode == http.StatusNotFound && tt.method == "DELETE" {
				mock.ExpectExec("DELETE FROM favourites").
					WillReturnResult(sqlmock.NewResult(0, 0))
			}

			var req *http.Request
			if tt.method == "POST" || tt.method == "PATCH" {
				body, _ := json.Marshal(map[string]string{"description": "test"})
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(body))
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}
			req.Header.Set("Accept", "application/json")
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			addAuthHeader(req, "user1")
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("expected status %d, got %d. Body: %s", tt.wantCode, rr.Code, rr.Body.String())
			}
			if tt.wantError != "" {
				var resp map[string]string
				json.Unmarshal(rr.Body.Bytes(), &resp)
				if resp["error"] != tt.wantError {
					t.Errorf("expected error %q, got %q", tt.wantError, resp["error"])
				}
			}
		})
	}
}
