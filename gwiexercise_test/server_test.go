package gwiexercise_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	gwiExercise "platform-go-challenge/gwiexercise"
)

func doReq(s *gwiExercise.Server, method, path string, body any) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	}
	r := httptest.NewRequest(method, path, reader)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	return w
}

func TestAddListUpdateDeleteFlow(t *testing.T) {
	store := gwiExercise.NewInMemoryStore()
	s := gwiExercise.NewServer(store)
	user := "u123"

	// Initially empty list
	w := doReq(s, http.MethodGet, "/users/"+user+"/favourites", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %d", len(list))
	}

	addReq := map[string]any{
		"type":        "chart",
		"description": "My chart",
		"payload": map[string]any{
			"title":      "Sales",
			"xAxisTitle": "Month",
			"yAxisTitle": "USD",
			"data":       []float64{1, 2, 3},
		},
	}
	w = doReq(s, http.MethodPost, "/users/"+user+"/favourites", addReq)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var created map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	id := created["id"].(string)
	if id == "" {
		t.Fatalf("expected id assigned")
	}

	// List and verify one item
	w = doReq(s, http.MethodGet, "/users/"+user+"/favourites", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200")
	}
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}

	// Update description
	w = doReq(s, http.MethodPatch, "/users/"+user+"/favourites/"+id, map[string]string{"description": "Updated"})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d: %s", w.Code, w.Body.String())
	}

	// Delete
	w = doReq(s, http.MethodDelete, "/users/"+user+"/favourites/"+id, nil)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 on delete, got %d", w.Code)
	}

	// List again empty
	w = doReq(s, http.MethodGet, "/users/"+user+"/favourites", nil)
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 0 {
		t.Fatalf("expected empty after deletion")
	}
}

func TestPOSTValidationCases(t *testing.T) {
	store := gwiExercise.NewInMemoryStore()
	s := gwiExercise.NewServer(store)
	user := "u1"

	tests := []struct {
		name   string
		body   any
		status int
	}{
		{
			name:   "missing fields",
			body:   map[string]any{},
			status: http.StatusBadRequest,
		},
		{
			name: "unsupported type",
			body: map[string]any{
				"type":        "unknown",
				"description": "d",
				"payload":     map[string]any{"x": 1},
			},
			status: http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := doReq(s, http.MethodPost, "/users/"+user+"/favourites", tc.body)
			if w.Code != tc.status {
				t.Fatalf("expected %d, got %d: %s", tc.status, w.Code, w.Body.String())
			}
		})
	}
}
