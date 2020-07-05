package io

import (
	"encoding/json"
	"hotels_api/api"
	"io/ioutil"
	"log"
	"os"
	"time"
)

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

func UpdateHotels(ch chan []byte, filepath string) {
	for {
		h, err := ioutil.ReadFile(filepath)
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

func ParseJSON(ch chan []byte, locations *api.Locations) {
	for {
		err := json.Unmarshal(<-ch, &locations)
		if err != nil {
			log.Fatal(err)
		}
	}
}
