package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper function: reset global map before each test
func reset() {
	allUserFavourites = make(map[string]map[string]Asset)
}

// Helper function: perform a request against the router
func performRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	rr := httptest.NewRecorder()
	usersHandler(rr, req)
	return rr
}

//
// ───────────────────────────────────────────────────────────────
//   TEST: GET /users/{user_id}/favourites
// ───────────────────────────────────────────────────────────────
//
func TestGetFavourites(t *testing.T) {
	tests := []struct {
		name string
		setup func()
		userID string
		expectedStatus int
	}{
		{
			name: "returns 404 if user does not exist",
			setup: func() {
				reset()
			},
			userID: "missing",
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "returns 200 with empty list if user exists but has no favourites",
			setup: func() {
				reset()
				allUserFavourites["user1"] = map[string]Asset{}
			},
			userID: "user1",
			expectedStatus: http.StatusOK,
		},
		{
			name: "returns 200 with list of favourites",
			setup: func() {
				reset()
				allUserFavourites["user1"] = map[string]Asset{
					"abc": {ID: "abc", Type: "chart"},
				}
			},
			userID: "user1",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			path := "/users/" + tt.userID + "/favourites"
			rr := performRequest("GET", path, nil)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

//
// ───────────────────────────────────────────────────────────────
//   TEST: POST /users/{user_id}/favourites/{asset_type}
// ───────────────────────────────────────────────────────────────
//

func TestPostAddAsset(t *testing.T) {
	tests := []struct {
		name string
		setup func()
		userID string
		assetType string
		body interface{}
		expectedStatus int
	}{
		{
			name: "returns 404 if user does not exist",
			setup: func() { 
				reset() 
			},
			userID: "missing",
			assetType: "chart",
			body: map[string]interface{}{
				"description": "Chart example",
				"data": map[string]interface{}{
					"title": "t",
					"x_axis": "x",
					"y_axis": "y",
					"data": []float64{1, 2},
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "adds chart successfully",
			setup: func() {
				reset()
				allUserFavourites["user1"] = map[string]Asset{}
			},
			userID: "user1",
			assetType: "chart",
			body: map[string]interface{}{
				"description": "Chart example",
				"data": map[string]interface{}{
					"title": "t",
					"x_axis": "x",
					"y_axis": "y",
					"data": []float64{1, 2},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "adds insight successfully",
			setup: func() {
				reset()
				allUserFavourites["user1"] = map[string]Asset{}
			},
			userID: "user1",
			assetType: "insight",
			body: map[string]interface{}{
				"description": "insight 1",
				"data": map[string]interface{}{
					"text": "test",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "adds audience successfully",
			setup: func() {
				reset()
				allUserFavourites["user1"] = map[string]Asset{}
			},
			userID: "user1",
			assetType: "audience",
			body: map[string]interface{}{
				"description": "aud1",
				"data": map[string]interface{}{
					"gender": "Male",
					"birth_country": "GR",
					"age_groups": []string{"18-24"},
					"daily_hours_on_social_media": 2,
					"purchases_last_month": 5,
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			path := "/users/" + tt.userID + "/favourites/" + tt.assetType
			rr := performRequest("POST", path, tt.body)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}


//
// ───────────────────────────────────────────────────────────────
//   TEST: PATCH /users/{user_id}/favourites/{asset_id}
// ───────────────────────────────────────────────────────────────
//

func TestPatchAsset(t *testing.T) {
	tests := []struct {
		name string
		setup func()
		userID string
		assetID string
		body interface{}
		expectedStatus int
	}{
		{
			name: "returns 404 if user does not exist",
			setup: func() { 
				reset() 
			},
			userID: "missing",
			assetID: "asset1",
			body: map[string]interface{}{"description": "x"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "returns 404 if asset does not exist",
			setup: func() {
				reset()
				allUserFavourites["user1"] = map[string]Asset{}
			},
			userID: "user1",
			assetID: "missing",
			body: map[string]interface{}{"description": "x"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "successfully updates description",
			setup: func() {
				reset()
				allUserFavourites["user1"] = map[string]Asset{
					"asset1": {ID: "asset1", Type: "chart", Description: "old"},
				}
			},
			userID: "user1",
			assetID: "asset1",
			body: map[string]interface{}{"description": "new"},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			path := "/users/" + tt.userID + "/favourites/" + tt.assetID
			rr := performRequest("PATCH", path, tt.body)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

//
// ───────────────────────────────────────────────────────────────
//   TEST: DELETE /users/{user_id}/favourites/{asset_id}
// ───────────────────────────────────────────────────────────────
//

func TestDeleteAsset(t *testing.T) {
	tests := []struct {
		name string
		setup func()
		userID string
		assetID string
		expectedStatus int
	}{
		{
			name: "returns 404 if user does not exist",
			setup: func() { 
				reset() 
			},
			userID: "missing",
			assetID: "asset1",
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "returns 404 if asset does not exist",
			setup: func() {
				reset()
				allUserFavourites["user1"] = map[string]Asset{}
			},
			userID: "u1",
			assetID: "missing",
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "successfully deletes existing asset",
			setup: func() {
				reset()
				allUserFavourites["user1"] = map[string]Asset{
					"asset1": {ID: "asset1", Type: "chart"},
				}
			},
			userID: "user1",
			assetID: "asset1",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			path := "/users/" + tt.userID + "/favourites/" + tt.assetID
			rr := performRequest("DELETE", path, nil)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}