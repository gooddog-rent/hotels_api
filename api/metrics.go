package api

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

// Metrics struct for prometheus
type Metrics struct {
	Counter prom.Counter // increment api requests
}

func NewMetrics() *Metrics {
	return &Metrics{
		prom.NewCounter(prom.CounterOpts{
			Name: "hotels_api_total_requests",
		}),
	}
}
