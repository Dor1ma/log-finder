package service

import (
	"sync"
	"time"
)

type TTLCache struct {
	cache map[string]cacheEntry
	mutex sync.RWMutex
	ttl   time.Duration
}

type cacheEntry struct {
	value      string
	expiration time.Time
}

func NewTTLCache(ttl time.Duration) *TTLCache {
	c := &TTLCache{
		cache: make(map[string]cacheEntry),
		ttl:   ttl,
	}
	go c.cleanup()
	return c
}

func (c *TTLCache) Get(key string) (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.cache[key]
	if !exists || time.Now().After(entry.expiration) {
		return "", false
	}
	return entry.value, true
}

func (c *TTLCache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[key] = cacheEntry{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}

func (c *TTLCache) cleanup() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, entry := range c.cache {
			if now.After(entry.expiration) {
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}
