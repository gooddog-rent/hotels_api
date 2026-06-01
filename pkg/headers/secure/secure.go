package secure

import (
	"net/http"
)

// SecureHeaders middleware handler setup CORS and other headers
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// HTTPS Enforcement
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Content Security Policies
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; frame-ancestors 'none'")

		// Prevent MIME Sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Referrer Policy
		w.Header().Set("Referrer-Policy", "no-referrer")

		// Cross-Origin Protection
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		// Cache Control
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, req)
	})
}
