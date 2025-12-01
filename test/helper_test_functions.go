package tests

import (
	"net/http/httptest"
	"testing"
	"encoding/json"
)

func testErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedMessage, expectedErrorCode string) {
	if w.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, w.Code)
	}

	// Check JSON response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if message := response["message"]; message != expectedMessage  {
		t.Errorf("Response message: %v", message)
	}

	if errorCode := response["error_code"]; errorCode != expectedErrorCode {
		t.Errorf("Response error_code: %v", errorCode)
	}
}