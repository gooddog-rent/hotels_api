package location

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type LocationService struct {
	HotelRepository
}

// GetHotels method handle requests and response with filtered hotels
func (l *LocationService) GetHotels(w http.ResponseWriter, req *http.Request) {

	// all request validations are located in separate validateRequest middleware
	// parse encoded query params to map
	params, _ := url.ParseQuery(req.URL.Query().Encode())

	// parse search query
	text := params.Get("query")

	// convert limit value to string
	limit, _ := strconv.Atoi(params.Get("limit"))

	filtered := l.Search(text, limit)

	out, err := json.Marshal(filtered)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error!", http.StatusInternalServerError)
		return
	}

	w.Write(out)
}
