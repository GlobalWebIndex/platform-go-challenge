package utils

import "fmt"

// Base APIError type
type APIError struct {
	Status  int
	Message string
}

// Implement error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

var (
	ErrBadRequest      = func(msg string) *APIError { return &APIError{Status: 400, Message: msg} }
	ErrNotFound        = func(msg string) *APIError { return &APIError{Status: 404, Message: msg} }
	ErrUnauthorized    = func(msg string) *APIError { return &APIError{Status: 401, Message: msg} }
	ErrUnauthenticated = func(msg string) *APIError { return &APIError{Status: 403, Message: msg} }
	ErrInternalServer  = func(msg string) *APIError { return &APIError{Status: 500, Message: msg} }
)
