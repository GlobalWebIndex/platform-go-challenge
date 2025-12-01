package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/models"
)

// Helper function to attach the role to the request context
func attachRole(r *http.Request, cfg *config.Config, role models.Role) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, cfg.ContextKeys.Role, role)
	return r.WithContext(ctx)
}

// Helper function to attach the user ID to the request context
func attachUserID(r *http.Request, cfg *config.Config, userID uuid.UUID) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, cfg.ContextKeys.UserID, userID) // Attach userID to context
	return r.WithContext(ctx)
}
