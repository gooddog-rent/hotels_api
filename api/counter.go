package api

import "sync"

type Counter struct {
	counter uint
	mu      sync.Mutex
	Metrics
}
