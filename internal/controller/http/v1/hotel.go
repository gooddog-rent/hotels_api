package http

import (
	"encoding/json"
	"hotels_api/internal/service"
	"log"
	"net/http"
)

type contextKey string

const contextQueryKey contextKey = "query"
const contextLimitKey contextKey = "limit"

type hotelRoutes struct {
	hotelService service.Hotel
}

// GetHotels method handle requests and response with filtered hotels
func (h *hotelRoutes) getHotels(w http.ResponseWriter, req *http.Request) {

	// all request validations are located in separate validateRequest middleware
	// get query and limit params from context
	text := req.Context().Value(contextQueryKey).(string)
	limit := req.Context().Value(contextLimitKey).(int)

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
