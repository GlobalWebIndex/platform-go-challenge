package gwiexercise_test

import (
	"encoding/json"
	"net/http"
	"testing"

	gwiExercise "platform-go-challenge/gwiexercise"
)

func TestBulkCreateOnly(t *testing.T) {
	store := gwiExercise.NewInMemoryStore()
	s := gwiExercise.NewServer(store)
	user := "u3"

	payload := []map[string]any{
		{
			"type":        "insight",
			"description": "a1",
			"payload":     map[string]any{"text": "hello"},
		},
		{
			"type":        "chart",
			"description": "a2",
			"payload":     map[string]any{"title": "t", "xAxisTitle": "x", "yAxisTitle": "y", "data": []float64{1}},
		},
	}
	w := doReq(s, http.MethodPost, "/users/"+user+"/favourites/bulk", payload)
	if w.Code != http.StatusCreated {
		t.Fatalf("bulk expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var created []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	if len(created) != 2 {
		t.Fatalf("expected 2 created, got %d", len(created))
	}
}

func TestDeleteNotFound(t *testing.T) {
	store := gwiExercise.NewInMemoryStore()
	s := gwiExercise.NewServer(store)
	user := "u4"

	w := doReq(s, http.MethodDelete, "/users/"+user+"/favourites/nope", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestPATCHValidationCases(t *testing.T) {
	store := gwiExercise.NewInMemoryStore()
	s := gwiExercise.NewServer(store)
	user := "ux"

	tests := []struct {
		name   string
		body   any
		status int
	}{
		{"empty string", map[string]string{"description": ""}, http.StatusBadRequest},
		{"only spaces", map[string]string{"description": "   "}, http.StatusBadRequest},
		{"missing field", map[string]any{}, http.StatusBadRequest},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := doReq(s, http.MethodPatch, "/users/"+user+"/favourites/anyid", tc.body)
			if w.Code != tc.status {
				t.Fatalf("expected %d, got %d: %s", tc.status, w.Code, w.Body.String())
			}
		})
	}
}

func TestGETPaginationLimits(t *testing.T) {
	store := gwiExercise.NewInMemoryStore()
	s := gwiExercise.NewServer(store)
	user := "u5"

	// Seed 2 assets
	payload := []map[string]any{
		{"type": "insight", "description": "a1", "payload": map[string]any{"text": "t1"}},
		{"type": "chart", "description": "a2", "payload": map[string]any{"title": "t", "xAxisTitle": "x", "yAxisTitle": "y", "data": []float64{1}}},
	}
	w := doReq(s, http.MethodPost, "/users/"+user+"/favourites/bulk", payload)
	if w.Code != http.StatusCreated {
		t.Fatalf("seed bulk expected 201, got %d: %s", w.Code, w.Body.String())
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"default params", "", 2},
		{"limit=1", "?limit=1&offset=0", 1},
		{"limit=0 -> default", "?limit=0", 2},
		{"limit=-5 -> default", "?limit=-5", 2},
		{"offset beyond length -> empty", "?offset=10", 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := doReq(s, http.MethodGet, "/users/"+user+"/favourites"+tc.query, nil)
			if w.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
			}
			var list []map[string]any
			_ = json.Unmarshal(w.Body.Bytes(), &list)
			if len(list) != tc.expected {
				t.Fatalf("expected %d items, got %d", tc.expected, len(list))
			}
		})
	}
}
