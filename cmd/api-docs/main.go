package main

import (
	"fmt"
	"hotels_api/pkg/logger"
	mid "hotels_api/pkg/middleware"
	"log"
	"net/http"
	"time"
)

// OpenAPI docs
func main() {

	// serve api folder for openapi docs
	mux := http.NewServeMux()
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.Handle("/api/", http.StripPrefix("/api/", http.FileServer(http.Dir("./api"))))

	// redirect to from / to /api/app/redoc.html
	mux.HandleFunc("/", mid.Bind(logger.Log, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/app/redoc.html", http.StatusMovedPermanently)
	}))

	const HTTP_PORT = ":3200"

	fmt.Printf("📄 Hotels API Docs server started on port: %s\n", HTTP_PORT)

	// best practice to use timeout
	server := &http.Server{
		Addr:              HTTP_PORT,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		Handler:           mux,
	}
	log.Fatal(server.ListenAndServe())
}
