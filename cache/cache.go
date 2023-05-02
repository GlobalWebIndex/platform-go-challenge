package cache

import (
	"errors"
	"time"
)

// Cache is the cache access controller
var Cache CacheHandler

// CacheHandler interface is impemented using all the required methods to
// interface with the various cache options
type CacheHandler interface {
	// Store create a new cache entry mapping the response to the user token
	// This entry is replaced on favourites list changes.
	Store(token string, response []byte, invalidator bool, expiration ...time.Time) error
	// Load retrieves the cached entry's response body.
	Load(token string) (bool, []byte, error)
	// Evict removes a cache entry.
	Evict(token string) error
	// Release frees the resources allocated for the cache.
	Release() error
}

// ConnectoToCache handles the initialization of the cache option
func ConnectToCache(option string) error {
	var err error
	cacheOption := option
	switch cacheOption {
	case "memory":
		Cache = createMemCacheHander()
	default:
		err = errors.New("unsupported cache option")
	}
	return err
}
