package repo

import (
	"errors"
	"sync"
	"sort"

	"github.com/KostasDasios/platform-go-challenge/internal/models"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	List(userID string) ([]*models.Favourite, error)
	Create(userID string, fav *models.Favourite) error
	Get(userID, favID string) (*models.Favourite, error)
	UpdateDescription(userID, favID, desc string) (*models.Favourite, error)
	Delete(userID, favID string) error
}

// InMemoryRepo is a thread-safe in-memory implementation intended for the assignment and unit tests.
// It is guarded by an RWMutex; production deployments would use an external store.
type InMemoryRepo struct {
	mu   sync.RWMutex
	data map[string]map[string]*models.Favourite // userID -> favID -> Favourite
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{data: make(map[string]map[string]*models.Favourite)}
}

// List returns all favourites of a given user in deterministic order.
// Results are sorted by creation time (newest first).
func (r *InMemoryRepo) List(userID string) ([]*models.Favourite, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m := r.data[userID]
	out := make([]*models.Favourite, 0, len(m))
	for _, f := range m {
		out = append(out, f)
	}
	// Sort newest first for deterministic output
    sort.Slice(out, func(i, j int) bool {
        return out[i].CreatedAt.After(out[j].CreatedAt)
    })
	return out, nil
}

func (r *InMemoryRepo) Create(userID string, fav *models.Favourite) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.data[userID] == nil {
		r.data[userID] = make(map[string]*models.Favourite)
	}
	r.data[userID][fav.ID] = fav
	return nil
}

func (r *InMemoryRepo) Get(userID, favID string) (*models.Favourite, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m := r.data[userID]
	if m == nil {
		return nil, ErrNotFound
	}
	f, ok := m[favID]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

func (r *InMemoryRepo) UpdateDescription(userID, favID, desc string) (*models.Favourite, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	m := r.data[userID]
	if m == nil {
		return nil, ErrNotFound
	}
	f, ok := m[favID]
	if !ok {
		return nil, ErrNotFound
	}
	f.Description = desc
	return f, nil
}

func (r *InMemoryRepo) Delete(userID, favID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m := r.data[userID]
	if m == nil {
		return ErrNotFound
	}
	if _, ok := m[favID]; !ok {
		return ErrNotFound
	}
	delete(m, favID)
	return nil
}
