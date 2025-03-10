package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultGetAllNumberOfWorkers = 4
	defaultGetAllTimeout         = 30 * time.Second // suitable for local testing
)

// Favorite is a user favorite asset with a custom description.
type Favorite struct {
	ID          int    `json:"id,omitempty"`
	AssetKey    string `json:"assetKey,omitempty"`
	Description string `json:"description,omitempty"`
}

// FavoriteWithAsset is a favorite plus the asset is refers to.
type FavoriteWithAsset struct {
	Favorite
	IncludesAsset bool
	Asset         any
}

// GetAllFavoritesResponse is the response type for getAll.
type GetAllFavoritesResponse struct {
	Total     int                 `json:"total"`
	Favorites []FavoriteWithAsset `json:"favorites"`
}

// FavoritesStorage is an inmemory store for favorites.
// It is a map indexed by users.
type FavoritesStorage struct {
	sync.Mutex
	favorites map[string][]*Favorite
}

// FavoritesService provides CRUD handler for favorites.
type FavoritesService struct {
	storage *FavoritesStorage
	fetcher *AssetsFetcher
}

// normalize trims space of the favorites attributes.
func (f *Favorite) normalize() {
	f.AssetKey = strings.TrimSpace(f.AssetKey)
	f.Description = strings.TrimSpace(f.Description)
}

// newFavoritesStorage returns a new, empty FavoritesStorage
func newFavoritesStorage() *FavoritesStorage {
	return &FavoritesStorage{
		favorites: make(map[string][]*Favorite),
	}
}

// add adds the favorite if not already there and returns the non-zero id.
func (fs *FavoritesStorage) add(owner string, fav Favorite) int {
	fs.Lock()
	defer fs.Unlock()

	var id int

	// check if there is already a favorite for the same asset
	i := slices.IndexFunc(fs.favorites[owner], func(f *Favorite) bool {
		return f != nil && f.AssetKey == fav.AssetKey
	})
	if i == -1 {
		// we don't allow the zero value for ID
		id = len(fs.favorites[owner]) + 1
		fs.favorites[owner] = append(fs.favorites[owner], &Favorite{
			ID:          id,
			AssetKey:    fav.AssetKey,
			Description: fav.Description,
		})
	} else {
		id = fs.favorites[owner][i].ID
	}

	return id
}

// getByID returns the user favorite with ID.
func (fs *FavoritesStorage) getByID(owner string, id int) (Favorite, bool) {
	fs.Lock()
	defer fs.Unlock()

	if i := fs.indexByID(owner, id); i != -1 {
		return *fs.favorites[owner][i], true
	}
	return Favorite{}, false
}

// getAll returns all the user favorites
func (fs *FavoritesStorage) getAll(owner string) []Favorite {
	fs.Lock()
	defer fs.Unlock()

	orig := fs.favorites[owner]
	favs := make([]Favorite, 0, len(orig))
	for _, f := range orig {
		if f != nil {
			favs = append(favs, *f)
		}
	}
	return favs
}

// del removes the favorite. It nils the slice slot so that ids are not reused.
func (fs *FavoritesStorage) del(owner string, id int) {
	fs.Lock()
	defer fs.Unlock()

	if i := fs.indexByID(owner, id); i != -1 {
		fs.favorites[owner][i] = nil
	}
}

// update updates the properties of the favorite
func (fs *FavoritesStorage) update(owner string, id int, fav Favorite) {
	fs.Lock()
	defer fs.Unlock()

	if i := fs.indexByID(owner, id); i != -1 {
		fs.favorites[owner][i].Description = fav.Description
	}
}

// indexByID returns the index of the favorite in the owners slice or -1 if not found
// Called should hold the lock.
func (fs *FavoritesStorage) indexByID(owner string, id int) int {
	// exploit that id is the index of the favorite + 1. See add.
	if i := id - 1; 0 <= i && i < len(fs.favorites[owner]) {
		fav := fs.favorites[owner][i]
		if fav != nil { // not deleted
			return i
		}
	}
	return -1
}

// newFavoritesService returns a new FavoritesService backedup by fs.
func newFavoritesService(fs *FavoritesStorage, af *AssetsFetcher) *FavoritesService {
	return &FavoritesService{fs, af}
}

