package cache

import (
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

// CacheResponse middleware
func (c *Cache) CacheResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		content := CacheStore.Get(req.RequestURI)
		if content != nil {
			// log.Println("Cached response.")
			_, err := w.Write(content)
			if err != nil {
				log.Println("Cache write error! ", err)
			}

		} else {
			rr := httptest.NewRecorder()

			next.ServeHTTP(rr, req)

			for k, v := range rr.Header() {
				w.Header()[k] = v
			}

			w.WriteHeader(rr.Code)
			content := rr.Body.Bytes()

			if d, err := time.ParseDuration(c.duration); err == nil {

				// log.Printf("New data cached: %s for %s\n", req.RequestURI, c.duration)
				CacheStore.Set(req.RequestURI, content, d)
			} else {
				log.Printf("Data not cached. err: %s\n", err)
			}

			_, err := w.Write(content)
			if err != nil {
				log.Println("Cache write error! ", err)
			}
		}
	})
}
