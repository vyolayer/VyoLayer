package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// V is the value type stored in the cache.
type Cache[V any] interface {
	// Get retrieves a value by key. Returns the value and true if found and not expired.
	Get(key string) (V, bool)

	// Set stores a value with the given TTL. A TTL of 0 means no expiry.
	Set(key string, value V, ttl time.Duration)

	// Delete removes a key from the cache.
	Delete(key string)

	// Flush removes all entries from the cache.
	Flush()
}

// memoryCache is the local RAM implementation of the Store interface.
type memoryCache struct {
	client *cache.Cache
}

// NewMemoryCache creates a fast, local RAM cache with background cleanup.
func NewMemoryCache(defaultExpiration, cleanupInterval time.Duration) Cache[any] {
	return &memoryCache{
		client: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (m *memoryCache) Get(key string) (any, bool) {
	return m.client.Get(key)
}

func (m *memoryCache) Set(key string, value any, duration time.Duration) {
	m.client.Set(key, value, duration)
}

func (m *memoryCache) Delete(key string) {
	m.client.Delete(key)
}

func (m *memoryCache) Flush() {
	m.client.Flush()
}
