package content

import (
	"net/http"
)

// CustomHeaders middleware handler setup CORS and other headers
func ContentTypeHeaders(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// JSON type headers
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, req)
	})
}
