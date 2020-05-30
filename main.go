// docker build -t hotels-api .
// docker run -p 4000:4000 hotels-api

// test requests from cli
// http GET http://localhost:4000/hotels query==resort limit==10
// curl -X GET 'http://localhost:4000/hotels?query=resort&limit=10'

//TODO: add tls (generate certificate with mkcert). List of hotels is a public non-sensitive data. Does it really needed?
//TODO: add a few required metrics for prometheus
//TODO: add API and services tests
//TODO: add API credentials such as JWT token

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"hotels_api/api"
	"hotels_api/io"

	"github.com/joho/godotenv"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// init is invoked before main()
func init() {

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	/* apiKey, exists := os.LookupEnv("HOTELS_API_KEY")
	if !exists {
		log.Println("Env variable HOTELS_API_KEY does not exist!")
	} */

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

	// update hotels
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
	mux.HandleFunc("/hotels", l.Hotels)
	mux.Handle("/metrics", promhttp.Handler())

	fmt.Printf("Hotels API listening requests on port %s\n", port)

	log.Fatal(http.ListenAndServe(port, api.Limit(mux)))
}
