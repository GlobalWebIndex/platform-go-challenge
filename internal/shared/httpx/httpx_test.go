package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	type payload struct {
		Message string `json:"message"`
	}

	rr := httptest.NewRecorder()
	body := payload{Message: "hello"}

	WriteJSON(rr, http.StatusCreated, body)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", got)
	}

	var decoded payload
	if err := json.NewDecoder(rr.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if decoded.Message != "hello" {
		t.Fatalf("expected message %q, got %q", "hello", decoded.Message)
	}
}

func TestWriteProblem(t *testing.T) {
	const (
		status   = http.StatusUnprocessableEntity
		typ      = "invalid_asset"
		title    = "Invalid asset"
		detail   = "asset payload is invalid"
		instance = "/v1/users/user-1/favourites"
	)

	rr := httptest.NewRecorder()
	WriteProblem(rr, status, typ, title, detail, instance)

	if rr.Code != status {
		t.Fatalf("expected status %d, got %d", status, rr.Code)
	}

	if got := rr.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Fatalf("expected Content-Type application/problem+json, got %q", got)
	}

	var p ProblemDetail
	if err := json.NewDecoder(rr.Body).Decode(&p); err != nil {
		t.Fatalf("failed to decode problem JSON: %v", err)
	}

	if p.Status != status {
		t.Fatalf("expected problem Status %d, got %d", status, p.Status)
	}
	if p.Type != typ {
		t.Fatalf("expected problem Type %q, got %q", typ, p.Type)
	}
	if p.Title != title {
		t.Fatalf("expected problem Title %q, got %q", title, p.Title)
	}
	if p.Detail != detail {
		t.Fatalf("expected problem Detail %q, got %q", detail, p.Detail)
	}
	if p.Instance != instance {
		t.Fatalf("expected problem Instance %q, got %q", instance, p.Instance)
	}
}
