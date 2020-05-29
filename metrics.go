package main

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	counter prom.Counter
}

func NewMetrics() *Metrics {
	return &Metrics{
		prom.NewCounter(prom.CounterOpts{
			Name: "hotels_api_total_requests",
		}),
	}
}
