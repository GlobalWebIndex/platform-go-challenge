package cache

import (
	"challenge/models"
	"errors"
	"runtime"
	"time"
)

type memCacheHandler struct {
	Block              map[string]models.CacheEntry // UserID: CacheEntry{}
	DefaultCacheExpiry time.Duration
	Released           bool
}

// createMemoryHandler initializes the memory structure for the users/assets.
func createMemCacheHander() *memCacheHandler {
	newMemCache := memCacheHandler{}
	newMemCache.Block = make(map[string]models.CacheEntry)
	// Setting default cache expiry at 20 minutes
	newMemCache.DefaultCacheExpiry = 20 * time.Minute
	return &newMemCache
}

// cacheEntryInvalidator checks the cache entry expiry and invalidates it when
// it expires.
func (m *memCacheHandler) cacheEntryInvalidator(entry *models.CacheEntry) {
	for {
		if m.Released {
			return
		}
		if entry.ExpiresAt.Before(time.Now()) {
			entry.Set = false
			Cache.Evict(entry.Token)
			return
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

// Store create a new cache entry mapping the response to the user token
// This entry is replaced on favourites list changes.
func (m *memCacheHandler) Store(token string, response []byte, invalidator bool, expiration ...time.Time) error {
	if m.Released {
		err := errors.New("cache resource is not allocated")
		return err
	}
	newCacheEntry := models.CacheEntry{
		Token:     token,
		Response:  response,
		Set:       true,
		ExpiresAt: time.Now().Add(m.DefaultCacheExpiry),
	}
	if len(expiration) > 0 {
		newCacheEntry.ExpiresAt = expiration[0]
	}
	// Start the invalidator counter (if requested)
	if invalidator {
		go m.cacheEntryInvalidator(&newCacheEntry)
	}
	m.Block[token] = newCacheEntry
	return nil
}

// Load retrieves the cached entry's response body.
func (m *memCacheHandler) Load(token string) (bool, []byte, error) {
	if m.Released {
		err := errors.New("cache resource is not allocated")
		return false, nil, err
	}
	entry := m.Block[token]
	if !entry.Set {
		return false, nil, nil
	}
	if time.Now().Before(entry.ExpiresAt) {
		return true, entry.Response, nil
	} else {
		return false, nil, nil
	}
}

// Evict removes a cache entry.
func (m *memCacheHandler) Evict(token string) error {
	if m.Released {
		err := errors.New("cache resource is not allocated")
		return err
	}
	delete(m.Block, token)
	return nil
}

// Release
func (m *memCacheHandler) Release() error {
	if m.Released {
		err := errors.New("cache resource already released")
		return err
	}
	for k := range m.Block {
		delete(m.Block, k)
	}
	m.Released = true
	runtime.GC()
	return nil
}

// MockCacheEntryData generate mock cache entries data.
func MockCacheEntryData() models.CacheEntry {
	newEntry := models.CacheEntry{
		Token:     "exampleToken",
		Response:  []byte("exampleResponse"),
		Set:       true,
		ExpiresAt: time.Now().Add(5 * time.Second),
	}
	return newEntry
}
