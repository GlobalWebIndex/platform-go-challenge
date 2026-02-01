package store

import (
	"errors"
	"sort"
	"sync"

	"github.com/kongrillas/platform-go-challenge/internal/model"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]map[string]model.Asset
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]map[string]model.Asset),
	}
}

func (s *MemoryStore) List(userID string) ([]model.Asset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userMap, ok := s.data[userID]
	if !ok {
		return []model.Asset{}, nil
	}

	out := make([]model.Asset, 0, len(userMap))
	for _, a := range userMap {
		out = append(out, a)
	}

	// deterministic output (nice for tests)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (s *MemoryStore) Add(userID string, a model.Asset) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[userID]; !ok {
		s.data[userID] = make(map[string]model.Asset)
	}
	if _, exists := s.data[userID][a.ID]; exists {
		return ErrAlreadyExists
	}
	s.data[userID][a.ID] = a
	return nil
}

func (s *MemoryStore) Remove(userID, assetID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	userMap, ok := s.data[userID]
	if !ok {
		return ErrNotFound
	}
	if _, ok := userMap[assetID]; !ok {
		return ErrNotFound
	}
	delete(userMap, assetID)
	return nil
}

func (s *MemoryStore) UpdateDescription(userID, assetID, description string) (model.Asset, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userMap, ok := s.data[userID]
	if !ok {
		return model.Asset{}, ErrNotFound
	}
	a, ok := userMap[assetID]
	if !ok {
		return model.Asset{}, ErrNotFound
	}
	a.Description = description
	userMap[assetID] = a
	return a, nil
}
