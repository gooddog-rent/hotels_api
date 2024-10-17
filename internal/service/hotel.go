package service

import (
	"hotels_api/internal/entity"
	"hotels_api/internal/infra"
)

type HotelService struct {
	hotelRepo infra.Hotel
}

func NewHotelService(hotelRepo infra.Hotel) *HotelService {
	return &HotelService{hotelRepo: hotelRepo}
}

func (h *HotelService) SearchHotels(text string, limit int) *entity.ResponseHotels {
	return h.hotelRepo.SearchHotels(text, limit)
}
