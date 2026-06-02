package internal

import (
	"context"
	"errors"
	"fmt"
	ctrl "hotels_api/internal/controller/http/v1"
	infra "hotels_api/internal/infrastructure"

	"hotels_api/internal/service"
	"hotels_api/pkg/cache"
	"hotels_api/pkg/headers/content"
	"hotels_api/pkg/headers/cors"
	"hotels_api/pkg/headers/secure"
	"hotels_api/pkg/limiter"
	"hotels_api/pkg/logger"
	"os"
	"time"

	mid "hotels_api/pkg/middleware"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/sh7dm/brotlihandler"
)

func Run() {

	// init config environments
	cfg := NewConfig()

	// init cache storage
	c := cache.NewCache("2m")
	cache.CacheStore = c

	// init repositories
	repositories, err := infra.NewRepositories(cfg.HOTELS_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer repositories.DBConnection.DB.Close()

	// init services dependencies
	deps := service.ServicesDependencies{
		Repos: repositories,
	}

	// init services
	services := service.NewServices(deps)

	// init controllers
	controller := ctrl.NewHotelController(services.Searcher)

	// init server
	mux := http.NewServeMux()

	// init middlewares per one route
	midMain := mid.Middlewares(
		c.CacheHeaders,
		content.ContentTypeHeaders,
		cors.CORSHeaders,
		secure.SecureHeaders,
		ctrl.ValidateRequest,
		// c.CacheResponse,
	)

	// init global middlewares
	middlewaresGlobal := mid.Middlewares(
		logger.Log,
		limiter.Limit,
		brotlihandler.CompressHandler,
	)

	mux.HandleFunc("GET /api/v1/hotels", mid.Bind(midMain, controller.GetHotels))

	// best practice to use timeout
	server := &http.Server{
		Addr:              ":" + cfg.HTTP_PORT,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		Handler:           middlewaresGlobal(mux),
	}

	// start server
	go func() {
		fmt.Printf("🏨 Hotels API service listening requests on port: %s\n", cfg.HTTP_PORT)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v\n", err)
		}
	}()

	// graceful shutdown
	ctx, exit := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer exit()

	<-ctx.Done()
	fmt.Println("graceful shutting down starting...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("failed to shutdown server: %v\n", err)
	} else {
		log.Println("Hotels API server successfully shutdown")
	}
}
