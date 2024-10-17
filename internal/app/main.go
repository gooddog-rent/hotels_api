package app

import (
	"hotels_api/internal/server/http"
)

func Run() {

	// init config environments
	cfg := NewConfig()

	// init logger

	// init db in-memory search tree

	// init repositories

	// init dependencies

	// init services

	// init server
	// run server
	http.RunServer(cfg.HTTP_PORT, cfg.HOTELS_PATH)
}
