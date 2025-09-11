package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ResponseSummaryDelete struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func TestDeleteAssets(t *testing.T) {
	deleteTests := []struct {
		name    string
		body    []byte
		respMsg string
	}{
		{
			name:    "Success deletion",
			body:    []byte(`{"user_id": 1, "asset_id": 2}`),
			respMsg: "Success deletion",
		},
		{
			name:    "failed to delete",
			body:    []byte(`{"user_id": 100, "asset_id": 2}`),
			respMsg: "No favorite asset found for the given attributes",
		},
		{
			name:    "Invalid json input",
			body:    []byte(`{"user_id": ""}`),
			respMsg: "Invalid JSON",
		},
	}

	for _, deleteTest := range deleteTests {
		request, _ := http.NewRequest("POST", "/deleteAssets", bytes.NewBuffer(deleteTest.body))
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		DeleteAsset(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", response.Code)
		}

		var summary ResponseSummaryEdit
		err := json.Unmarshal(response.Body.Bytes(), &summary)
		if err != nil {
			t.Fatalf("error unmarshaling response: %v", err)
		}
		got := summary.Msg
		want := deleteTest.respMsg

		if got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}
