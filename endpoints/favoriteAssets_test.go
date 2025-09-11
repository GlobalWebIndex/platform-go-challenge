package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ResponseSummary struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func TestFavoriteAssets(t *testing.T) {
	tests := []struct {
		name    string
		body    []byte
		respMsg string
	}{
		{
			name:    "valid user_id",
			body:    []byte(`{"user_id": 1}`),
			respMsg: "Success",
		},
		{
			name:    "nil user_id or does not have favorite data",
			body:    []byte(`{"user_id": null}`),
			respMsg: "User has not any favorite asset",
		},
		{
			name:    "Invalid json input",
			body:    []byte(`{"user_id": ""}`),
			respMsg: "Invalid JSON",
		},
	}

	for _, test := range tests {
		request, _ := http.NewRequest("POST", "/favAssets", bytes.NewBuffer(test.body))
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		FavoriteAssets(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", response.Code)
		}

		var summary ResponseSummary
		err := json.Unmarshal(response.Body.Bytes(), &summary)
		if err != nil {
			t.Fatalf("error unmarshaling response: %v", err)
		}
		got := summary.Msg
		want := test.respMsg

		if got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}
