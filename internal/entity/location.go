package entity

// Location struct store Region and list of hotels
type Location struct {
	Region string   `json:"region"`
	Hotels []string `json:"hotels"`
}

// Locations struct store slice of Locations struct
type LocationsStore struct {
	Locations []Location `json:"locations"`
}

// ResponseHotels struct store output response hotels data
type ResponseHotels struct {
	Hotels map[string][]string `json:"hotels"`
}
