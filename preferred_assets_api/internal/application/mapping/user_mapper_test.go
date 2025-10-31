package mapping_test

import (
	"strings"
	"testing"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/mapping"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

// --- Test UserReqToDomain ---
func TestUserReqToDomain_Success(t *testing.T) {
	// Arrange
	req := dto.CreateUserRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	}

	// Act
	user, err := mapping.UserReqToDomain(req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != req.Name {
		t.Errorf("expected Name '%s', got '%s'", req.Name, user.Name)
	}
	if user.Email != req.Email {
		t.Errorf("expected Email '%s', got '%s'", req.Email, user.Email)
	}
	if user.Password != req.Password {
		t.Errorf("expected Password '%s', got '%s'", req.Password, user.Password)
	}
	if strings.TrimSpace(user.Id) == "" {
		t.Error("expected Id to be set")
	}
}

// --- Unhappy path for UserReqToDomain: missing fields ---
func TestUserReqToDomain_EmptyRequest(t *testing.T) {
	// Arrange
	req := dto.CreateUserRequest{}

	// Act
	user, err := mapping.UserReqToDomain(req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err) // still returns no error in current implementation
	}
	if user.Name != "" || user.Email != "" || user.Password != "" {
		t.Error("expected fields to be empty")
	}
}

// --- Test UpdateReqToDomain ---
func TestUpdateReqToDomain_Success(t *testing.T) {
	// Arrange
	existing := &domain.User{Name: "Alice", Email: "alice@example.com", Password: "oldpass"}
	req := dto.UpdateUserRequest{Name: "Bob", Password: "newpass123"}

	// Act
	updated := mapping.UpdateReqToDomain(existing, req)

	// Assert
	if updated.Name != "Bob" {
		t.Errorf("expected Name 'Bob', got '%s'", updated.Name)
	}
	if updated.Email != "alice@example.com" {
		t.Errorf("expected Email unchanged, got '%s'", updated.Email)
	}
	if updated.Password != "newpass123" {
		t.Errorf("expected Password 'newpass123', got '%s'", updated.Password)
	}
}

// --- Unhappy path for UpdateReqToDomain: empty request ---
func TestUpdateReqToDomain_EmptyRequest(t *testing.T) {
	// Arrange
	existing := &domain.User{Name: "Alice", Email: "alice@example.com", Password: "oldpass"}
	req := dto.UpdateUserRequest{}

	// Act
	updated := mapping.UpdateReqToDomain(existing, req)

	// Assert
	if updated.Name != "Alice" || updated.Email != "alice@example.com" || updated.Password != "oldpass" {
		t.Error("expected existing user to remain unchanged")
	}
}

// --- Test DomainToUserRes ---
func TestDomainToUserRes(t *testing.T) {
	// Arrange
	user := domain.User{Id: "123", Name: "Alice", Email: "alice@example.com"}

	// Act
	res := mapping.DomainToUserRes(user)

	// Assert
	if res.ID != "123" || res.Name != "Alice" || res.Email != "alice@example.com" {
		t.Errorf("unexpected UserResponse: %+v", res)
	}
}

// --- Test UserReqToResponseList ---
func TestUserReqToResponseList(t *testing.T) {
	// Arrange
	users := []domain.User{
		{Id: "1", Name: "Alice", Email: "a@example.com"},
		{Id: "2", Name: "Bob", Email: "b@example.com"},
	}

	// Act
	responses := mapping.UserReqToResponseList(users)

	// Assert
	if len(responses) != 2 {
		t.Errorf("expected 2 responses, got %d", len(responses))
	}
	if responses[0].ID != "1" || responses[1].ID != "2" {
		t.Errorf("unexpected IDs in responses: %+v", responses)
	}
}
