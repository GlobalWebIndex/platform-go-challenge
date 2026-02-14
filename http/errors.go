package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
	ErrCodeBadRequest       = "BAD_REQUEST"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	ErrCodeInternal         = "INTERNAL_ERROR"
	ErrCodeRateLimit        = "RATE_LIMIT_EXCEEDED"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func WriteError(w http.ResponseWriter, statusCode int, code, message string) {
	writeErrorResponse(w, statusCode, code, message)
}

// WriteErrorWithLog logs the full error server-side; sends a generic message to the client (no internal details).
func WriteErrorWithLog(w http.ResponseWriter, r *http.Request, statusCode int, code, publicMessage string, internalErr error) {
	if internalErr != nil {
		if r != nil {
			slog.Error("api error", "method", r.Method, "path", r.URL.Path, "code", code, "status", statusCode, "err", internalErr)
		} else {
			slog.Error("api error", "code", code, "status", statusCode, "err", internalErr)
		}
	}
	writeErrorResponse(w, statusCode, code, publicMessage)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(newErrorResponse(code, message))
}

func newErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{Code: code, Message: message},
	}
}
