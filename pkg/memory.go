package gwi

import (
	"context"
	"errors"
	"log"
	"sync"
)

// Store is a map using the userid as main key. It has a child structure formed by maps of strings (fav type) and Asseter arrays
// ["userid"] =
// 							["chart"] = [{Asseter}, {Asseter}, ...]
// 							["audience"] = [{Asseter}, {Asseter}, ...]
// 							["insight"] = [{Asseter}, {Asseter}, ...]
type Store map[string]map[AssetType][]Asseter

// Index is used to store the relation of an assetid and the userid to allow to check for existance of an asset in a quick way
type Index map[string]bool

// MemoryRepository is a struct to handle store, index and mux properties
// db is a map structure to store the favourites asociated to an user by its type
// index is helper to allow to check if an asset is already present in the user
// lock and lockIndex are used to block/unblock db store to allow concurrency
type MemoryRepository struct {
	db        Store
	index     Index
	lock      *sync.RWMutex
	lockIndex *sync.RWMutex
}

// NewMemoryRepository initialize the memoryrepository struct
func NewMemoryRepository() (*MemoryRepository, error) {
	return &MemoryRepository{
		db:        make(Store),
		index:     Index{},
		lock:      &sync.RWMutex{},
		lockIndex: &sync.RWMutex{},
	}, nil
}

// ExistAssetInFavs block/unblock db store when checking if an asset already exists asociated to an user
func (mer MemoryRepository) ExistAssetInFavs(ctx context.Context, assetid string) bool {
	log.Println("Checking if asset already exists")

	mer.lockIndex.RLock()
	defer mer.lockIndex.RUnlock()

	return mer.existAssetInFavs(assetid)
}

// existAssetInFavs checks in index structure if an asset id is associated to an user
func (mer MemoryRepository) existAssetInFavs(assetid string) bool {
	a, ok := mer.index[assetid]
	if !ok {
		return false
	}
	return a
}

// RetrieveFavs block/unblock db store when retrieving favs form an user
func (mer MemoryRepository) RetrieveFavs(ctx context.Context, userid string) (map[AssetType][]Asseter, error) {
	log.Println("Retrieving fav list for user")

	mer.lock.RLock()
	defer mer.lock.RUnlock()

	return mer.retrieveFavs(userid)
}

// retrieveFavs returns fav list of an user
func (mer MemoryRepository) retrieveFavs(userid string) (map[AssetType][]Asseter, error) {
	a, ok := mer.db[userid]
	if !ok {
		return make(map[AssetType][]Asseter), errors.New("no results")
	}
	return a, nil
}

// AddAssetToFavs block/unblock db store when adding a fav to an user
func (mer MemoryRepository) AddAssetToFavs(ctx context.Context, userid string, asset *Asset) (bool, error) {
	log.Println("Adding asset to fav list")

	mer.lock.Lock()
	defer mer.lock.Unlock()

	return mer.addAssetToFavs(ctx, userid, asset)
}

// addAssetToFavs adds a fav if not previously exists
func (mer MemoryRepository) addAssetToFavs(ctx context.Context, userid string, asset *Asset) (bool, error) {
	if mer.db[userid] == nil {
		mer.db[userid] = make(map[AssetType][]Asseter)
		mer.db[userid][TypeChart] = []Asseter{}
		mer.db[userid][TypeInsight] = []Asseter{}
		mer.db[userid][TypeAudience] = []Asseter{}
	}
	data := mer.db[userid]

	switch asset.AssetType {
	case "chart":

		ok := mer.ExistAssetInFavs(ctx, asset.Chart.Id)
		if ok {
			return false, nil
		}

		data[asset.AssetType] = append(data[asset.AssetType], &asset.Chart)
		mer.index[asset.Chart.Id] = true
	case "insight":
		ok := mer.existAssetInFavs(asset.Insight.Id)
		if ok {
			return false, nil
		}
		data[asset.AssetType] = append(data[asset.AssetType], &asset.Insight)
		mer.index[asset.Insight.Id] = true
	case "audience":
		ok := mer.existAssetInFavs(asset.Audience.Id)
		if ok {
			return false, nil
		}
		data[asset.AssetType] = append(data[asset.AssetType], &asset.Audience)
		mer.index[asset.Audience.Id] = true
	}

	mer.db[userid] = data
	return true, nil
}

// AddAssetToFavs block/unblock db store when updating a fav from an user
func (mer MemoryRepository) UpdateFav(ctx context.Context, userid string, asset *Asset) (bool, error) {
	log.Println("Updating asset in fav list")

	mer.lock.Lock()
	defer mer.lock.Unlock()

	return mer.updateFav(ctx, userid, asset)
}

// updateFav updates a fav
func (mer MemoryRepository) updateFav(ctx context.Context, userid string, asset *Asset) (bool, error) {
	assetid := ""
	switch asset.AssetType {
	case "chart":
		assetid = asset.Chart.Id
	case "insight":
		assetid = asset.Insight.Id
	case "audience":
		assetid = asset.Audience.Id
	}

	ok := mer.ExistAssetInFavs(ctx, assetid)

	if ok {
		pos := mer.getIndexPosition(userid, assetid)
		if pos >= 0 {
			data := mer.db[userid]
			switch asset.AssetType {
			case "chart":
				data[asset.AssetType][pos] = &asset.Chart
			case "insight":
				data[asset.AssetType][pos] = &asset.Insight
			case "audience":
				data[asset.AssetType][pos] = &asset.Audience
			}
			return true, nil
		}
	}
	return false, errors.New("asset does not exist")
}

// RemoveFav block/unblock db store when deleting a fav from an user
func (mer MemoryRepository) RemoveFav(ctx context.Context, userid string, asset *Asset) (bool, error) {
	log.Println("Removing asset from fav list")

	mer.lock.Lock()
	defer mer.lock.Unlock()

	return mer.removeFav(ctx, userid, asset)
}

// removeFav deletes a fav from its type array if exists
func (mer MemoryRepository) removeFav(ctx context.Context, userid string, asset *Asset) (bool, error) {
	assetid := ""
	switch asset.AssetType {
	case "chart":
		assetid = asset.Chart.Id
	case "insight":
		assetid = asset.Insight.Id
	case "audience":
		assetid = asset.Audience.Id
	}

	ok := mer.ExistAssetInFavs(ctx, assetid)

	if ok {
		pos := mer.getIndexPosition(userid, assetid)
		if pos >= 0 {
			data := mer.db[userid]
			aux := append(data[asset.AssetType][:pos], data[asset.AssetType][pos+1:]...)
			mer.db[userid][asset.AssetType] = aux
			return true, nil
		}
	}
	return false, errors.New("not found")
}

// getIndexPosition returns the position of an asset if exists
func (mer MemoryRepository) getIndexPosition(userid string, assetid string) int {
	for _, asset := range mer.db[userid] {
		for j, assettype := range asset {
			if assettype.GetId() == assetid {
				return j
			}
		}

	}
	return -1
}
