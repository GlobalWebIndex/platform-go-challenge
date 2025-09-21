package gwiexercise

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type InMemoryStore struct {
	mu    sync.RWMutex
	users map[string]map[string]Asset // userID -> assetID -> Asset
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{users: make(map[string]map[string]Asset)}
}

func (s *InMemoryStore) List(userID string) []Asset {
	s.mu.RLock()
	defer s.mu.RUnlock()
	umap := s.users[userID]
	if umap == nil {
		return []Asset{}
	}
	res := make([]Asset, 0, len(umap))
	for _, a := range umap {
		res = append(res, a)
	}
	return res
}

func (s *InMemoryStore) ensureUser(userID string) map[string]Asset {
	m, ok := s.users[userID]
	if !ok {
		m = make(map[string]Asset)
		s.users[userID] = m
	}
	return m
}

func (s *InMemoryStore) Add(userID string, a Asset) (Asset, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a.GetID() == "" {
		setID(a, newID())
	}
	// set CreatedAt if zero
	setCreatedAtIfZero(a, time.Now())
	umap := s.ensureUser(userID)
	if _, exists := umap[a.GetID()]; exists {
		return nil, errors.New("asset with same id already exists for user")
	}
	umap[a.GetID()] = a
	return a, nil
}

func (s *InMemoryStore) Remove(userID, assetID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	umap := s.users[userID]
	if umap == nil {
		return errors.New("user not found")
	}
	if _, ok := umap[assetID]; !ok {
		return errors.New("asset not found")
	}
	delete(umap, assetID)
	return nil
}

func (s *InMemoryStore) UpdateDescription(userID, assetID, description string) (Asset, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	umap := s.users[userID]
	if umap == nil {
		return nil, errors.New("user not found")
	}
	a, ok := umap[assetID]
	if !ok {
		return nil, errors.New("asset not found")
	}
	a.SetDescription(description)
	umap[assetID] = a
	return a, nil
}

// Helpers

func newID() string {
	now := time.Now().UnixNano()
	return fmtID(now, rng.Int63())
}

func fmtID(x, y int64) string {
	return base36(uint64(x)) + base36(uint64(y))
}

func base36(u uint64) string {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	if u == 0 {
		return "0"
	}
	buf := make([]byte, 0, 13)
	for u > 0 {
		r := int(u % 36)
		buf = append([]byte{alphabet[r]}, buf...)
		u /= 36
	}
	return string(buf)
}

// setID sets the ID on the concrete asset via type assertion.
func setID(a Asset, id string) {
	switch v := a.(type) {
	case *ChartAsset:
		v.ID = id
	case *InsightAsset:
		v.ID = id
	case *AudienceAsset:
		v.ID = id
	}
}

func setCreatedAtIfZero(a Asset, t time.Time) {
	switch v := a.(type) {
	case *ChartAsset:
		if v.CreatedAt.IsZero() {
			v.CreatedAt = t
		}
	case *InsightAsset:
		if v.CreatedAt.IsZero() {
			v.CreatedAt = t
		}
	case *AudienceAsset:
		if v.CreatedAt.IsZero() {
			v.CreatedAt = t
		}
	}
}
