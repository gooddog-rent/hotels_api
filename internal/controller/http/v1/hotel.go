package http

import (
	"encoding/json"
	"net/http"

	"github.com/gooddog-rent/hotels_api/internal/service"
)

type contextKey string

const contextQueryKey contextKey = "query"
const contextLimitKey contextKey = "limit"

type hotelRoutes struct {
	hotelService service.Hotel
}

// ResponseLocations struct dto store output response regions and hotels data
type ResponseLocations struct {
	Hotels map[string][]string `json:"hotels"`
}

func NewHotelController(hotelService service.Hotel) *hotelRoutes {
	return &hotelRoutes{
		hotelService: hotelService,
	}
}

// GetHotels method handle requests and response with filtered hotels
func (h *hotelRoutes) GetHotels(w http.ResponseWriter, req *http.Request) {

	// all request validations are located in separate validateRequest middleware
	// get query and limit params from context
	text, ok := req.Context().Value(contextQueryKey).(string)
	if !ok {
		http.Error(w, "query value must be string type!", http.StatusBadRequest)
		return
	}

	limit, ok := req.Context().Value(contextLimitKey).(int)
	if !ok {
		http.Error(w, "limit value must be integer type!", http.StatusBadRequest)
		return
	}

	filtered, err := h.hotelService.SearchHotels(req.Context(), text, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// convert to response dto
	responseHotels := &ResponseLocations{
		Hotels: filtered.Hotels,
	}

	out, err := json.Marshal(responseHotels)
	if err != nil {
		http.Error(w, "500 Internal Server Error!", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
