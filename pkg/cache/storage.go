package cache

import "time"

// Storage interface for cache middleware implementation
type Storage interface {
	Get(key string) []byte
	Set(key string, content []byte, duration time.Duration)
}
