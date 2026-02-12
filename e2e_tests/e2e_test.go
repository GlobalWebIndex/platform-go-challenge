// Package e2e_tests provides end-to-end tests for the favourites API
// running against real Docker Compose services (app + PostgreSQL).
//
// Usage:
//
//	./run.sh              # starts compose, runs tests, tears down
//	go test -v -count=1   # if services are already running
//
// Override the default base URL with:
//
//	API_BASE_URL=http://localhost:9000 go test -v -count=1
package e2e_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultAPIBase    = "http://localhost:8000"
	defaultHealthBase = "http://localhost:8001"
)

func apiBase() string {
	if v := os.Getenv("API_BASE_URL"); v != "" {
		return v
	}
	return defaultAPIBase
}

func healthBase() string {
	if v := os.Getenv("HEALTH_BASE_URL"); v != "" {
		return v
	}
	return defaultHealthBase
}

// favouritesURL builds the full URL for the favourites endpoint.
func favouritesURL() string {
	return fmt.Sprintf("%s/api/v1/favourites", apiBase())
}

// favouriteURL builds the full URL for a specific favourite asset.
func favouriteURL(assetID string) string {
	return fmt.Sprintf("%s/api/v1/favourites/%s", apiBase(), assetID)
}

// jwtSecret returns the JWT secret used by the running service.
// When empty, the service accepts unsigned tokens (alg=none).
func jwtSecret() string {
	return os.Getenv("JWT_SECRET")
}

// tokenForUser returns a Bearer token string for the given user.
func tokenForUser(userID string) string {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	secret := jwtSecret()
	if secret == "" {
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		s, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		return s
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(secret))
	return s
}

// ---------- TestMain: wait for services before running ----------

func TestMain(m *testing.M) {
	if err := waitForReady(60 * time.Second); err != nil {
		fmt.Fprintf(os.Stderr, "services not ready: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func waitForReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	url := healthBase() + "/health/ready"

	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timed out waiting for %s", url)
}

// ---------- Helpers ----------

type apiResponse struct {
	StatusCode int
	Body       []byte
}

func doRequest(t *testing.T, method, url, userID string, body any) apiResponse {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if userID != "" {
		req.Header.Set("Authorization", "Bearer "+tokenForUser(userID))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s failed: %v", method, url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	return apiResponse{StatusCode: resp.StatusCode, Body: respBody}
}

func requireStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("expected status %d, got %d", want, got)
	}
}

// cleanup removes a favourite, ignoring 404s (already deleted).
func cleanup(t *testing.T, userID, assetID string) {
	t.Helper()
	doRequest(t, http.MethodDelete, favouriteURL(assetID), userID, nil)
}

// ---------- Test data ----------

func chartPayload(id string) map[string]any {
	return map[string]any{
		"asset_type":  "chart",
		"description": "Monthly sales data",
		"asset_data": map[string]any{
			"id":           id,
			"title":        "Sales Chart",
			"x_axis_title": "Month",
			"y_axis_title": "Revenue",
			"data":         map[string]any{"Jan": 100, "Feb": 200},
		},
	}
}

func insightPayload(id string) map[string]any {
	return map[string]any{
		"asset_type":  "insight",
		"description": "Social media usage insight",
		"asset_data": map[string]any{
			"id":   id,
			"text": "40% of millennials spend more than 3 hours on social media daily",
		},
	}
}

func audiencePayload(id string) map[string]any {
	return map[string]any{
		"asset_type":  "audience",
		"description": "Tech-savvy millennials",
		"asset_data": map[string]any{
			"id":                       id,
			"gender":                   []string{"Male"},
			"birth_country":            []string{"US", "UK"},
			"age_groups":               []string{"25-34"},
			"social_media_hours_daily": "3-5",
			"purchases_last_month":     5,
		},
	}
}

// ---------- Tests ----------

func TestAddFavourite(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		assetID    string
		payload    map[string]any
		wantStatus int
	}{
		{
			name:       "chart",
			userID:     "e2e-add-1",
			assetID:    "e2e-chart-1",
			payload:    chartPayload("e2e-chart-1"),
			wantStatus: http.StatusCreated,
		},
		{
			name:       "insight",
			userID:     "e2e-add-1",
			assetID:    "e2e-insight-1",
			payload:    insightPayload("e2e-insight-1"),
			wantStatus: http.StatusCreated,
		},
		{
			name:       "audience",
			userID:     "e2e-add-1",
			assetID:    "e2e-audience-1",
			payload:    audiencePayload("e2e-audience-1"),
			wantStatus: http.StatusCreated,
		},
		{
			name:    "invalid asset type",
			userID:  "e2e-add-2",
			assetID: "",
			payload: map[string]any{
				"asset_type": "unknown",
				"asset_data": map[string]any{"id": "x"},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.assetID != "" {
				t.Cleanup(func() { cleanup(t, tt.userID, tt.assetID) })
			}
			resp := doRequest(t, http.MethodPost, favouritesURL(), tt.userID, tt.payload)
			requireStatus(t, resp.StatusCode, tt.wantStatus)
		})
	}
}

func TestAddFavourite_Duplicate(t *testing.T) {
	const userID, assetID = "e2e-dup-1", "e2e-dup-chart"
	t.Cleanup(func() { cleanup(t, userID, assetID) })

	resp := doRequest(t, http.MethodPost, favouritesURL(), userID, chartPayload(assetID))
	requireStatus(t, resp.StatusCode, http.StatusCreated)

	resp = doRequest(t, http.MethodPost, favouritesURL(), userID, chartPayload(assetID))
	requireStatus(t, resp.StatusCode, http.StatusConflict)
}

