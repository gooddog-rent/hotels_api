package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hotels_api/cmd/app/internal/domain"
	"hotels_api/cmd/app/internal/services/location"
	"hotels_api/cmd/app/internal/services/watcher"
	"hotels_api/internal/validator"
	"hotels_api/pkg/cache"
	"hotels_api/pkg/headers/content"
	"hotels_api/pkg/headers/cors"
	"hotels_api/pkg/headers/secure"
	"hotels_api/pkg/limiter"
	"hotels_api/pkg/logger"

	"github.com/sh7dm/brotlihandler"
)

func Run(cfg *Config) {

	w := &watcher.Watch{
		Filepath:   cfg.HOTELS_PATH,
		UpdateTime: 15 * time.Second,
		Ch:         make(chan []byte),
	}

	l := &location.LocationService{}

	// watch hotels.json file change and auto update hotels
	go w.Update()

	locations := domain.Locations{Locations: []domain.Location{}}

	// parse hotels json file
	go w.Parse(&locations)

	// init cache storage
	c := cache.NewCache()
	cache.CacheStore = c

	mux := http.NewServeMux()
	mux.HandleFunc("/hotels", logger.Log(c.CacheHeaders(content.ContentTypeHeaders(cors.CORSHeaders(secure.SecureHeaders(validator.ValidateRequest(c.CacheResponse(l.GetHotels))))))))

	fmt.Printf("Hotels API listening requests on port: %s\n", cfg.HTTP_PORT)

	log.Fatal(http.ListenAndServe(":"+cfg.HTTP_PORT, limiter.Limit(brotlihandler.CompressHandler(mux))))

	// graceful shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
}
