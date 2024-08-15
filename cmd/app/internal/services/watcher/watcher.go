package watcher

import (
	"encoding/json"
	"hotels_api/cmd/app/internal/domain"
	"hotels_api/cmd/app/internal/usecases"
	"log"
	"os"
	"time"
)

// static interface implementation check for convinience
var _ FileWatcher = (*Watch)(nil)

type Watch struct {
	Filepath   string
	UpdateTime time.Duration
	Ch         chan []byte
	// Locations    domain.Locations
	// HotelUsecase usecases.HotelUsecase
}

// watchFile function watch for file changes
func (w *Watch) watchFile() error {

	initialStat, err := os.Stat(w.Filepath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(w.Filepath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}
		time.Sleep(w.UpdateTime)
	}
	return nil
}

// Update function watch to hotels.json file changes
// and check json is valid then send updated data into channel
func (w *Watch) Update() {

	for {
		h, err := os.ReadFile(w.Filepath)
		if err != nil {
			log.Fatal(err)
		}
		if json.Valid(h) {
			w.Ch <- h
		} else {
			log.Println("Invalid JSON syntax in hotels.json file!")
		}

		err = w.watchFile()
		if err != nil {
			log.Println(err)
		}

		log.Println("File has been changed")
	}
}

// Parse function gets data byte channel with locations
// and unmarshal data in json
func (w *Watch) Parse(locations *domain.Locations) {

	for {
		err := json.Unmarshal(<-w.Ch, locations)
		if err != nil {
			log.Fatal(err)
		}

		var hotelUsecase usecases.HotelUsecase
		hotelUsecase.Build()
	}
}
