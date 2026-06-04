package service

import (
	"context"

	"github.com/gooddog-rent/hotels_api/internal/domain"
	infra "github.com/gooddog-rent/hotels_api/internal/infrastructure"
)

type HotelService struct {
	hotelRepo infra.Searcher
}

func NewHotelService(hotelRepo infra.Searcher) *HotelService {
	return &HotelService{hotelRepo: hotelRepo}
}

func (h *HotelService) SearchHotels(ctx context.Context, text string, limit int) (*domain.Locations, error) {
	return h.hotelRepo.SearchHotels(ctx, text, limit)
}
