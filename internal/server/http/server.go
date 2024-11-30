package http

import (
	"context"
	"fmt"
	"hotels_api/internal/entity"
	"hotels_api/internal/infra"

	"hotels_api/internal/service"
	"hotels_api/pkg/cache"
	"hotels_api/pkg/headers/content"
	"hotels_api/pkg/headers/cors"
	"hotels_api/pkg/headers/secure"
	"hotels_api/pkg/limiter"
	"hotels_api/pkg/logger"
	"hotels_api/pkg/suffixtree"
	"os"
	"time"

	mid "hotels_api/pkg/middleware"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/sh7dm/brotlihandler"
)

func RunServer(HTTP_PORT, HOTELS_PATH string) {

	// init cache storage
	c := cache.NewCache()
	cache.CacheStore = c

	// init db in-memory search tree
	tree := suffixtree.NewGeneralizedSuffixTree()

	// init repositories
	repositories := infra.NewRepositories(tree, HOTELS_PATH)

	locations := &entity.LocationsStore{}

	go repositories.Hotel.Parse(locations)
	go repositories.Hotel.Update()

	// init services dependencies
	deps := service.ServicesDependencies{
		Repos: repositories,
	}

	// init services
	services := service.NewServices(deps)

	r := &hotelRoutes{
		hotelService: services.Hotel,
	}

	// init server
	mux := http.NewServeMux()

	// init middlewares per one route
	middlewaresPerRoute := mid.Middlewares(
		logger.Log,
		c.CacheHeaders,
		content.ContentTypeHeaders,
		cors.CORSHeaders,
		secure.SecureHeaders,
		ValidateRequest,
		c.CacheResponse,
	)

	// init global middlewares
	middlewaresGlobal := mid.Middlewares(
		limiter.Limit,
		brotlihandler.CompressHandler,
	)

	mux.HandleFunc("/hotels", mid.Bind(middlewaresPerRoute, r.getHotels))
	fmt.Printf("🏨 Hotels API service listening requests on port: %s\n", HTTP_PORT)

	// best practice to use timeout
	server := &http.Server{
		Addr:              ":" + HTTP_PORT,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           middlewaresGlobal(mux),
	}
	log.Fatal(server.ListenAndServe())

	// graceful shutdown
	ctx, exit := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer exit()
	<-ctx.Done()
}
