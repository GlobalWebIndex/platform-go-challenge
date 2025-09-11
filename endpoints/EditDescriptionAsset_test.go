package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ResponseSummaryEdit struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func TestEditDescriptionAsset(t *testing.T) {
	testsEdit := []struct {
		name    string
		body    []byte
		respMsg string
	}{
		{
			name:    "Success edit",
			body:    []byte(`{"asset_id": 2, "new_description": "Annual inflation rate by country"}`),
			respMsg: "Success edit of an asset",
		},
		{
			name:    "Failed edit",
			body:    []byte(`{"asset_id": 20, "new_description": "Annual inflation rate by country"}`),
			respMsg: "Failed to edit asset",
		},
		{
			name:    "Invalid json input",
			body:    []byte(`{"asset_id": ""}`),
			respMsg: "Invalid JSON",
		},
	}

	for _, testEdit := range testsEdit {
		request, _ := http.NewRequest("POST", "/editAssets", bytes.NewBuffer(testEdit.body))
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		EditAssets(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", response.Code)
		}

		var summary ResponseSummaryEdit
		err := json.Unmarshal(response.Body.Bytes(), &summary)
		if err != nil {
			t.Fatalf("error unmarshaling response: %v", err)
		}
		got := summary.Msg
		want := testEdit.respMsg

		if got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}
