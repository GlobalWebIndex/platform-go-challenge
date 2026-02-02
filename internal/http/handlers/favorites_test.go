package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kongrillas/platform-go-challenge/internal/http/router"
	"github.com/kongrillas/platform-go-challenge/internal/store"
)

func newTestServer() *httptest.Server {
	s := store.NewMemoryStore()
	h := NewFavoritesHandler(s)
	return httptest.NewServer(router.New(h))
}

func TestAddAndListFavorites(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	body := `{"id":"a1","type":"insight","description":"test","payload":{"text":"hello"}}`
	res, err := http.Post(ts.URL+"/users/u1/favorites", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}

	res2, err := http.Get(ts.URL + "/users/u1/favorites")
	if err != nil {
		t.Fatal(err)
	}
	if res2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res2.StatusCode)
	}
}

func TestAddInvalidPayload(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	// insight without text
	body := `{"id":"a1","type":"insight","payload":{"text":""}}`
	res, err := http.Post(ts.URL+"/users/u1/favorites", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestUpdateDescription(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	body := `{"id":"a1","type":"insight","description":"test","payload":{"text":"hello"}}`
	res, err := http.Post(ts.URL+"/users/u1/favorites", "application/json", bytes.NewBufferString(body))
	if err != nil || res.StatusCode != http.StatusCreated {
		t.Fatalf("setup failed, status=%v err=%v", res.StatusCode, err)
	}

	req := httptest.NewRequest(http.MethodPatch, ts.URL+"/users/u1/favorites/a1", bytes.NewBufferString(`{"description":"updated"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Use handler directly would require path parsing; easiest: real client
	client := &http.Client{}
	req2, _ := http.NewRequest(http.MethodPatch, ts.URL+"/users/u1/favorites/a1", bytes.NewBufferString(`{"description":"updated"}`))
	req2.Header.Set("Content-Type", "application/json")
	res2, err := client.Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	if res2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res2.StatusCode)
	}

	_ = rr // silence unused if kept
	_ = req
}

func TestDeleteFavorite(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	body := `{"id":"a1","type":"insight","description":"test","payload":{"text":"hello"}}`
	res, err := http.Post(ts.URL+"/users/u1/favorites", "application/json", bytes.NewBufferString(body))
	if err != nil || res.StatusCode != http.StatusCreated {
		t.Fatalf("setup failed, status=%v err=%v", res.StatusCode, err)
	}

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/users/u1/favorites/a1", nil)
	client := &http.Client{}
	res2, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res2.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", res2.StatusCode)
	}
}
