package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Hotels method handle requests and response with filtered hotels
func (l *Locations) Hotels(w http.ResponseWriter, req *http.Request) {

	t := time.Now()

	l.mu.Lock()
	l.counter++
	l.mu.Unlock()

	// prometheus api conuter metric increment
	l.Metrics.Counter.Inc()

	// check on correct request method
	if req.Method != http.MethodGet {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.RemoteAddr, http.StatusMethodNotAllowed, time.Since(t))
		http.Error(w, "Wrong method used! Only GET method allowed.", http.StatusMethodNotAllowed)
		return
	}

	// parse encoded query params to map
	params, err := url.ParseQuery(req.URL.Query().Encode())
	if err != nil {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.RemoteAddr, http.StatusBadRequest, time.Since(t))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if "query" word is exists
	if _, ok := params["query"]; !ok {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.RemoteAddr, http.StatusBadRequest, time.Since(t))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if "limit" word is exists
	if _, ok := params["limit"]; !ok {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.RemoteAddr, http.StatusBadRequest, time.Since(t))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// parse search query
	text := params.Get("query")
	if text == "" {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.RemoteAddr, http.StatusBadRequest, time.Since(t))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// convert limit value to string
	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.RemoteAddr, http.StatusBadRequest, time.Since(t))
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filtered := l.filterByHotels(text, limit)

	out, err := json.Marshal(filtered)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // change to host domain
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Write(out)
}
