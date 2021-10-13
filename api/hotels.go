package api

import (
	"github.com/ljfuyuan/suffixtree"
)

// Location struct store Region and slice of hotels
type Location struct {
	Region string   `json:"region"`
	Hotels []string `json:"hotels"`
}

// Locations struct store slice of Locations and Counter struct
type Locations struct {
	Locations  []Location `json:"locations"`
	Tree       *suffixtree.GeneralizedSuffixTree
	HotelsList []string
	Result     []string
}

// FilteredHotels struct store output response hotels data
type FilteredHotels struct {
	Hotels map[string][]string `json:"hotels"`
}
