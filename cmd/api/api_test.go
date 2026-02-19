package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_Contract_GoldenPath(t *testing.T) {
	app := newContractTestApplication(t)
	mux := app.mount()

	userID := "user-123"

	// 1) Add favourite -> { "data": { ... } }
	addBody := []byte(`{
		"asset_id":"asset-1",
		"type":"chart",
		"description":"Initial description",
		"title":"Sales Chart",
		"x_axis_title":"Month",
		"y_axis_title":"Revenue",
		"data":[]
	}`)

	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/users/"+userID+"/favourites", bytes.NewBuffer(addBody))
	req1.Header.Set("Content-Type", "application/json")

	rr1 := executeRequest(req1, mux)
	checkResponseCode(t, http.StatusCreated, rr1.Code)

	asset := envelopeObject(t, rr1)
	assetID, _ := asset["asset_id"].(string)
	assert.Equal(t, "asset-1", assetID)

	// 2) List favourites (paginated) -> { "data": [...], "meta": {...} }
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/users/"+userID+"/favourites?limit=20", nil)
	rr2 := executeRequest(req2, mux)
	checkResponseCode(t, http.StatusOK, rr2.Code)

	items, meta := envelopePage(t, rr2)
	assert.GreaterOrEqual(t, len(items), 1)

	assert.Equal(t, 20, int(meta["limit"].(float64)))
	assert.Equal(t, false, meta["has_more"].(bool)) // with 1 item should be false

	// 3) Get specific favourite
	req3, _ := http.NewRequest(http.MethodGet, "/api/v1/users/"+userID+"/favourites/"+assetID, nil)
	rr3 := executeRequest(req3, mux)
	checkResponseCode(t, http.StatusOK, rr3.Code)

	obj := envelopeObject(t, rr3)
	assert.Equal(t, "asset-1", obj["asset_id"])

	// 4) Update description
	descBody, _ := json.Marshal(map[string]string{"description": "New desc"})
	req4, _ := http.NewRequest(http.MethodPut, "/api/v1/users/"+userID+"/favourites/"+assetID+"/description", bytes.NewBuffer(descBody))
	req4.Header.Set("Content-Type", "application/json")

	rr4 := executeRequest(req4, mux)
	checkResponseCode(t, http.StatusNoContent, rr4.Code)

	// 5) Verify description updated
	req5, _ := http.NewRequest(http.MethodGet, "/api/v1/users/"+userID+"/favourites/"+assetID, nil)
	rr5 := executeRequest(req5, mux)
	checkResponseCode(t, http.StatusOK, rr5.Code)

	obj2 := envelopeObject(t, rr5)
	assert.Equal(t, "New desc", obj2["description"])

	// 6) Delete favourite
	req6, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/"+userID+"/favourites/"+assetID, nil)
	rr6 := executeRequest(req6, mux)
	checkResponseCode(t, http.StatusNoContent, rr6.Code)
}

func TestAPI_Contract_PaginationThroughManyfavourites_KeysetCursor(t *testing.T) {
	app := newContractTestApplication(t)
	mux := app.mount()

	userID := "user-scale-1"

	// Add many favourites with unique external asset_ids
	const total = 250
	for i := 0; i < total; i++ {
		body := []byte(fmt.Sprintf(`{
			"asset_id":"asset-%d",
			"type":"chart",
			"title":"Test Chart",
			"data":[]
		}`, i))

		r, _ := http.NewRequest(http.MethodPost, "/api/v1/users/"+userID+"/favourites", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		resp := executeRequest(r, mux)
		if resp.Code != http.StatusCreated {
			t.Fatalf("add favourite failed: status=%d body=%s", resp.Code, resp.Body.String())
		}
	}

	// Paginate through all results
	seen := make(map[string]bool, total)
	cursor := ""
	limit := 50
	pages := 0

	for {
		u := "/api/v1/users/" + userID + "/favourites?limit=" + fmt.Sprintf("%d", limit)
		if cursor != "" {
			// cursor contains ':' and '|', must be URL-escaped
			u += "&cursor=" + url.QueryEscape(cursor)
		}

		r, _ := http.NewRequest(http.MethodGet, u, nil)
		resp := executeRequest(r, mux)
		checkResponseCode(t, http.StatusOK, resp.Code)

		items, meta := envelopePage(t, resp)

		gotLimit := int(meta["limit"].(float64))
		assert.Equal(t, limit, gotLimit)

		// collect asset_ids, ensure no duplicates
		for _, it := range items {
			obj, ok := it.(map[string]any)
			if !ok {
				t.Fatalf("expected item to be object, got %T", it)
			}

			assetID, _ := obj["asset_id"].(string)
			if assetID == "" {
				t.Fatalf("missing item.asset_id")
			}

			if seen[assetID] {
				t.Fatalf("duplicate asset_id seen across pages: %s", assetID)
			}
			seen[assetID] = true
		}

		hasMore, ok := meta["has_more"].(bool)
		if !ok {
			t.Fatalf("expected meta.has_more bool, got %T", meta["has_more"])
		}

		nextCursor := ""
		if v, ok := meta["next_cursor"]; ok && v != nil {
			nextCursor, _ = v.(string)
		}

		pages++
		if pages > 50 {
			t.Fatalf("pagination exceeded 50 pages; cursor=%q", cursor)
		}

		if !hasMore {
			break
		}

		if nextCursor == "" {
			t.Fatalf("has_more=true but next_cursor is empty")
		}
		cursor = nextCursor
	}

	assert.Equal(t, total, len(seen), "should retrieve all favourites via pagination")
}

func TestAPI_Contract_InvalidCursor_Returns400(t *testing.T) {
	app := newContractTestApplication(t)
	mux := app.mount()

	userID := "user-123"

	u := "/api/v1/users/" + userID + "/favourites?limit=10&cursor=" + url.QueryEscape("not-a-cursor")

	r, _ := http.NewRequest(http.MethodGet, u, nil)
	resp := executeRequest(r, mux)

	checkResponseCode(t, http.StatusBadRequest, resp.Code)
}
