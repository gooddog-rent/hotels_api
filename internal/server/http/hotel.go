package http

import (
	"encoding/json"
	"hotels_api/internal/service"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type hotelRoutes struct {
	hotelService service.Hotel
}

// GetHotels method handle requests and response with filtered hotels
func (h *hotelRoutes) getHotels(w http.ResponseWriter, req *http.Request) {

	// all request validations are located in separate validateRequest middleware
	// parse encoded query params to map
	params, _ := url.ParseQuery(req.URL.Query().Encode())

	// parse search query
	text := params.Get("query")

	// convert limit value to string
	limit, _ := strconv.Atoi(params.Get("limit"))

	filtered := h.hotelService.SearchHotels(text, limit)

	out, err := json.Marshal(filtered)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error!", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(out)
	if err != nil {
		log.Println("Response result write error!", err)
	}
}
