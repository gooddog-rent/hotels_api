package api

import "strings"

// Location struct store Region and slice of hotels
type Location struct {
	Region string   `json:"region"`
	Hotels []string `json:"hotels"`
}

// Locations struct store slice of Locations and Counter struct
type Locations struct {
	Locations []Location `json:"locations"`
	Counter
}

// FilteredHotels struct store output response hotels data
type FilteredHotels struct {
	Hotels map[string][]string `json:"hotels"`
}

// filterByHotels method returns filtered request with hotels
func (l *Locations) filterByHotels(text string, limit int) FilteredHotels {

	filtered := FilteredHotels{
		make(map[string][]string),
	}

	words := strings.Split(text, " ")
	limitMapLength := 0

	for _, location := range l.Locations {
		for _, hotel := range location.Hotels {

			ctr := 0
			for _, word := range words {

				if strings.Count(strings.ToLower(hotel), strings.ToLower(word)) > 0 {
					ctr++
				}

				if ctr > len(words)-1 && limitMapLength < limit {
					filtered.Hotels[location.Region] = append(filtered.Hotels[location.Region], hotel)
					limitMapLength++
				}
			}
		}
	}
	return filtered
}
