package cache

import "time"

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
