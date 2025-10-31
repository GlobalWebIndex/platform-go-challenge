package middleware_test

import (
	"net/http"
	"testing"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/http/middleware"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
)

// MockBodyGetter lets you return any DTO for testing
type MockBodyGetter struct {
	MockedBody    any
	ShouldSucceed bool
}

func (m MockBodyGetter) GetValidatedBody(r *http.Request) (any, bool) {
	return m.MockedBody, m.ShouldSucceed
}
func TestHandler_CreateUser(t *testing.T) {
	// Arrange
	mockUser := dto.CreateUserRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	}
	middleware.Body = MockBodyGetter{
		MockedBody:    mockUser,
		ShouldSucceed: true,
	}

	req := &http.Request{} // can be a real request
	// Act
	val, ok := middleware.GetValidatedBody[dto.CreateUserRequest](req)

	// Assert
	if !ok {
		t.Fatal("expected validated body")
	}
	if val.Name != "Alice" {
		t.Errorf("expected Name 'Alice', got '%s'", val.Name)
	}
	if val.Email != "alice@example.com" {
		t.Errorf("expected Email 'alice@example.com', got '%s'", val.Email)
	}
}
