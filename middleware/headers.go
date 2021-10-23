package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// CustomHeaders middleware handler setup CORS and other headers
func (c Cache) CustomHeaders(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		duration, err := time.ParseDuration(c.duration)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // change to host domain for private
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// JSON and Cache headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "private, max-age="+fmt.Sprintf("%.0f", duration.Seconds()))

		// secure headers
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, req)
	})
}