func TestAddFavourite_InvalidBody(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, favouritesURL(),
		bytes.NewReader([]byte("not json")))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenForUser("e2e-add-3"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	requireStatus(t, resp.StatusCode, http.StatusBadRequest)
}

func TestGetFavourites(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		seed      []map[string]any // payloads to add before GET
		seedIDs   []string         // asset IDs for cleanup
		wantCount int
	}{
		{
			name:      "returns seeded favourites",
			userID:    "e2e-get-1",
			seed:      []map[string]any{chartPayload("e2e-get-chart"), insightPayload("e2e-get-insight")},
			seedIDs:   []string{"e2e-get-chart", "e2e-get-insight"},
			wantCount: 2,
		},
		{
			name:      "empty for non-existent user",
			userID:    "e2e-get-nonexistent",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, id := range tt.seedIDs {
				t.Cleanup(func() { cleanup(t, tt.userID, id) })
			}
			for _, payload := range tt.seed {
				doRequest(t, http.MethodPost, favouritesURL(), tt.userID, payload)
			}

			resp := doRequest(t, http.MethodGet, favouritesURL(), tt.userID, nil)
			requireStatus(t, resp.StatusCode, http.StatusOK)

			var favourites []map[string]any
			if err := json.Unmarshal(resp.Body, &favourites); err != nil {
				t.Fatalf("unmarshal response: %v", err)
			}
			if len(favourites) != tt.wantCount {
				t.Fatalf("expected %d favourites, got %d", tt.wantCount, len(favourites))
			}

			// Verify response structure for non-empty results
			for _, fav := range favourites {
				for _, field := range []string{"id", "user_id", "asset_type", "description", "data", "created_at", "updated_at"} {
					if _, ok := fav[field]; !ok {
						t.Errorf("favourite missing field %q", field)
					}
				}
				if fav["user_id"] != tt.userID {
					t.Errorf("expected user_id=%q, got %q", tt.userID, fav["user_id"])
				}
			}
		})
	}
}

func TestUpdateDescription(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		assetID     string
		seed        bool // whether to seed a favourite first
		description string
		wantStatus  int
	}{
		{
			name:        "valid update",
			userID:      "e2e-upd-1",
			assetID:     "e2e-upd-chart",
			seed:        true,
			description: "Updated via e2e test",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "empty description returns 400",
			userID:      "e2e-upd-2",
			assetID:     "e2e-upd-val",
			seed:        true,
			description: "",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "non-existent favourite returns 404",
			userID:      "e2e-upd-3",
			assetID:     "nonexistent",
			seed:        false,
			description: "does not matter",
			wantStatus:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.seed {
				t.Cleanup(func() { cleanup(t, tt.userID, tt.assetID) })
				doRequest(t, http.MethodPost, favouritesURL(), tt.userID, chartPayload(tt.assetID))
			}

			resp := doRequest(t, http.MethodPatch, favouriteURL(tt.assetID), tt.userID,
				map[string]string{"description": tt.description})
			requireStatus(t, resp.StatusCode, tt.wantStatus)

			// Verify the update persisted on success
			if tt.wantStatus == http.StatusOK {
				getResp := doRequest(t, http.MethodGet, favouritesURL(), tt.userID, nil)
				var favourites []map[string]any
				json.Unmarshal(getResp.Body, &favourites)
				if len(favourites) != 1 {
					t.Fatalf("expected 1 favourite, got %d", len(favourites))
				}
				if favourites[0]["description"] != tt.description {
					t.Errorf("expected description %q, got %q", tt.description, favourites[0]["description"])
				}
			}
		})
	}
}

func TestDeleteFavourite(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		assetID    string
		seed       bool
		wantStatus int
	}{
		{
			name:       "existing favourite",
			userID:     "e2e-del-1",
			assetID:    "e2e-del-chart",
			seed:       true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existent favourite returns 404",
			userID:     "e2e-del-2",
			assetID:    "nonexistent",
			seed:       false,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.seed {
				doRequest(t, http.MethodPost, favouritesURL(), tt.userID, chartPayload(tt.assetID))
			}

			resp := doRequest(t, http.MethodDelete, favouriteURL(tt.assetID), tt.userID, nil)
			requireStatus(t, resp.StatusCode, tt.wantStatus)

			// Verify deletion on success
			if tt.wantStatus == http.StatusOK {
				getResp := doRequest(t, http.MethodGet, favouritesURL(), tt.userID, nil)
				var favourites []any
				json.Unmarshal(getResp.Body, &favourites)
				if len(favourites) != 0 {
					t.Fatalf("expected 0 favourites after deletion, got %d", len(favourites))
				}
			}
		})
	}
}

func TestHealth(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		wantStatus int
	}{
		{
			name:       "liveness",
			endpoint:   "/health/live",
			wantStatus: http.StatusOK,
		},
		{
			name:       "readiness",
			endpoint:   "/health/ready",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(healthBase() + tt.endpoint)
			if err != nil {
				t.Fatalf("%s request failed: %v", tt.endpoint, err)
			}
			defer resp.Body.Close()
			requireStatus(t, resp.StatusCode, tt.wantStatus)
		})
	}
}
