package cache

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// CustomHeaders middleware handler setup CORS and other headers
func (c Cache) CacheHeaders(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		duration, err := time.ParseDuration(c.duration)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Cache-Control", "private, max-age="+fmt.Sprintf("%.0f", duration.Seconds()))

		next.ServeHTTP(w, req)
	})
}
