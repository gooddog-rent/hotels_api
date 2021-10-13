package middleware

import (
	"hotels_api/api"
	"log"
	"net/http"
	"time"
)

// showLog middleware handler shows network data log info
func ShowLog(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := time.Now()

		sr := api.NewStatusHTTP(w)
		next.ServeHTTP(sr, req)

		statusCode := sr.StatusCode
		log.Printf("%s %s - %s - [%d: %s] - %v\n", req.Method, req.URL.String(), req.RemoteAddr, statusCode, http.StatusText(statusCode), time.Since(t))
	})
}
