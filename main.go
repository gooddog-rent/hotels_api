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
	"github.com/sh7dm/brotlihandler"
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

	// parse hotels json file
	go io.ParseJSON(ch, &l)

	// init cache storage
	cache := mid.NewCache()
	mid.CacheStore = cache

	mux := http.NewServeMux()
	mux.HandleFunc("/hotels", mid.ShowLog(cache.CustomHeaders(mid.ValidateRequest(cache.CacheResponse(l.GetHotels)))))

	fmt.Printf("Hotels API listening requests on port: %s\n", cfg.HTTP_PORT)

	log.Fatal(http.ListenAndServe(":"+cfg.HTTP_PORT, mid.Limit(brotlihandler.CompressHandler(mux))))
}
