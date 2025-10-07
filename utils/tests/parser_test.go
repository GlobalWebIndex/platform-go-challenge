package utils_test

import (
	"gwi/configuration"
	"gwi/utils"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJwtToken(userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
	})
	tokenString, _ := token.SignedString(configuration.JwtSecret)
	return tokenString
}

func TestExtractUserID_ValidToken(t *testing.T) {
	userID := 1
	authorizationHeader := "Bearer " + GenerateJwtToken(userID)

	extractedID, err := utils.ExtractUserID(authorizationHeader)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if extractedID != userID {
		t.Errorf("expected user_id %d, got %d", userID, extractedID)
	}
}

func TestExtractUserID_InvalidToken(t *testing.T) {
	authorizationHeader := "Bearer invalid-token"

	extractedID, err := utils.ExtractUserID(authorizationHeader)

	// Assert
	if err == nil {
		t.Fatalf("expected an error, got none")
	}
	if extractedID != 0 {
		t.Errorf("expected user_id 0, got %d", extractedID)
	}
}

func TestIsUserAuthorized_ValidAuthorization(t *testing.T) {
	userID := 1
	authorizationHeader := "Bearer " + GenerateJwtToken(userID)

	authorized := utils.IsUserAuthorized(userID, authorizationHeader)

	// Assert
	if !authorized {
		t.Errorf("expected user to be authorized, got %v", authorized)
	}
}

func TestIsUserAuthorized_InvalidAuthorization(t *testing.T) {
	userID := 1
	authorizationHeader := "Bearer invalid-token"

	authorized := utils.IsUserAuthorized(userID, authorizationHeader)

	// Assert
	if authorized {
		t.Errorf("expected user to be unauthorized, got %v", authorized)
	}
}
