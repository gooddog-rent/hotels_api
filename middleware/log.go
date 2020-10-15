package middleware

import (
	"log"
	"net/http"
	"time"
)

// showLog middleware handler shows network data log info
func ShowLog(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.RemoteAddr, http.StatusOK, time.Since(t))
	})
}
