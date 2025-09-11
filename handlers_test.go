package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"

	"github.com/google/uuid"
)

func resetStore() {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.users = make(map[uuid.UUID]*User)
}

func TestHandleFavourites(t *testing.T) {
	resetStore()
	userID := uuid.New()
	user := &User{ID: userID}
	chart := &Chart{ID: uuid.New(), Title: "Chart1", Favorite: true}
	insight := &Insight{ID: uuid.New(), Text: "Insight1", Favorite: false}
	user.Favourites = []Asset{chart, insight}
	store.AddUser(user)
	token, _ := GenerateJWT(userID)
	req := httptest.NewRequest("GET", "/favourites?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	AuthMiddleware(handleFavourites)(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("Expected 1 favourite, got %d", len(resp))
	}
	// Assert that the returned asset matches the created chart
	if title, ok := resp[0]["Title"].(string); !ok || title != "Chart1" {
		t.Errorf("Expected asset Title to be 'Chart1', got %v", resp[0]["Title"])
	}

	// Test 1b: Pagination with limit=1, offset=0
	reqPag := httptest.NewRequest("GET", "/favourites?limit=1&offset=0", nil)
	reqPag.Header.Set("Authorization", "Bearer "+token)
	wPag := httptest.NewRecorder()
	AuthMiddleware(handleFavourites)(wPag, reqPag)
	if wPag.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", wPag.Code)
	}
	var respPag []map[string]interface{}
	if err := json.Unmarshal(wPag.Body.Bytes(), &respPag); err != nil {
		t.Fatalf("Failed to decode paginated response: %v", err)
	}
	if len(respPag) != 1 {
		t.Fatalf("Expected 1 favourite with pagination, got %d", len(respPag))
	}
	if title, ok := respPag[0]["Title"].(string); !ok || title != "Chart1" {
		t.Errorf("Expected paginated asset Title to be 'Chart1', got %v", respPag[0]["Title"])
	}

	// Test 1c: Pagination with limit=1, offset=1 (should be empty)
	reqPag2 := httptest.NewRequest("GET", "/favourites?limit=1&offset=1", nil)
	reqPag2.Header.Set("Authorization", "Bearer "+token)
	wPag2 := httptest.NewRecorder()
	AuthMiddleware(handleFavourites)(wPag2, reqPag2)
	if wPag2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", wPag2.Code)
	}
	var respPag2 []map[string]interface{}
	if err := json.Unmarshal(wPag2.Body.Bytes(), &respPag2); err != nil {
		t.Fatalf("Failed to decode paginated response: %v", err)
	}
	if len(respPag2) != 0 {
		t.Fatalf("Expected 0 favourites with pagination offset, got %d", len(respPag2))
	}

	// Test 2: User with no favorite assets
	user2ID := uuid.New()
	user2 := &User{ID: user2ID}
	user2.Favourites = []Asset{&Insight{ID: uuid.New(), Text: "Insight2", Favorite: false}}
	store.AddUser(user2)
	token2, _ := GenerateJWT(user2ID)
	req2 := httptest.NewRequest("GET", "/favourites?limit=10&offset=0", nil)
	req2.Header.Set("Authorization", "Bearer "+token2)
	w2 := httptest.NewRecorder()
	AuthMiddleware(handleFavourites)(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}
	var resp2 []map[string]interface{}
	if err := json.Unmarshal(w2.Body.Bytes(), &resp2); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(resp2) != 0 {
		t.Fatalf("Expected 0 favourites, got %d", len(resp2))
	}

	// Test 2b: Pagination with limit=1, offset=0 (should be empty)
	req2Pag := httptest.NewRequest("GET", "/favourites?limit=1&offset=0", nil)
	req2Pag.Header.Set("Authorization", "Bearer "+token2)
	w2Pag := httptest.NewRecorder()
	AuthMiddleware(handleFavourites)(w2Pag, req2Pag)
	if w2Pag.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2Pag.Code)
	}
	var resp2Pag []map[string]interface{}
	if err := json.Unmarshal(w2Pag.Body.Bytes(), &resp2Pag); err != nil {
		t.Fatalf("Failed to decode paginated response: %v", err)
	}
	if len(resp2Pag) != 0 {
		t.Fatalf("Expected 0 favourites with pagination, got %d", len(resp2Pag))
	}
}

