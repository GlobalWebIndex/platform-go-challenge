package cache

import (
	"log"

	lru "github.com/hashicorp/golang-lru/v2"
)

func InitLRUCache[K comparable, V any](size int) *lru.Cache[K, V] {
	c, err := lru.New[K, V](size)
	if err != nil {
		log.Fatalf("failed to initialize LRU cache: %v", err)
	}
	return c
}

func InitLRUCacheWithEvict[K comparable, V any](size int) *lru.Cache[K, V] {
	c, err := lru.NewWithEvict(size, func(key K, value V) {
		log.Printf("Evicted key: %v, value: %v\n", key, value)
	})
	if err != nil {
		log.Fatalf("failed to initialize LRU cache: %v", err)
	}
	return c
}
