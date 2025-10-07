package utils

import (
	"errors"
	"fmt"
	"gwi/configuration"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// ExtractUserID extracts the user ID from the given authorization header.
// It expects the header to contain a Bearer token. If the token is valid and
// contains a user_id claim, it returns the user ID as an integer. If the token
// is invalid or the user_id claim is not found, it returns an error.
//
// Parameters:
//   - authorizationHeader: A string containing the authorization header
//
// Returns:
//   - int: The extracted user ID
//   - error: An error if the token is invalid or user_id is not found
func ExtractUserID(authorizationHeader string) (int, error) {
	// Remove "Bearer " prefix
	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Optional: check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return configuration.JwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user_id
		if uid, ok := claims["user_id"].(float64); ok {
			return int(uid), nil
		}
		return 0, errors.New("user_id not found in token")
	}

	return 0, errors.New("invalid token")
}

// IsUserAuthorized checks if the provided user ID matches the user ID extracted from the
// given authorization header. It returns true if the user is authorized, otherwise false.
//
// Parameters:
//   - userId: The ID of the user to check authorization for.
//   - authorizationHeader: The JWT authorization header from which the user ID is extracted.
//
// Returns:
//   - bool: true if the user is authorized, false otherwise.
func IsUserAuthorized(userId int, authorizationHeader string) bool {
	// fetch user_id from JWT
	jwtUser, _ := ExtractUserID(authorizationHeader)

	return userId == jwtUser
}
