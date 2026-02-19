package store

import (
	"errors"
	"time"
)

var (
	// ErrNotFound indicates the requested resource was not found
	ErrNotFound = errors.New("resource not found")

	// ErrConflict indicates the resource already exists
	ErrConflict = errors.New("resource already exists")

	// QueryTimeoutDuration is the default timeout for database queries
	QueryTimeoutDuration = time.Second * 5
)

// Storage groups all stores together providing a single access point
// for all data operations. This makes it easy to inject into services
// and ensures consistent store initialization.
type Storage struct {
	Favourites FavouriteStore
}

// NewStorage creates a new Storage instance with all stores initialized
// Currently uses in-memory implementations, but can easily be swapped
// for database implementations without changing the interface.
func NewStorage() *Storage {
	return &Storage{
		Favourites: NewInMemoryFavouriteStore(),
	}
}
