package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// GetHotels method handle requests and response with filtered hotels
func (l *Locations) GetHotels(w http.ResponseWriter, req *http.Request) {

	// check on correct request method
	if req.Method != http.MethodGet {
		http.Error(w, "Wrong method used! Only GET method allowed.", http.StatusMethodNotAllowed)
		return
	}

	// parse encoded query params to map
	params, err := url.ParseQuery(req.URL.Query().Encode())
	if err != nil {
		http.Error(w, "Invalid query parameters!", http.StatusBadRequest)
		return
	}

	// check if "query" word is exists
	if _, ok := params["query"]; !ok {
		http.Error(w, "'query' parameter does not exist or bad value!", http.StatusBadRequest)
		return
	}

	// check if "limit" word is exists
	if _, ok := params["limit"]; !ok {
		http.Error(w, "'limit' parameter does not exist or bad value!", http.StatusBadRequest)
		return
	}

	// parse search query
	text := params.Get("query")
	if text == "" {
		http.Error(w, "'query' value is empty!", http.StatusBadRequest)
		return
	}

	// convert limit value to string
	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil {
		log.Println(err)
		http.Error(w, "'limit' value must be integer!", http.StatusBadRequest)
		return
	}

	filtered := l.SearchInSuffixTree(text, limit)

	out, err := json.Marshal(filtered)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // change to host domain for private
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Write(out)
}
