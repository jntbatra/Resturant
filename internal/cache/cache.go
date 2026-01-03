package cache

import (
	"sync"
	"time"
)

// CacheEntry represents a cached value with expiration
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired checks if the cache entry has expired
func (ce *CacheEntry) IsExpired() bool {
	return time.Now().After(ce.ExpiresAt)
}

// Cache provides a simple in-memory caching with TTL
// TODO: IMPLEMENTATION - Integrate this cache into services for:
//   - Menu items (long TTL: 30-60 min) - rarely change
//   - Active sessions (medium TTL: 2-10 min) - change moderately
//   - Paginated results (short TTL: 1-5 min) - frequently accessed
//
// Usage: Initialize cache in main.go and pass to service constructors
type Cache struct {
	data    map[string]*CacheEntry
	mu      sync.RWMutex
	ttl     time.Duration
	cleanup time.Duration
	ticker  *time.Ticker
	stop    chan struct{}
}

// NewCache creates a new cache with the specified TTL
func NewCache(ttl time.Duration, cleanupInterval time.Duration) *Cache {
	if ttl == 0 {
		ttl = 5 * time.Minute
	}
	if cleanupInterval == 0 {
		cleanupInterval = 1 * time.Minute
	}

	c := &Cache{
		data:    make(map[string]*CacheEntry),
		ttl:     ttl,
		cleanup: cleanupInterval,
		ticker:  time.NewTicker(cleanupInterval),
		stop:    make(chan struct{}),
	}

	// Start cleanup goroutine
	go c.cleanupExpired()

	return c
}

// Set stores a value in cache with default TTL
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// SetWithTTL stores a value in cache with custom TTL
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a value from cache if it exists and hasn't expired
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if entry.IsExpired() {
		return nil, false
	}

	return entry.Value, true
}

// Delete removes a value from cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

// Clear removes all values from cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*CacheEntry)
}

// cleanupExpired periodically removes expired entries
func (c *Cache) cleanupExpired() {
	for range c.ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.data {
			if now.After(entry.ExpiresAt) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

// Stop stops the cache cleanup goroutine
func (c *Cache) Stop() {
	c.ticker.Stop()
	close(c.stop)
}

// Size returns the number of entries in cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}

// CacheKey generates a cache key for database queries
func CacheKey(prefix string, id string) string {
	return prefix + ":" + id
}

// MenuItemCacheKey generates cache key for menu items
func MenuItemCacheKey(id string) string {
	return CacheKey("menu_item", id)
}

// CategoryCacheKey generates cache key for categories
func CategoryCacheKey(name string) string {
	return CacheKey("category", name)
}

// SessionCacheKey generates cache key for sessions
func SessionCacheKey(id string) string {
	return CacheKey("session", id)
}

// OrderCacheKey generates cache key for orders
func OrderCacheKey(id string) string {
	return CacheKey("order", id)
}

// ListCacheKey generates cache key for list queries
func ListCacheKey(resource string, filters string) string {
	if filters == "" {
		return resource + ":list:all"
	}
	return resource + ":list:" + filters
}
