package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/app"
)

// Happy-path integration: add a favourite and list it back.
func TestFavouritesIntegration_AddAndList(t *testing.T) {
	t.Parallel()

	router := app.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()
	baseURL := server.URL

	// --- 1. Add a favourite ---

	createBody := map[string]any{
		"asset": map[string]any{
			"type":        "insight",
			"description": "An example insight",
			"insight": map[string]any{
				"text": "40% of millennials spend more than 3 hours on social media daily",
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(createBody); err != nil {
		t.Fatalf("encode request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/v1/users/user-1/favourites", &buf)
	if err != nil {
		t.Fatalf("create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "user-1")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST favourites request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var createResp struct {
		UserID    string `json:"userId"`
		AssetID   string `json:"assetId"`
		CreatedAt string `json:"createdAt"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if createResp.UserID != "user-1" {
		t.Fatalf("expected userId %q, got %q", "user-1", createResp.UserID)
	}
	if createResp.AssetID == "" {
		t.Fatalf("expected non-empty assetId")
	}

	// --- 2. List favourites and verify shape ---

	listReq, err := http.NewRequest(http.MethodGet, baseURL+"/v1/users/user-1/favourites", nil)
	if err != nil {
		t.Fatalf("create GET request: %v", err)
	}
	listReq.Header.Set("X-User-ID", "user-1")

	listResp, err := client.Do(listReq)
	if err != nil {
		t.Fatalf("GET favourites request error: %v", err)
	}
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", listResp.StatusCode)
	}

	var listBody struct {
		Items []struct {
			UserID    string `json:"userId"`
			AssetID   string `json:"assetId"`
			CreatedAt string `json:"createdAt"`
			Asset     struct {
				ID          string `json:"id"`
				Type        string `json:"type"`
				Description string `json:"description"`
			} `json:"asset"`
		} `json:"items"`
		Meta struct {
			Total int `json:"total"`
			Count int `json:"count"`
		} `json:"meta"`
	}

	if err := json.NewDecoder(listResp.Body).Decode(&listBody); err != nil {
		t.Fatalf("decode list response: %v", err)
	}

	if listBody.Meta.Total != 1 || listBody.Meta.Count != 1 {
		t.Fatalf("expected meta.total=1, meta.count=1, got total=%d, count=%d",
			listBody.Meta.Total, listBody.Meta.Count)
	}

	if len(listBody.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(listBody.Items))
	}

	item := listBody.Items[0]
	if item.UserID != "user-1" {
		t.Fatalf("expected item.userId %q, got %q", "user-1", item.UserID)
	}
	if item.AssetID != createResp.AssetID {
		t.Fatalf("expected item.assetId %q, got %q", createResp.AssetID, item.AssetID)
	}
	if item.Asset.Type != "insight" {
		t.Fatalf("expected asset.type %q, got %q", "insight", item.Asset.Type)
	}
	if item.Asset.Description != "An example insight" {
		t.Fatalf("unexpected asset.description %q", item.Asset.Description)
	}
}

// Error-path integration: invalid type query parameter should return
// a Problem Details 400 with the expected shape.
func TestFavouritesIntegration_ListInvalidType(t *testing.T) {
	t.Parallel()

	router := app.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()
	baseURL := server.URL

	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/users/user-1/favourites?type=invalid-type", nil)
	if err != nil {
		t.Fatalf("create GET request: %v", err)
	}
	req.Header.Set("X-User-ID", "user-1")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET favourites with invalid type error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/problem+json") {
		t.Fatalf("expected Content-Type application/problem+json, got %q", ct)
	}

	var problem struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Status   int    `json:"status"`
		Detail   string `json:"detail"`
		Instance string `json:"instance"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem response: %v", err)
	}

	if problem.Status != http.StatusBadRequest {
		t.Fatalf("expected problem.status %d, got %d", http.StatusBadRequest, problem.Status)
	}
	if problem.Type != "invalid_type" {
		t.Fatalf("expected problem.type %q, got %q", "invalid_type", problem.Type)
	}
	if problem.Title == "" {
		t.Fatalf("expected non-empty problem.title")
	}
	if problem.Instance == "" {
		t.Fatalf("expected non-empty problem.instance")
	}
}

// Error-path integration: missing X-User-ID header should return 401.
func TestFavouritesIntegration_AuthMissingHeader(t *testing.T) {
	t.Parallel()

	router := app.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()
	baseURL := server.URL

	resp, err := client.Get(baseURL + "/v1/users/user-1/favourites")
	if err != nil {
		t.Fatalf("GET favourites without auth error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/problem+json") {
		t.Fatalf("expected Content-Type application/problem+json, got %q", ct)
	}
}

// Error-path integration: X-User-ID header present but does not match {userId}
// should return 403 Forbidden.
func TestFavouritesIntegration_AuthUserMismatch(t *testing.T) {
	t.Parallel()

	router := app.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()
	baseURL := server.URL

	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/users/user-1/favourites", nil)
	if err != nil {
		t.Fatalf("create GET request: %v", err)
	}
	// Authenticated as user-2, but requesting /users/user-1/...
	req.Header.Set("X-User-ID", "user-2")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET favourites with user mismatch error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/problem+json") {
		t.Fatalf("expected Content-Type application/problem+json, got %q", ct)
	}

	var problem struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Status   int    `json:"status"`
		Detail   string `json:"detail"`
		Instance string `json:"instance"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem response: %v", err)
	}

	if problem.Status != http.StatusForbidden {
		t.Fatalf("expected problem.status %d, got %d", http.StatusForbidden, problem.Status)
	}
	if problem.Type != "forbidden" {
		t.Fatalf("expected problem.type %q, got %q", "forbidden", problem.Type)
	}
	if problem.Title == "" {
		t.Fatalf("expected non-empty problem.title")
	}
	if problem.Instance == "" {
		t.Fatalf("expected non-empty problem.instance")
	}
}
