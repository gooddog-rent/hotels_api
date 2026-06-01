package content

import (
	"net/http"
)

// ContentTypeHeaders middleware handler setup content type headers
func ContentTypeHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// JSON type headers
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, req)
	})
}
