package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"platform-go-challenge/internal/common"
)

// TestIsValidRole tests the IsValidRole function
func TestIsValidRole(t *testing.T) {
	tests := []struct {
		role     string
		expected bool
	}{
		{"user", true},
		{"admin", true},
		{"guest", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			assert.Equal(t, tt.expected, isValidRole(tt.role))
		})
	}
}

// TestCreateUser tests the CreateUser function
func TestCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		role        Role
		expectError bool
		errorMsg    string
	}{
		{"valid user", "alice", "password123", RoleUser, false, ""},
		{"valid admin", "bob", "adminpass", RoleAdmin, false, ""},
		{"empty username", "", "pass", RoleUser, true, common.ErrEmptyUsername},
		{"empty password", "charlie", "", RoleUser, true, common.ErrEmptyPassword},
		{"invalid role", "dave", "pass", Role("invalid"), true, common.ErrInvalidRole},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := CreateUser(tt.username, tt.password, tt.role)

			if tt.expectError {
				assert.EqualError(t, err, tt.errorMsg)
				return
			}

			// For valid cases
			assert.NoError(t, err)
			assert.Equal(t, tt.username, user.Username)
			assert.Equal(t, tt.role, user.UserRole)
			assert.NotEqual(t, tt.password, user.Password)

			// Verify hashed password matches
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.password))
			assert.NoError(t, err)
		})
	}
}
