package usecases

import "hotels_api/cmd/app/internal/domain"

type HotelUsecase interface {
	Build()
	Search(text string, limit int) domain.ResponseHotels
}
