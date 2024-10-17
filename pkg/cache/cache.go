package cache

import (
	"sync"
	"time"
)

// static interface implementation check for convinience
var _ Storage = new(Cache)

var CacheStore Storage

// Item is a cached reference
type Item struct {
	Content    []byte
	Expiration int64
}

// Expired returns true if the item has expired.
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// Cache struct for caching strings in memory
type Cache struct {
	items    map[string]Item
	mu       *sync.RWMutex
	duration string
}

// Get a cached content by key
func (c Cache) Get(key string) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item := c.items[key]
	if item.Expired() {
		delete(c.items, key)
		return nil
	}
	return item.Content
}

// Set a cached content by key
func (c Cache) Set(key string, content []byte, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Content:    content,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}

// NewCache creates a new in memory Cache
func NewCache() *Cache {
	return &Cache{
		items:    make(map[string]Item),
		mu:       &sync.RWMutex{},
		duration: "2m",
	}
}
