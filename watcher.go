package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func watchFile(filePath string) error {

	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func updateHotels(ch chan []byte, filepath string) {
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

		err = watchFile(filepath)
		if err != nil {
			log.Println(err)
		}

		log.Println("File has been changed")
	}
}

func (l *Locations) parseJSON(ch chan []byte) {
	for {
		err := json.Unmarshal(<-ch, &l)
		if err != nil {
			log.Fatal(err)
		}
	}
}
