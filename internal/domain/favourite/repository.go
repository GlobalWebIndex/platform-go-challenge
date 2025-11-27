package favourite

import (
	"github.com/google/uuid"
)

// Repository Interface for favorites
type Repository interface {
	GetByID(userID uuid.UUID, favoriteID uuid.UUID) (*Favorite, error)
	GetAll(userID uuid.UUID) ([]Favorite, error)
	Add(userID uuid.UUID, favorite Favorite) error
	Update(userID uuid.UUID, favorite Favorite) error
	Delete(userID uuid.UUID, favoriteID uuid.UUID) error
}
