package http

import (
	"bytes"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	amem "github.com/gwi/platform-go-challenge/assets/adapters/memory"
	"github.com/gwi/platform-go-challenge/favourites/adapters/assetcatalog"
	"github.com/gwi/platform-go-challenge/favourites/adapters/memory"
	"github.com/gwi/platform-go-challenge/favourites/application"
)


const testBasePath = "/v1"

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	signer, verifier := NewHMACAuth([]byte("test-secret"))
	repo := memory.NewRepository()
	catalog := amem.NewInMemoryCatalog()
	assetCatalogClient := assetcatalog.NewClient(catalog)
	svc := application.NewFavouritesService(repo, assetCatalogClient)
	openAPIServer := &OpenAPIServer{
		UC:          svc,
		EventsUC:    svc,
		Signer:      signer,
		Verifier:    verifier,
		TokenExpiry: time.Hour,
	}
	mux := stdhttp.NewServeMux()
	handler := HandlerFromMuxWithBaseURL(openAPIServer, mux, testBasePath)
	return httptest.NewServer(handler)
}

func getToken(t *testing.T, baseURL string, userID string) string {
	t.Helper()
	body := []byte(`{"user_id":"` + userID + `"}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, baseURL+testBasePath+"/auth/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := stdhttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("auth request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusOK {
		t.Fatalf("auth: expected 200, got %d", resp.StatusCode)
	}
	var out struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode token: %v", err)
	}
	return out.Token
}

func TestOpenAPIServer_AuthToken_and_List(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	req, _ := stdhttp.NewRequest(stdhttp.MethodGet, srv.URL+testBasePath+"/users/u1/favourites", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusOK {
		t.Errorf("list: expected 200, got %d", resp.StatusCode)
	}
}

func TestOpenAPIServer_List_Unauthorized(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	req, _ := stdhttp.NewRequest(stdhttp.MethodGet, srv.URL+testBasePath+"/users/u1/favourites", nil)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", resp.StatusCode)
	}
}

// TestOpenAPIServer_AddFullBody_Get_Update_Delete adds a favourite with full body, gets it, updates description, then deletes.
func TestOpenAPIServer_AddFullBody_Get_Update_Delete(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	addBody := []byte(`{"type":"chart","description":"My chart","chart":{"title":"Usage","x_axis_title":"Month","y_axis_title":"Hours","data":[1,2,3]}}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/u1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusCreated {
		t.Fatalf("add: expected 201, got %d", resp.StatusCode)
	}
	var addResult struct {
		Id *string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&addResult); err != nil || addResult.Id == nil {
		t.Fatalf("add: decode response: %v", err)
	}
	assetID := *addResult.Id

	// GET same asset
	req2, _ := stdhttp.NewRequest(stdhttp.MethodGet, srv.URL+testBasePath+"/users/u1/favourites/"+assetID, nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	resp2, err := srv.Client().Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != stdhttp.StatusOK {
		t.Fatalf("get: expected 200, got %d", resp2.StatusCode)
	}

	// PATCH description
	patchBody := []byte(`{"description":"Updated description"}`)
	req3, _ := stdhttp.NewRequest(stdhttp.MethodPatch, srv.URL+testBasePath+"/users/u1/favourites/"+assetID, bytes.NewReader(patchBody))
	req3.Header.Set("Authorization", "Bearer "+token)
	req3.Header.Set("Content-Type", "application/json")
	resp3, err := srv.Client().Do(req3)
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != stdhttp.StatusOK {
		t.Fatalf("patch: expected 200, got %d", resp3.StatusCode)
	}

	// DELETE
	req4, _ := stdhttp.NewRequest(stdhttp.MethodDelete, srv.URL+testBasePath+"/users/u1/favourites/"+assetID, nil)
	req4.Header.Set("Authorization", "Bearer "+token)
	resp4, err := srv.Client().Do(req4)
	if err != nil {
		t.Fatal(err)
	}
	defer resp4.Body.Close()
	if resp4.StatusCode != stdhttp.StatusNoContent {
		t.Fatalf("delete: expected 204, got %d", resp4.StatusCode)
	}
}

// TestOpenAPIServer_AddByReference adds a favourite by source_asset_id + type (catalog has chr-001 for user-1).
func TestOpenAPIServer_AddByReference(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "user-1")
	addBody := []byte(`{"type":"chart","source_asset_id":"chr-001"}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/user-1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusCreated {
		t.Fatalf("add by reference: expected 201, got %d", resp.StatusCode)
	}
	var result struct {
		Chart *struct {
			Title *string `json:"title"`
		} `json:"chart"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result.Chart == nil || result.Chart.Title == nil || *result.Chart.Title != "Monthly rainfall index" {
		t.Errorf("expected chart title from catalog (Monthly rainfall index), got %v", result.Chart)
	}
}

// TestOpenAPIServer_Forbidden asserts 403 when path userID does not match token subject.
func TestOpenAPIServer_Forbidden(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	req, _ := stdhttp.NewRequest(stdhttp.MethodGet, srv.URL+testBasePath+"/users/u2/favourites", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusForbidden {
		t.Errorf("expected 403 when path user != token subject, got %d", resp.StatusCode)
	}
}

// TestOpenAPIServer_Add_InvalidTypePayload asserts 400 when type requires a payload and none is provided (no source_asset_id).
func TestOpenAPIServer_Add_InvalidTypePayload(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	addBody := []byte(`{"type":"chart"}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/u1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusBadRequest {
		t.Fatalf("expected 400 for type chart without chart payload, got %d", resp.StatusCode)
	}
	var errResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if errResp.Error.Code != ErrCodeBadRequest {
		t.Errorf("expected code BAD_REQUEST, got %s", errResp.Error.Code)
	}
}

// TestOpenAPIServer_Add_SourceAssetNotFound asserts 400 when adding by reference with source_asset_id not in catalog for this user.
func TestOpenAPIServer_Add_SourceAssetNotFound(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "user-1")
	addBody := []byte(`{"type":"chart","source_asset_id":"chr-999"}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/user-1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusBadRequest {
		t.Fatalf("expected 400 when source_asset_id not in catalog, got %d", resp.StatusCode)
	}
	var errResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if errResp.Error.Code != ErrCodeBadRequest {
		t.Errorf("expected code BAD_REQUEST, got %s", errResp.Error.Code)
	}
}

// TestOpenAPIServer_HealthEndpoint tests the health check endpoint.
func TestOpenAPIServer_HealthEndpoint(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	req, _ := stdhttp.NewRequest(stdhttp.MethodGet, srv.URL+testBasePath+"/health", nil)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result struct {
		Status    *string `json:"status"`
		Timestamp *string `json:"timestamp"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result.Status == nil || *result.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %v", result.Status)
	}
	if result.Timestamp == nil {
		t.Error("expected timestamp to be present")
	}
}

// TestOpenAPIServer_Add_DescriptionTooLong tests validation of description length.
func TestOpenAPIServer_Add_DescriptionTooLong(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	longDesc := strings.Repeat("a", 501) // 501 characters
	addBody := []byte(`{"type":"chart","description":"` + longDesc + `","chart":{"title":"T","x_axis_title":"X","y_axis_title":"Y","data":[1,2,3]}}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/u1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusBadRequest {
		t.Fatalf("expected 400 for description too long, got %d", resp.StatusCode)
	}
}

// TestOpenAPIServer_Add_ChartTitleEmpty tests validation of empty chart title.
func TestOpenAPIServer_Add_ChartTitleEmpty(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	addBody := []byte(`{"type":"chart","chart":{"title":"   ","x_axis_title":"X","y_axis_title":"Y","data":[1,2,3]}}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/u1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusBadRequest {
		t.Fatalf("expected 400 for empty chart title, got %d", resp.StatusCode)
	}
}

// TestOpenAPIServer_Add_InsightTextEmpty tests validation of empty insight text.
func TestOpenAPIServer_Add_InsightTextEmpty(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	addBody := []byte(`{"type":"insight","insight":{"text":"   "}}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/u1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusBadRequest {
		t.Fatalf("expected 400 for empty insight text, got %d", resp.StatusCode)
	}
}

// TestOpenAPIServer_Add_AudienceNegativeHours tests validation of negative hours.
func TestOpenAPIServer_Add_AudienceNegativeHours(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	addBody := []byte(`{"type":"audience","audience":{"hours_social_media_daily":-1}}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/u1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusBadRequest {
		t.Fatalf("expected 400 for negative hours, got %d", resp.StatusCode)
	}
}

// TestOpenAPIServer_Add_AudienceExcessiveHours tests validation of hours > 24.
func TestOpenAPIServer_Add_AudienceExcessiveHours(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	addBody := []byte(`{"type":"audience","audience":{"hours_social_media_daily":25}}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/u1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusBadRequest {
		t.Fatalf("expected 400 for hours > 24, got %d", resp.StatusCode)
	}
}

// TestOpenAPIServer_Add_AudienceNegativePurchases tests validation of negative purchases.
func TestOpenAPIServer_Add_AudienceNegativePurchases(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	addBody := []byte(`{"type":"audience","audience":{"purchases_last_month":-5}}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/u1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusBadRequest {
		t.Fatalf("expected 400 for negative purchases, got %d", resp.StatusCode)
	}
}

// TestOpenAPIServer_UpdateDescription_TooLong tests validation when updating description.
func TestOpenAPIServer_UpdateDescription_TooLong(t *testing.T) {
	srv := setupTestServer(t)
	defer srv.Close()

	token := getToken(t, srv.URL, "u1")
	// First add a favourite
	addBody := []byte(`{"type":"chart","chart":{"title":"T","x_axis_title":"X","y_axis_title":"Y","data":[1,2,3]}}`)
	req, _ := stdhttp.NewRequest(stdhttp.MethodPost, srv.URL+testBasePath+"/users/u1/favourites", bytes.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusCreated {
		t.Fatalf("add failed: expected 201, got %d", resp.StatusCode)
	}
	var addResult struct {
		Id *string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&addResult); err != nil || addResult.Id == nil {
		t.Fatalf("add: decode response: %v", err)
	}
	assetID := *addResult.Id

	// Try to update with too long description
	longDesc := strings.Repeat("b", 501)
	patchBody := []byte(`{"description":"` + longDesc + `"}`)
	req2, _ := stdhttp.NewRequest(stdhttp.MethodPatch, srv.URL+testBasePath+"/users/u1/favourites/"+assetID, bytes.NewReader(patchBody))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := srv.Client().Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != stdhttp.StatusBadRequest {
		t.Fatalf("expected 400 for description too long, got %d", resp2.StatusCode)
	}
}
