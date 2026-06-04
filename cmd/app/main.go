package main

import (
	app "github.com/gooddog-rent/hotels_api/internal"
)

func main() {

	// init .env file
	app.InitEnv(".env")

	// run app
	app.Run()
}
