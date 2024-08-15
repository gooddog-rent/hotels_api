package main

import (
	"hotels_api/cmd/app"
)

func main() {

	// init config environments
	cfg := app.NewConfig()

	app.Run(cfg)
}
