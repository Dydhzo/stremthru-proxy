package cache

import (
	"sync"
	"time"
)

// CacheConfig holds cache configuration settings
type CacheConfig struct {
	Name     string
	Lifetime time.Duration
}

type cacheItem[T any] struct {
	value   T
	expires time.Time
}

// Cache provides thread-safe generic caching with expiration
type Cache[T any] struct {
	mu       sync.RWMutex
	items    map[string]cacheItem[T]
	lifetime time.Duration
}

// NewCache creates new cache instance with given config
func NewCache[T any](config *CacheConfig) Cache[T] {
	return Cache[T]{
		items:    make(map[string]cacheItem[T]),
		lifetime: config.Lifetime,
	}
}

// Get retrieves item from cache, checking expiration
func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		var zero T
		return zero, false
	}

	if time.Now().After(item.expires) {
		var zero T
		return zero, false
	}

	return item.value, true
}

// Set stores item in cache with TTL expiration
func (c *Cache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem[T]{
		value:   value,
		expires: time.Now().Add(c.lifetime),
	}
}

// Delete removes item from cache
func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}