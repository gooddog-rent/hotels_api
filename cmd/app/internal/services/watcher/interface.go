package watcher

import (
	"hotels_api/cmd/app/internal/domain"
)

type Watcher interface {
	watchFile() error
}

type Parser interface {
	Update()
	Parse(locations *domain.Locations)
}

type FileWatcher interface {
	Watcher
	Parser
}
