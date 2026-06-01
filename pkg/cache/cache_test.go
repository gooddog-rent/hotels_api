package cache

import (
	"testing"

	"time"

	"bytes"
)

// helper function
func parse(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

// helper function
func assertContentEquals(t *testing.T, content, expected []byte) {
	if !bytes.Equal(content, expected) {
		t.Errorf("content should '%s', but was '%s'", expected, content)
	}
}

func TestGetEmpty(t *testing.T) {

	// Arrange
	storage := NewCache("5s")

	// Act
	content := storage.Get("MY_KEY")

	// Assert
	assertContentEquals(t, content, []byte(""))
}

func TestGetValue(t *testing.T) {

	// Arrange
	storage := NewCache("5s")
	storage.Set("MY_KEY", []byte("123456"), parse("5s"))

	// Act
	content := storage.Get("MY_KEY")

	// Assert
	assertContentEquals(t, content, []byte("123456"))
}

func TestGetExpiredValue(t *testing.T) {

	// Arrange
	storage := NewCache("5s")
	storage.Set("MY_KEY", []byte("123456"), parse("1s"))
	time.Sleep(parse("1s200ms"))

	// Act
	content := storage.Get("MY_KEY")

	// Assert
	assertContentEquals(t, content, []byte(""))
}

func TestGetValueAfterSet(t *testing.T) {

	// Arrange
	storage := NewCache("5s")
	storage.Set("MY_KEY", []byte("123456"), parse("5s"))
	time.Sleep(parse("1s"))

	// Act
	content := storage.Get("MY_KEY")

	// Assert
	assertContentEquals(t, content, []byte("123456"))
}
