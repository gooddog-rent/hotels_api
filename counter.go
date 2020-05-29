package main

import "sync"

type Counter struct {
	counter uint
	mu sync.Mutex
}