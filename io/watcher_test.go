package io

import (
	"os"
	"testing"
	"time"
)

const (
	TEMP_DIR     = "./"
	TEMP_FILE    = "temp.*.json"
	TEST_PATTERN = "test"
)

func createTempFolderAndFile(t *testing.T) *os.File {

	dirname, err := os.MkdirTemp(TEMP_DIR, TEST_PATTERN)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dirname)

	f, err := os.CreateTemp(dirname, TEMP_FILE)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())

	return f
}

func TestWatchFileValidChange(t *testing.T) {

	f := createTempFolderAndFile(t)

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
}

func TestWatchFileWrongFile(t *testing.T) {

	_ = createTempFolderAndFile(t)

	ok := watchFile("", 0)
	if ok == nil {
		t.Errorf("Invalid file. Error: %v\n", ok)
	}
}
