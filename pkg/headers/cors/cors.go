package cors

import (
	"net/http"
)

// CustomHeaders middleware handler setup CORS and other headers
func CORSHeaders(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // change to host domain for private
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

		next.ServeHTTP(w, req)
	})
}
