package favourite

import (
	"errors"

	"github.com/google/uuid"
)

type FavouriteEntry struct {
	Description  string
	FavouriteId  uuid.UUID
	ResourceType string
}

type FavouriteEntries []FavouriteEntry

type Repository struct {
	// keys in the map are the uniqe ids of the users
	favourites map[uuid.UUID]FavouriteEntries
}

func NewRepository() *Repository {
	return &Repository{
		favourites: make(map[uuid.UUID]FavouriteEntries, 100),
	}
}

func (r *Repository) GetFavourites(userid uuid.UUID, offset int, pageSize int) FavouriteEntries {
	if result, ok := r.favourites[userid]; !ok {
		return make(FavouriteEntries, 0, 100)
	} else {
		if offset > len(result) {
			return FavouriteEntries{}
		}

		high := offset + pageSize
		if high > len(result) {
			high = len(result)
		}
		return result[offset:high]
	}
}

func (r *Repository) AddFavourite(userId uuid.UUID, description string, favouriteId uuid.UUID, favouriteType string) {
	if _, ok := r.favourites[userId]; !ok {
		r.favourites[userId] = make(FavouriteEntries, 0, 100)
	}

	// not checking for duplicate entries
	r.favourites[userId] = append(r.favourites[userId], FavouriteEntry{
		Description:  description,
		FavouriteId:  favouriteId,
		ResourceType: favouriteType,
	})
}

func (r *Repository) DeleteFavourite(userId uuid.UUID, favouriteId uuid.UUID) error {
	if _, ok := r.favourites[userId]; !ok {
		// user doesn't exist. Deletion succeeds silently.
		return nil
	}

	// find element with favouriteId, replace it with last element in slice
	// then return a 1-shorter slice
	for i, favourite := range r.favourites[userId] {
		if favourite.FavouriteId == favouriteId {
			l := len(r.favourites[userId])
			r.favourites[userId][i] = r.favourites[userId][l-1]
			r.favourites[userId] = r.favourites[userId][:l-1]
		}
	}
	// no actual errors in this mock implementation, a real implementation
	// would contact a DB etc. and potentially return an error
	return nil
}

func (r *Repository) UpdateFavourite(userId uuid.UUID, description string, favouriteId uuid.UUID, favouriteType string) error {
	if _, ok := r.favourites[userId]; !ok {
		// user doesn't exist. Deletion succeeds silently.
		return errors.New("user not found")
	}

	// update element with favouriteId
	found := false
	for i, favourite := range r.favourites[userId] {
		if favourite.FavouriteId == favouriteId {
			r.favourites[userId][i].Description = description
			r.favourites[userId][i].ResourceType = favouriteType
			found = true
		}
	}
	if !found {
		return errors.New("favourite not found")
	}

	return nil
}
