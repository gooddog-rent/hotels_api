package api

import "sync"

// Counter struct store metrics for prometheus
// and sync mutex for mutate api counter
type Counter struct {
	counter uint
	mu      sync.Mutex
	Metrics
}
