package service

import (
	"context"

	"github.com/gooddog-rent/hotels_api/internal/domain"
	infra "github.com/gooddog-rent/hotels_api/internal/infrastructure"
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
