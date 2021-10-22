package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

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

// NewCache creates a new in memory Cache
func NewCache() *Cache {
	return &Cache{
		items:    make(map[string]Item),
		mu:       &sync.RWMutex{},
		duration: "10s",
	}
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

//Set a cached content by key
func (c Cache) Set(key string, content []byte, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Content:    content,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}

// CacheResponse middleware
func (c *Cache) CacheResponse(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		content := CacheStore.Get(req.RequestURI)
		if content != nil {
			// log.Println("Cached response.")
			w.Write(content)
		} else {
			rr := httptest.NewRecorder()

			//TODO: fix to dublicate here from customHeaders middleware
			rr.Header().Set("Content-Type", "application/json")

			next(rr, req)

			for k, v := range rr.HeaderMap {
				w.Header()[k] = v
			}

			w.WriteHeader(rr.Code)
			content := rr.Body.Bytes()

			if d, err := time.ParseDuration(c.duration); err == nil {

				// log.Printf("New page cached: %s for %s\n", req.RequestURI, c.duration)
				CacheStore.Set(req.RequestURI, content, d)
			} else {
				log.Printf("Page not cached. err: %s\n", err)
			}
			w.Write(content)
		}
	})
}
