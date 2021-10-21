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
	items map[string]Item
	mu    *sync.RWMutex
}

// NewCache creates a new in memory Cache
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]Item),
		mu:    &sync.RWMutex{},
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
func CacheResponse(duration string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		d, err := time.ParseDuration(duration)
		if err != nil {
			log.Println(err)
			http.Error(w, "500 Internal Server Error!", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Cache-Control", "private, max-age="+d.String())

		content := CacheStore.Get(req.RequestURI)
		if content != nil {
			log.Println("Cached response.")
			w.Write(content)
		} else {
			c := httptest.NewRecorder()
			next(c, req)

			for k, v := range c.HeaderMap {
				w.Header()[k] = v
			}

			w.WriteHeader(c.Code)
			content := c.Body.Bytes()

			if d, err := time.ParseDuration(duration); err == nil {
				log.Printf("New page cached: %s for %s\n", req.RequestURI, duration)
				CacheStore.Set(req.RequestURI, content, d)
			} else {
				log.Printf("Page not cached. err: %s\n", err)
			}
			w.Write(content)
		}
	})
}
