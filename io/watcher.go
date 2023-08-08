package io

import (
	"encoding/json"
	"hotels_api/api"
	"log"
	"os"
	"time"
)

// watchFile function watch for file changes
func watchFile(filepath string, updatetime time.Duration) error {

	initialStat, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filepath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}
		time.Sleep(updatetime)
	}
	return nil
}

// UpdateHotels function watch to hotels.json file changes
// and check json is valid then send updateed data into channel
func UpdateHotels(ch chan []byte, filepath string) {
	for {
		h, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		if json.Valid(h) {
			ch <- h
		} else {
			log.Println("Invalid JSON syntax in hotels.json file!")
		}

		err = watchFile(filepath, 15*time.Second)
		if err != nil {
			log.Println(err)
		}

		log.Println("File has been changed")
	}
}

// ParseJSON function gets data byte channel with locations
// and unmarshal data in json
func ParseJSON(ch chan []byte, locations *api.Locations) {

	for {
		err := json.Unmarshal(<-ch, &locations)
		if err != nil {
			log.Fatal(err)
		}

		locations.BuildSuffixTree()
	}
}
