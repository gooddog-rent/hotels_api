package infrastructure

import (
	"encoding/json"
	"errors"
	"hotels_api/internal/entity"
	"hotels_api/pkg/suffixtree"
	"log"
	"os"
	"sync"
	"time"
)

// static interface implementation check for convinience
var _ FileWatcher = (*WatcherRepo)(nil)

type WatcherRepo struct {
	tree *TreeRepo

	filepath     string
	updateTime   time.Duration
	dataFromFile chan []byte

	mu sync.RWMutex
}

func NewWatcherRepo(tree *suffixtree.SuffixTree, filepath string, updateTime time.Duration) *WatcherRepo {
	return &WatcherRepo{
		tree: &TreeRepo{
			searchTree: &entity.SearchTree{
				Tree:       tree,
				HotelsList: []string{},
			},
		},
		filepath:     filepath,
		updateTime:   updateTime,
		dataFromFile: make(chan []byte),
	}
}

// watchFile function watch for file changes
func (w *WatcherRepo) WatchFile() error {

	initialStat, err := os.Stat(w.filepath)
	if err != nil {
		return errors.Unwrap(err)
	}

	for {
		stat, err := os.Stat(w.filepath)
		if err != nil {
			return errors.Unwrap(err)
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}
		time.Sleep(w.updateTime)
	}
	return nil
}

// Update function watch to hotels.json file changes
// and check json is valid then send updated data into channel
func (w *WatcherRepo) Update() {

	for {
		content, err := os.ReadFile(w.filepath)
		if err != nil {
			log.Fatal(err)
		}
		if json.Valid(content) {
			w.dataFromFile <- content
		} else {
			log.Println("Invalid JSON syntax in hotels.json file!")
		}

		err = w.WatchFile()
		if err != nil {
			log.Println(err)
		}

		log.Println("File has been changed")
	}
}

// Parse function gets data byte channel with locations
// and unmarshal data in json
func (w *WatcherRepo) Parse(locations *entity.LocationsStore) {

	for {
		err := json.Unmarshal(<-w.dataFromFile, &locations)
		if err != nil {
			log.Fatal(err)
		}
		w.mu.Lock()
		w.tree.searchTree = w.tree.BuildTree(locations)
		w.mu.Unlock()
	}
}
