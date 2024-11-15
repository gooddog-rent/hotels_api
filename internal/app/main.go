package app

import (
	"hotels_api/internal/server/http"
)

func Run() {

	// init config environments
	cfg := NewConfig()

	// init server
	http.RunServer(cfg.HTTP_PORT, cfg.HOTELS_PATH)
}
