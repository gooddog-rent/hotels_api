package service

import (
	"hotels_api/internal/entity"
	"hotels_api/internal/infra"
)

type Hotel interface {
	SearchHotels(text string, limit int) *entity.ResponseHotels
}

type Services struct {
	Hotel Hotel
	// Watcher Watcher
}

type ServicesDependencies struct {
	Repos *infra.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Hotel: NewHotelService(deps.Repos.Hotel),
		// Watcher: NewWatcherService(deps.Repos.FileWatcher),
	}
}
