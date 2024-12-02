package service

import (
	"hotels_api/internal/entity"
	infra "hotels_api/internal/infrastructure"
)

type Hotel interface {
	SearchHotels(text string, limit int) *entity.ResponseHotels
}

type Services struct {
	Hotel Hotel
}

type ServicesDependencies struct {
	Repos *infra.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Hotel: NewHotelService(deps.Repos.Hotel),
	}
}
