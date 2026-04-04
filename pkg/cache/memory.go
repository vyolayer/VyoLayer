package cache

// type entry[V any] struct {
// 	value     V
// 	expiresAt time.Time // zero value means no expiry
// }

// func (e *entry[V]) isExpired() bool {
// 	return !e.expiresAt.IsZero() && time.Now().After(e.expiresAt)
// }

// // MemoryCache is a thread-safe, generic in-memory cache with TTL support.
// type MemoryCache[V any] struct {
// 	mu      sync.RWMutex
// 	entries map[string]*entry[V]
// }

// // NewMemoryCache creates a new in-memory cache.
// // Optionally start a background cleanup goroutine by calling StartCleanup.
// func NewMemoryCache[V any]() *MemoryCache[V] {
// 	return &MemoryCache[V]{
// 		entries: make(map[string]*entry[V]),
// 	}
// }

// // Get retrieves a value by key. Returns zero value and false if not found or expired.
// func (c *MemoryCache[V]) Get(key string) (V, bool) {
// 	c.mu.RLock()
// 	e, ok := c.entries[key]
// 	c.mu.RUnlock()

// 	if !ok || e.isExpired() {
// 		var zero V
// 		return zero, false
// 	}

// 	return e.value, true
// }

// // Set stores a value with the given TTL. A TTL of 0 means no expiry.
// func (c *MemoryCache[V]) Set(key string, value V, ttl time.Duration) {
// 	var expiresAt time.Time
// 	if ttl > 0 {
// 		expiresAt = time.Now().Add(ttl)
// 	}

// 	c.mu.Lock()
// 	c.entries[key] = &entry[V]{value: value, expiresAt: expiresAt}
// 	c.mu.Unlock()
// }

// // Delete removes a key from the cache.
// func (c *MemoryCache[V]) Delete(key string) {
// 	c.mu.Lock()
// 	delete(c.entries, key)
// 	c.mu.Unlock()
// }

// // Flush removes all entries from the cache.
// func (c *MemoryCache[V]) Flush() {
// 	c.mu.Lock()
// 	c.entries = make(map[string]*entry[V])
// 	c.mu.Unlock()
// }

// // StartCleanup launches a background goroutine that periodically removes
// // expired entries. The goroutine stops when the provided stop channel is closed.
// func (c *MemoryCache[V]) StartCleanup(interval time.Duration, stop <-chan struct{}) {
// 	go func() {
// 		ticker := time.NewTicker(interval)
// 		defer ticker.Stop()
// 		for {
// 			select {
// 			case <-ticker.C:
// 				c.deleteExpired()
// 			case <-stop:
// 				return
// 			}
// 		}
// 	}()
// }

// func (c *MemoryCache[V]) deleteExpired() {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	for key, e := range c.entries {
// 		if e.isExpired() {
// 			delete(c.entries, key)
// 		}
// 	}
// }
