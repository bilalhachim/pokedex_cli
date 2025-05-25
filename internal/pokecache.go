package internal

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntry map[string]cacheEntry
	mu         *sync.RWMutex // Use RWMutex for better read performance
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// NewCache creates a new Cache instance and starts the cleanup loop
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cacheEntry: make(map[string]cacheEntry),
		mu:         &sync.RWMutex{},
	}
	go cache.reapLoop(interval)
	return cache
}

// Add adds a new entry to the cache
func (cache *Cache) Add(key string, val []byte) {
	cacheEntry := cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.cacheEntry[key] = cacheEntry
}

// Get retrieves an entry from the cache
func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.RLock() // Use read lock for better concurrency
	defer cache.mu.RUnlock()
	entry, ok := cache.cacheEntry[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

// reapLoop periodically removes expired entries from the cache
func (cache *Cache) reapLoop(interval time.Duration) {
	for {
		time.Sleep(interval) // Wait before starting the cleanup
		cache.mu.Lock()
		for key, entry := range cache.cacheEntry {
			if time.Since(entry.createdAt) >= interval {
				delete(cache.cacheEntry, key)
			}
		}
		cache.mu.Unlock()
	}
}
