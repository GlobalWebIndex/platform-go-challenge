package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError represents a structured, actionable error for API consumers.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) GetMessage() string {
	return e.Message
}

func (e *APIError) GetField() string {
	return e.Field
}

func (e *APIError) GetCode() int {
	return e.Code
}

// NewAPIError creates a new structured error.
func NewAPIError(code int, message, field string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Field:   field,
	}
}

// ErrorHandler is a middleware to catch and handle errors.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var apiErr *APIError
			if errors.As(err, &apiErr) {
				c.JSON(apiErr.Code, apiErr)
			} else {
				// Default to a 500 Internal Server Error for unhandled errors
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "An unexpected error occurred",
				})
			}
		}
	}
}