func TestHandleAddFavourite_Chart(t *testing.T) {
	resetStore()
	userID := uuid.New()
	store.AddUser(&User{ID: userID})

	cases := []struct {
		name     string
		favorite bool
	}{
		{"Chart Favorite True", true},
		{"Chart Favorite False", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chart := map[string]interface{}{
				"id":          uuid.Nil.String(),
				"title":       "Test Chart",
				"description": "Chart Desc",
			}
			assetBytes, _ := json.Marshal(chart)
			body := map[string]interface{}{
				"type":     ChartType,
				"asset":    json.RawMessage(assetBytes),
				"favorite": tc.favorite,
			}
			bodyBytes, _ := json.Marshal(body)

			token, _ := GenerateJWT(userID)
			req := httptest.NewRequest(http.MethodPost, "/favourites/add", bytes.NewReader(bodyBytes))
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			AuthMiddleware(handleAddFavourite)(w, req)

			if w.Code != http.StatusCreated {
				t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
			}

			var resp Chart
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if resp.Title != "Test Chart" {
				t.Errorf("expected title 'Test Chart', got '%s'", resp.Title)
			}
			if resp.Description != "Chart Desc" {
				t.Errorf("expected description 'Chart Desc', got '%s'", resp.Description)
			}
			if resp.Favorite != tc.favorite {
				t.Errorf("expected favorite %v, got %v", tc.favorite, resp.Favorite)
			}
			if resp.ID == uuid.Nil {
				t.Error("expected non-nil UUID for chart asset")
			}
		})
	}
}

func TestHandleAddFavourite_Insight(t *testing.T) {
	resetStore()
	userID := uuid.New()
	store.AddUser(&User{ID: userID})

	cases := []struct {
		name     string
		favorite bool
	}{
		{"Insight Favorite True", true},
		{"Insight Favorite False", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			insight := map[string]interface{}{
				"id":          uuid.Nil.String(),
				"text":        "Test Insight",
				"description": "Insight Desc",
			}
			assetBytes, _ := json.Marshal(insight)
			body := map[string]interface{}{
				"type":     InsightType,
				"asset":    json.RawMessage(assetBytes),
				"favorite": tc.favorite,
			}
			bodyBytes, _ := json.Marshal(body)

			token, _ := GenerateJWT(userID)
			req := httptest.NewRequest(http.MethodPost, "/favourites/add", bytes.NewReader(bodyBytes))
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			AuthMiddleware(handleAddFavourite)(w, req)

			if w.Code != http.StatusCreated {
				t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
			}

			var resp Insight
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if resp.Text != "Test Insight" {
				t.Errorf("expected text 'Test Insight', got '%s'", resp.Text)
			}
			if resp.Description != "Insight Desc" {
				t.Errorf("expected description 'Insight Desc', got '%s'", resp.Description)
			}
			if resp.Favorite != tc.favorite {
				t.Errorf("expected favorite %v, got %v", tc.favorite, resp.Favorite)
			}
			if resp.ID == uuid.Nil {
				t.Error("expected non-nil UUID for insight asset")
			}
		})
	}
}

func TestHandleAddFavourite_Audience(t *testing.T) {
	resetStore()
	userID := uuid.New()
	store.AddUser(&User{ID: userID})

	cases := []struct {
		name     string
		favorite bool
	}{
		{"Audience Favorite True", true},
		{"Audience Favorite False", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			audience := map[string]interface{}{
				"id":          uuid.Nil.String(),
				"description": "Audience Desc",
			}
			assetBytes, _ := json.Marshal(audience)
			body := map[string]interface{}{
				"type":     AudienceType,
				"asset":    json.RawMessage(assetBytes),
				"favorite": tc.favorite,
			}
			bodyBytes, _ := json.Marshal(body)

			token, _ := GenerateJWT(userID)
			req := httptest.NewRequest(http.MethodPost, "/favourites/add", bytes.NewReader(bodyBytes))
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			AuthMiddleware(handleAddFavourite)(w, req)

			if w.Code != http.StatusCreated {
				t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
			}

			var resp Audience
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if resp.Description != "Audience Desc" {
				t.Errorf("expected description 'Audience Desc', got '%s'", resp.Description)
			}
			if resp.Favorite != tc.favorite {
				t.Errorf("expected favorite %v, got %v", tc.favorite, resp.Favorite)
			}
			if resp.ID == uuid.Nil {
				t.Error("expected non-nil UUID for audience asset")
			}
		})
	}
}

