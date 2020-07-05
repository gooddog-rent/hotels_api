package io

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestWatchFileValidChange(t *testing.T) {

	tempdir := "./"
	tempfile := "temp.*.json"

	dirname, err := ioutil.TempDir(tempdir, "test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dirname)

	f, err := ioutil.TempFile(dirname, tempfile)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())

	done := make(chan struct{}, 2)

	f.WriteString("some new data in json file")

	go func() {

		ok := watchFile(f.Name(), 150*time.Millisecond)
		if ok != nil {
			t.Errorf("Invalid file change. Should be nil return. Error: %v\n", ok)
		}
		done <- struct{}{}
	}()

	time.Sleep(100 * time.Millisecond)
	go func() {

		// write to new temp file some data
		f.Truncate(5)
		done <- struct{}{}
	}()

	<-done
	<-done

	// f.WriteString(`{"locations":[{"region":"Bavaro","hotels":["Airbnb Apartments","Catalonia Bavaro"]}]}`)
}

func TestWatchFileWrongFile(t *testing.T) {

	tempdir := "./"
	tempfile := "temp.*.json"

	dirname, err := ioutil.TempDir(tempdir, "test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dirname)

	f, err := ioutil.TempFile(dirname, tempfile)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())

	ok := watchFile("", 0)
	if ok == nil {
		t.Errorf("Invalid file. Error: %v\n", ok)
	}
}
