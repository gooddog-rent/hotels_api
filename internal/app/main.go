package app

import (
	"hotels_api/internal/controller/http/v1"
)

func Run() {

	// init config environments
	cfg := NewConfig()

	// init server
	http.RunServer(cfg.HTTP_PORT, cfg.HOTELS_PATH)
}