func TestHandleRemoveFavourite_TrueToFalse(t *testing.T) {
	resetStore()
	userID := uuid.New()
	chart := &Chart{ID: uuid.New(), Title: "Chart1", Favorite: true}
	user := &User{ID: userID, Favourites: []Asset{chart}}
	store.AddUser(user)

	reqBody := map[string]interface{}{"favorite": false}
	bodyBytes, _ := json.Marshal(reqBody)
	token, _ := GenerateJWT(userID)
	req := httptest.NewRequest(http.MethodPut, "/favourites/remove?asset_id="+chart.ID.String(), bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	AuthMiddleware(handleRemoveFavourite)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var resp Chart
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Favorite != false {
		t.Errorf("expected favorite false, got %v", resp.Favorite)
	}
}

func TestHandleRemoveFavourite_FalseToTrue(t *testing.T) {
	resetStore()
	userID := uuid.New()
	chart := &Chart{ID: uuid.New(), Title: "Chart2", Favorite: false}
	user := &User{ID: userID, Favourites: []Asset{chart}}
	store.AddUser(user)

	reqBody := map[string]interface{}{"favorite": true}
	bodyBytes, _ := json.Marshal(reqBody)
	token, _ := GenerateJWT(userID)
	req := httptest.NewRequest(http.MethodPut, "/favourites/remove?asset_id="+chart.ID.String(), bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	AuthMiddleware(handleRemoveFavourite)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var resp Chart
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Favorite != true {
		t.Errorf("expected favorite true, got %v", resp.Favorite)
	}
}

func TestHandleEditFavourite_ChangeDescription(t *testing.T) {
	resetStore()
	userID := uuid.New()
	chart := &Chart{ID: uuid.New(), Title: "Chart3", Description: "Old Description", Favorite: true}
	user := &User{ID: userID, Favourites: []Asset{chart}}
	store.AddUser(user)

	newDesc := "New Description"
	reqBody := map[string]interface{}{"description": newDesc}
	bodyBytes, _ := json.Marshal(reqBody)
	token, _ := GenerateJWT(userID)
	req := httptest.NewRequest(http.MethodPut, "/favourites/edit?asset_id="+chart.ID.String(), bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	AuthMiddleware(handleEditFavourite)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var resp Chart
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Description != newDesc {
		t.Errorf("expected description '%s', got '%s'", newDesc, resp.Description)
	}
}

func TestHandleDeleteFavourite_DeleteAsset(t *testing.T) {
	resetStore()
	userID := uuid.New()
	chart := &Chart{ID: uuid.New(), Title: "Chart4", Favorite: true}
	insight := &Insight{ID: uuid.New(), Text: "Insight4", Favorite: true}
	user := &User{ID: userID, Favourites: []Asset{chart, insight}}
	store.AddUser(user)

	// Check favourites before deletion
	token, _ := GenerateJWT(userID)
	reqList := httptest.NewRequest(http.MethodGet, "/favourites?limit=10&offset=0", nil)
	reqList.Header.Set("Authorization", "Bearer "+token)
	wList := httptest.NewRecorder()
	AuthMiddleware(handleFavourites)(wList, reqList)
	if wList.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, wList.Code)
	}
	var respList []map[string]interface{}
	if err := json.NewDecoder(wList.Body).Decode(&respList); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(respList) != 2 {
		t.Fatalf("expected 2 assets before deletion, got %d", len(respList))
	}

	req := httptest.NewRequest(http.MethodDelete, "/favourites/delete?asset_id="+chart.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	AuthMiddleware(handleDeleteFavourite)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var resp []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 asset remaining, got %d", len(resp))
	}
	if text, ok := resp[0]["Text"].(string); !ok || text != "Insight4" {
		t.Errorf("expected remaining asset to be 'Insight4', got %v", resp[0]["Text"])
	}

	// Check favourites after deletion
	reqList2 := httptest.NewRequest(http.MethodGet, "/favourites?limit=10&offset=0", nil)
	reqList2.Header.Set("Authorization", "Bearer "+token)
	wList2 := httptest.NewRecorder()
	AuthMiddleware(handleFavourites)(wList2, reqList2)
	if wList2.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, wList2.Code)
	}
	var respList2 []map[string]interface{}
	if err := json.NewDecoder(wList2.Body).Decode(&respList2); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(respList2) != 1 {
		t.Fatalf("expected 1 asset after deletion, got %d", len(respList2))
	}
	if text, ok := respList2[0]["Text"].(string); !ok || text != "Insight4" {
		t.Errorf("expected remaining asset to be 'Insight4', got %v", respList2[0]["Text"])
	}
}