// add accepts a json Favorite and adds it to the store.
func (fs *FavoritesService) add(w http.ResponseWriter, r *http.Request) {
	var fav Favorite
	if err := json.NewDecoder(r.Body).Decode(&fav); err != nil {
		http.Error(w, "bad input: "+err.Error(), http.StatusBadRequest)
		return
	}
	fav.normalize()

	if fav.AssetKey == "" {
		http.Error(w, "empty asset key", http.StatusBadRequest)
		return
	}

	// TODO: should it call the asset server?

	owner := r.Context().Value(keyUser).(User).Email
	id := fs.storage.add(owner, fav)
	w.Header().Set("Location", fmt.Sprintf("/favorites/%d", id))
	w.WriteHeader(http.StatusCreated)
}

// del removes the favorite with id.
func (fs *FavoritesService) del(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	owner := r.Context().Value(keyUser).(User).Email
	if _, present := fs.storage.getByID(owner, id); present {
		fs.storage.del(owner, id)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "not found", http.StatusNotFound)
	}
}

// update updates the favorite with id.
func (fs *FavoritesService) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var fav Favorite
	if err := json.NewDecoder(r.Body).Decode(&fav); err != nil {
		http.Error(w, "bad input: "+err.Error(), http.StatusBadRequest)
		return
	}
	fav.normalize()

	owner := r.Context().Value(keyUser).(User).Email
	if prev, present := fs.storage.getByID(owner, id); present {
		fs.storage.update(owner, id, fav)
		w.WriteHeader(http.StatusOK)
		fav.ID = prev.ID
		fav.AssetKey = prev.AssetKey
		json.NewEncoder(w).Encode(fav)
	} else {
		http.Error(w, "not found", http.StatusNotFound)
	}
}

// get returns the favorite with id.
func (fs *FavoritesService) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// TODO: should it call the asset server?

	owner := r.Context().Value(keyUser).(User).Email
	if fav, present := fs.storage.getByID(owner, id); present {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fav)
	} else {
		http.Error(w, "not found", http.StatusNotFound)
	}
}

// getAll returns all the favorites of the users including the asset.
// It calls the asset server to fetch them.
// It returns each favorite with the corresponding asset if it was
// fetched successfully, otherwise the asset is null.
func (fs *FavoritesService) getAll(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(keyUser).(User).Email
	favorites := fs.storage.getAll(owner)
	// need an index by key because the fetch pool only works with keys
	favoritesMap := make(map[string]Favorite)
	for _, f := range favorites {
		favoritesMap[f.AssetKey] = f
	}

	ctx, cancel := context.WithTimeout(r.Context(), defaultGetAllTimeout)
	defer cancel()
	snd, rcv := fs.fetcher.newFetcher(ctx, defaultGetAllNumberOfWorkers)

	// start a goroutine to send the keys to the fetch pool
	go func() {
		defer close(snd)
		for _, f := range fs.storage.getAll(owner) {
			select {
			case snd <- f.AssetKey:
				//
			case <-ctx.Done():
				log.Printf("producer stops, context done: %s", ctx.Err())
				return
			}
		}
	}()

	// read results from the fetch pool
	var response GetAllFavoritesResponse
	response.Total = len(favorites)
	for res := range rcv {
		var body []byte
		if res.Err == nil {
			body = res.Body
		} else {
			log.Printf("fetch error for %s: %v", res.Key, res.Err)
			body = nil
		}
		var err error
		if strings.Contains(res.Key, "audience") {
			err = addAsset[Audience](&response.Favorites, favoritesMap[res.Key], body)
		} else if strings.Contains(res.Key, "charts") {
			err = addAsset[Chart](&response.Favorites, favoritesMap[res.Key], body)
		} else if strings.Contains(res.Key, "insights") {
			err = addAsset[Insight](&response.Favorites, favoritesMap[res.Key], body)
		} else {
			err = errors.New("unknown asset key")
		}
		if err != nil {
			log.Printf("encoding error for %s: %v", res.Key, err)
		}
	}

	if ctx.Err() != nil {
		log.Printf("consumer ended prematurely, context done: %s", ctx.Err())
		// TODO: check with favoritesMap for favorites not included in the response
	}

	json.NewEncoder(w).Encode(response)
}

// addAsset is a helper. Add the favorite in the response slice, including the
// data asset if it can be json decoded.
func addAsset[T any](s *[]FavoriteWithAsset, fav Favorite, data []byte) error {
	var f FavoriteWithAsset
	f.Favorite = fav

	var err error
	if data != nil {
		var t T
		if err = json.NewDecoder(bytes.NewBuffer(data)).Decode(&t); err == nil {
			f.IncludesAsset = true
			f.Asset = t
		}
	}

	*s = append(*s, f)
	return err
}
