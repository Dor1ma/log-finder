package repository

import (
	"container/list"
	"os"
	"sync"
	"time"

	"github.com/Dor1ma/log-finder/pkg/mmap"
)

type fileCache struct {
	maxSize int
	ttl     time.Duration
	cache   map[string]*cacheEntry
	lruList *list.List
	mutex   sync.Mutex
}

type cacheEntry struct {
	path      string
	data      []byte
	expiresAt time.Time
	element   *list.Element
}

func NewFileCache(maxSize int, ttl time.Duration) *fileCache {
	return &fileCache{
		maxSize: maxSize,
		ttl:     ttl,
		cache:   make(map[string]*cacheEntry),
		lruList: list.New(),
	}
}

func (c *fileCache) Get(path string) ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if entry, exists := c.cache[path]; exists {
		if time.Now().After(entry.expiresAt) {
			c.removeEntry(entry)
			return nil, os.ErrNotExist
		}
		c.lruList.MoveToFront(entry.element)
		return entry.data, nil
	}

	data, err := mmap.MapFile(path)
	if err != nil {
		return nil, err
	}

	entry := &cacheEntry{
		path:      path,
		data:      data,
		expiresAt: time.Now().Add(c.ttl),
		element:   c.lruList.PushFront(path),
	}

	c.cache[path] = entry

	if len(c.cache) > c.maxSize {
		c.evictOldest()
	}

	return data, nil
}

func (c *fileCache) evictOldest() {
	oldest := c.lruList.Back()
	if oldest != nil {
		c.removeEntry(c.cache[oldest.Value.(string)])
	}
}

func (c *fileCache) removeEntry(entry *cacheEntry) {
	delete(c.cache, entry.path)
	c.lruList.Remove(entry.element)
	mmap.Unmap(entry.data)
}
