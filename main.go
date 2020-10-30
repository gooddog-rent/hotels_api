// docker build -t iqhater/hotels_api .
// docker run -p 4040:4000 iqhater/hotels_api:latest

// test requests from cli
// http GET http://localhost:4000/hotels query==resort limit==10
// curl -X GET 'http://localhost:4000/hotels?query=resort&limit=10'

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"hotels_api/api"
	"hotels_api/io"
	mid "hotels_api/middleware"

	"github.com/joho/godotenv"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// init is invoked before main()
func init() {

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found. Env variables should be loaded.")
	}
}

func main() {

	// get port number
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Println("Env variable PORT does not exist!")
	}
	port = ":" + port

	// get hotels file
	filepath, exists := os.LookupEnv("HOTELS_PATH")
	if !exists {
		log.Println("Env variable HOTELS_PATH does not exist!")
	}

	ch := make(chan []byte)

	// watch hotels.json file change and auto update hotels
	go io.UpdateHotels(ch, filepath)

	l := api.Locations{
		Counter: api.Counter{
			Metrics: api.Metrics{
				Counter: api.NewMetrics().Counter,
			},
		},
	}
	prom.MustRegister(l.Metrics.Counter)

	// parse json
	go io.ParseJSON(ch, &l)

	mux := http.NewServeMux()
	mux.HandleFunc("/hotels", mid.ShowLog(l.Hotels))
	mux.Handle("/metrics", promhttp.Handler()) // prometheus metrics

	fmt.Printf("Hotels API listening requests on port %s\n", port)

	log.Fatal(http.ListenAndServe(port, api.Limit(mux)))
}
