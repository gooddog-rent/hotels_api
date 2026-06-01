package service

import (
	"context"
	"hotels_api/internal/domain"
	infra "hotels_api/internal/infrastructure"
)

type Hotel interface {
	SearchHotels(ctx context.Context, text string, limit int) (*domain.Locations, error)
}

type Services struct {
	Searcher infra.Searcher
}

type ServicesDependencies struct {
	Repos *infra.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Searcher: NewHotelService(deps.Repos.Searcher),
	}
}
