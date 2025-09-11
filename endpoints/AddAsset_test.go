package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ResponseSummaryAdd struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func TestAddAsset(t *testing.T) {
	addTests := []struct {
		name    string
		body    []byte
		respMsg string
	}{
		{
			name:    "New asset",
			body:    []byte(`{"user_id": 1, "asset_id": 1}`),
			respMsg: "New asset added succesfully",
		},
		{
			name:    "Invalid json input",
			body:    []byte(`{"user_id": ""}`),
			respMsg: "Invalid JSON",
		},
	}

	// body := []byte(`{"user_id": 1, "asset_id": 1}`)
	for _, addTest := range addTests {
		request, _ := http.NewRequest("POST", "/deleteAssets", bytes.NewBuffer(addTest.body))
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		AddAsset(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", response.Code)
		}
		var summary ResponseSummaryAdd
		err := json.Unmarshal(response.Body.Bytes(), &summary)
		if err != nil {
			t.Fatalf("error unmarshaling response: %v", err)
		}
		got := summary.Msg
		want := addTest.respMsg

		if got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}
