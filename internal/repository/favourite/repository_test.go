package favourite

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetNonexistingUser(t *testing.T) {
	// if user is not in the data set, empty map[uuid]string should be returned

	// GIVEN an empty repo
	r := NewRepository()

	// WHEN I retrieve a user's favourites
	f := r.GetFavourites(uuid.New(), 0, 100)

	// THEN I get an empty map
	assert.NotNil(t, f, "GetFavourites must return empty result for non-existent user uuid")
}

func TestAddAndRetrieveCorrectFavourite(t *testing.T) {
	// if I add a favourite for two users, the favourite for the correct user should be returned
	user1Id := uuid.New()
	favourite1Id := uuid.New()

	user2Id := uuid.New()
	favourite2Id := uuid.New()

	// GIVEN a repo with two users' favourites in it
	r := NewRepository()
	r.AddFavourite(user1Id, "chart1", favourite1Id, "chart")
	r.AddFavourite(user2Id, "insight1", favourite2Id, "insight")

	// WHEN I retrieve each user's favourites
	f1 := r.GetFavourites(user1Id, 0, 100)
	f2 := r.GetFavourites(user2Id, 0, 100)

	// THEN I get the correct favourites for each user
	assert.Equal(t, f1, FavouriteEntries{
		{
			Description:  "chart1",
			FavouriteId:  favourite1Id,
			ResourceType: "chart",
		},
	},
	)

	assert.Equal(t, f2, FavouriteEntries{
		{
			Description:  "insight1",
			FavouriteId:  favourite2Id,
			ResourceType: "insight",
		},
	},
	)
}

func TestPaging(t *testing.T) {
	// Add 25 favourites for one user and retrieve pages of ten.
	// Repo should return first ten, next ten, last five favourites

	// GIVEN a repo with user's 25 favourites
	r := NewRepository()

	userId := uuid.New()
	favourites := make([]uuid.UUID, 25) // used for assertions later
	for i := range 25 {
		favouriteId := uuid.New()
		resourceType := []string{"audience", "chart", "insight"}[rand.IntN(3)]
		favourites[i] = favouriteId
		r.AddFavourite(userId, fmt.Sprintf("my favourite item #%d", i), favouriteId, resourceType)
	}

	// WHEN I retrieve three pages of ten results
	page1 := r.GetFavourites(userId, 0, 10)
	page2 := r.GetFavourites(userId, 10, 10)
	page3 := r.GetFavourites(userId, 20, 5)

	// THEN I get the first ten, second ten, last five results
	page1keys := make([]uuid.UUID, 0, len(page1))
	for _, f := range page1 {
		page1keys = append(page1keys, f.FavouriteId)
	}
	assert.Equal(t, favourites[0:10], page1keys)

	page2keys := make([]uuid.UUID, 0, len(page2))
	for _, f := range page2 {
		page2keys = append(page2keys, f.FavouriteId)
	}
	assert.Equal(t, favourites[10:20], page2keys)

	page3keys := make([]uuid.UUID, 0, len(page3))
	for _, f := range page3 {
		page3keys = append(page3keys, f.FavouriteId)
	}
	assert.Equal(t, favourites[20:25], page3keys)
}

func TestPagingStartOutOfBounds(t *testing.T) {
	// GIVEN a repo with user's 5 favourites
	r := NewRepository()

	userId := uuid.New()
	for _ = range 5 {
		favouriteId := uuid.New()
		resourceType := []string{"audience", "chart", "insight"}[rand.IntN(3)]
		r.AddFavourite(userId, "Description", favouriteId, resourceType)
	}

	// WHEN I retrieve item n+1
	f := r.GetFavourites(userId, 6, 1)

	// I get an empty list of favourites
	assert.Equal(t, 0, len(f))
}

func TestDeleteFavourite(t *testing.T) {
	// if I add three favourites and delete the middle one,
	// a read should return the first and last ones
	userId := uuid.New()
	favourite1Id := uuid.New()
	favourite2Id := uuid.New()
	favourite3Id := uuid.New()

	// GIVEN a repo with two users' favourites in it
	r := NewRepository()
	r.AddFavourite(userId, "audience1", favourite1Id, "audience")
	r.AddFavourite(userId, "chart1", favourite2Id, "chart")
	r.AddFavourite(userId, "insight1", favourite3Id, "insight")

	// WHEN I delete the second favourite
	r.DeleteFavourite(userId, favourite2Id)

	// THEN I retrieve the first and third favourites
	favourites := r.GetFavourites(userId, 0, 10)

	assert.Equal(t, favourites, FavouriteEntries{
		{
			Description:  "audience1",
			FavouriteId:  favourite1Id,
			ResourceType: "audience",
		},
		{
			Description:  "insight1",
			FavouriteId:  favourite3Id,
			ResourceType: "insight",
		},
	},
	)
}
