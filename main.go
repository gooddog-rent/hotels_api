// docker build -t hotels-api .
// docker run -p 4000:4000 hotels-api

//TODO: add tls (generate certificate with mkcert)
//TODO: add logger exporter for prometheus system
//TODO: add API and services tests
//TODO: add API credentials such as JWT token
//TODO: add .env file for external config (port, path, api key and so on)

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	ch := make(chan []byte)

	// update hotels
	go updateHotels(ch)

	l := Locations{}

	// parse json
	go l.parseJSON(ch)

	mux := http.NewServeMux()
	mux.HandleFunc("/hotels", l.hotels)

	const port = ":4000"
	fmt.Printf("Hotels API listening requests on port %s\n", port)

	log.Fatal(http.ListenAndServe(port, limit(mux)))
}
