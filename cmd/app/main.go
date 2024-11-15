package main

import (
	"hotels_api/internal/app"
)

func main() {

	// init .env file
	app.InitEnv(".env")

	app.Run()
}
