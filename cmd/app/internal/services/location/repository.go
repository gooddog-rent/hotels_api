package location

import (
	"hotels_api/cmd/app/internal/domain"
	"net/http"
)

type HotelRepository interface {
	GetHotels(w http.ResponseWriter, req *http.Request)
	Build()
	Search(text string, limit int) domain.ResponseHotels
}
