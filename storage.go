package main

import (
	"sync"

	"github.com/google/uuid"
)

type Storage struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*User
}

var store = Storage{
	users: make(map[uuid.UUID]*User),
}

func (s *Storage) GetUser(id uuid.UUID) *User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[id]
}

func (s *Storage) AddUser(u *User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.ID] = u
}
