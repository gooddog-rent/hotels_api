package main

import (
	app "hotels_api/internal"
)

func main() {

	// init .env file
	app.InitEnv(".env")

	// run app
	app.Run()
}
