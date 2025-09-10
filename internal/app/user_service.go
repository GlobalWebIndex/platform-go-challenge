package app

import (
	"context"

	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/ports"
	"github.com/google/uuid"
)

// UserService coordinates user-related operations.
// Depends only on the UserRepository ports.
type UserService struct {
	userRepo ports.UserRepository
}

func NewUserService(users ports.UserRepository) *UserService {
	return &UserService{userRepo: users}
}

// Get returns a user or domain.ErrUserNotFound.
func (userService *UserService) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return userService.userRepo.Get(ctx, id)
}

// Exists reports whether a user exists.
func (userService *UserService) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return userService.userRepo.Exists(ctx, id)
}
