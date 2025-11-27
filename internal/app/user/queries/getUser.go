package queries

import (
	"fmt"
	"github.com/akazantzidis/gwi-ass/internal/domain/user"
	"github.com/google/uuid"
)

// GetUserRequest represents the query input
type GetUserRequest struct {
	ID       uuid.UUID
	Username string // optional: either ID or Username
}

// GetUserResult represents the returned user
type GetUserResult struct {
	ID       string
	Username string
	Roles    []string
}

// GetUserHandler interface
type GetUserHandler interface {
	Handle(req GetUserRequest) (*GetUserResult, error)
}

type getUserHandler struct {
	repo user.Repository
}

func NewGetUserHandler(userRepo user.Repository) GetUserHandler {
	return &getUserHandler{repo: userRepo}
}

// Handle fetches a user by ID or username
func (h *getUserHandler) Handle(req GetUserRequest) (*GetUserResult, error) {
	var u *user.User
	var err error

	if req.ID.String() != "" {
		u, err = h.repo.GetByID(req.ID)
	} else if req.Username != "" {
		u, err = h.repo.GetByUsername(req.Username)
	} else {
		return nil, fmt.Errorf("no identifier provided")
	}

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &GetUserResult{
		ID:       u.ID.String(),
		Username: u.Username,
		Roles:    u.Roles,
	}, nil
}
