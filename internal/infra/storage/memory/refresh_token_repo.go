package memory

import (
	"github.com/akazantzidis/gwi-ass/internal/domain/token"
	"sync"
)

type RefreshRepo struct {
	mu     sync.RWMutex
	tokens map[string]token.RefreshRecord
}

func NewRefreshRepo() *RefreshRepo {
	return &RefreshRepo{
		tokens: make(map[string]token.RefreshRecord),
	}
}

func (r *RefreshRepo) Save(token string, rec token.RefreshRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[token] = rec
}

func (r *RefreshRepo) Get(token string) (token.RefreshRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.tokens[token]
	return rec, ok
}

func (r *RefreshRepo) Delete(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tokens, token)
}
