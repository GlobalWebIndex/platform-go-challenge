package entadapter

import (
	"context"

	"github.com/SokratisChaimanas/platform-go-challenge/ent"
	"github.com/SokratisChaimanas/platform-go-challenge/ent/user"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/ports"
	"github.com/google/uuid"
)

// Compile time safety for ports.UserRepository implementation
var _ ports.UserRepository = (*UserRepo)(nil)

// UserRepo implements ports.UserRepository using Ent.
type UserRepo struct {
	client *ent.Client
}

func NewUserRepo(client *ent.Client) *UserRepo {
	return &UserRepo{client: client}
}

// Get returns the user with the given id.
func (userRepo *UserRepo) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, err := userRepo.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
	}, nil
}

// Exists checks if a user with this id exists in DB.
func (userRepo *UserRepo) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return userRepo.client.User.
		Query().
		Where(user.ID(id)).
		Exist(ctx)
}
