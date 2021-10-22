package middleware

import (
	"log"
	"net/http"
	"time"
)

// CustomHeaders middleware handler setup CORS and other headers
func CustomHeaders(cacheDuration string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		d, err := time.ParseDuration(cacheDuration)
		if err != nil {
			log.Println(err)
			http.Error(w, "500 Internal Server Error!", http.StatusInternalServerError)
			return
		}

		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // change to host domain for private
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// JSON and Cache headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "private, max-age="+d.String())

		next.ServeHTTP(w, req)
	})
}
