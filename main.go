// docker build -t iqhater/hotels_api .
// docker run -p 4040:4000 iqhater/hotels_api:latest

// test requests from cli
// http GET http://localhost:4000/hotels query==resort limit==10
// curl -X GET 'http://localhost:4000/hotels?query=resort&limit=10'

package main

import (
	"fmt"
	"log"
	"net/http"

	"hotels_api/api"
	"hotels_api/config"
	"hotels_api/io"
	mid "hotels_api/middleware"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found. Env variables should be loaded.")
	}
}

func main() {

	// init config environments
	cfg := config.NewConfig()

	ch := make(chan []byte)

	// watch hotels.json file change and auto update hotels
	go io.UpdateHotels(ch, cfg.HOTELS_PATH)

	l := api.Locations{}

	// parse json
	go io.ParseJSON(ch, &l)

	mux := http.NewServeMux()
	mux.HandleFunc("/hotels", mid.ShowLog(l.GetHotels))

	fmt.Printf("Hotels API listening requests on port: %s\n", cfg.HTTP_PORT)

	log.Fatal(http.ListenAndServe(":"+cfg.HTTP_PORT, mid.Limit(mux)))
}
