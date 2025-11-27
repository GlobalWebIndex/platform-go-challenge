package user

import "github.com/google/uuid"

type Repository interface {
	GetByID(id uuid.UUID) (*User, error)
	GetByUsername(username string) (*User, error)

	Add(username, plainPassword string, roles []string) (*User, error)
}
