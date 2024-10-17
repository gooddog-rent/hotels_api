package infra

import (
	"hotels_api/internal/entity"
	"hotels_api/internal/infra/treedb"
	"hotels_api/pkg/suffixtree"
	"time"
)

type Tree interface {
	BuildTree(l *entity.LocationsStore) *entity.SearchTree
}

type Hotel interface {
	SearchHotels(text string, limit int) *entity.ResponseHotels
	FileWatcher
}

type Watcher interface {
	WatchFile() error
}

type Parser interface {
	Update()
	Parse(locations *entity.LocationsStore)
}

type FileWatcher interface {
	Watcher
	Parser
}

type Repositories struct {
	// Tree
	Hotel
	// FileWatcher
}

func NewRepositories(tree *suffixtree.SuffixTree, HOTELS_PATH string) *Repositories {
	return &Repositories{
		// Tree:        treedb.NewTreeRepo(tree),
		Hotel: treedb.NewWatcherRepo(tree, HOTELS_PATH, 15*time.Second),
		// FileWatcher: treedb.NewWatcherRepo(tree, HOTELS_PATH, 15*time.Second),
	}
}
