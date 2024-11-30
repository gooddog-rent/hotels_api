package treedb

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	TEMP_DIR     = "../../../test/temp/"
	TEMP_FILE    = "temp.*.json"
	TEST_PATTERN = "test_watcher_service_"
)

func createTempFolderAndFile(t *testing.T) *os.File {

	tempDir, err := filepath.Abs(TEMP_DIR)
	if err != nil {
		t.Error(err)
	}

	dirname, err := os.MkdirTemp(tempDir, TEST_PATTERN)
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

	_, _ = f.WriteString("some new data in json file")

	w := WatcherRepo{
		filepath:   f.Name(),
		updateTime: 150 * time.Millisecond,
	}

	go func() {

		ok := w.WatchFile()
		if ok != nil {
			t.Errorf("Invalid file change. Should be nil return. Error: %v\n", ok)
		}
		done <- struct{}{}
	}()

	time.Sleep(100 * time.Millisecond)
	go func() {

		// write to new temp file some data
		_ = f.Truncate(5)
		done <- struct{}{}
	}()

	<-done
	<-done
}

func TestWatchFileWrongFile(t *testing.T) {

	_ = createTempFolderAndFile(t)

	w := WatcherRepo{
		filepath:   "",
		updateTime: 0,
	}

	ok := w.WatchFile()
	if ok == nil {
		t.Errorf("Invalid file. Error: %v\n", ok)
	}
}
