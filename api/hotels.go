package api

type Location struct {
	Region string   `json:"region"`
	Hotels []string `json:"hotels"`
}

type Locations struct {
	Locations []Location `json:"locations"`
	Counter
}
