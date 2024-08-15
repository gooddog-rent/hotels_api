package domain

import "hotels_api/pkg/suffixtree"

// Location struct store Region and slice of hotels
type Location struct {
	Region string   `json:"region"`
	Hotels []string `json:"hotels"`
}

// Locations struct store slice of Locations struct
type Locations struct {
	Locations  []Location `json:"locations"`
	Tree       *suffixtree.SuffixTree
	HotelsList []string
	Result     []string
}

// ResponseHotels struct store output response hotels data
type ResponseHotels struct {
	Hotels map[string][]string `json:"hotels"`
}
