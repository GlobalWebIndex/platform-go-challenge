package memory

import (
	"errors"
	"github.com/akazantzidis/gwi-ass/internal/domain/user"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo struct {
	mu    sync.RWMutex
	users map[string]*user.User // key = username
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users: make(map[string]*user.User),
	}
}

func (r *UserRepo) Add(username, plainPassword string, roles []string) (*user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		ID:       uuid.New(),
		Username: username,
		Password: string(hashed),
		Roles:    roles,
	}

	r.users[username] = u
	return u, nil
}

func (r *UserRepo) GetByID(id uuid.UUID) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (r *UserRepo) GetByUsername(username string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[username]
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}
