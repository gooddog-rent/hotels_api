package api

import (
	"strings"

	"hotels_api/suffixtree"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Location struct store Region and slice of hotels
type Location struct {
	Region string   `json:"region"`
	Hotels []string `json:"hotels"`
}

// Locations struct store slice of Locations and Counter struct
type Locations struct {
	Locations  []Location `json:"locations"`
	Tree       *suffixtree.SuffixTree
	HotelsList []string
	Result     []string
}

// FilteredHotels struct store output response hotels data
type FilteredHotels struct {
	Hotels map[string][]string `json:"hotels"`
}

// BuildSuffixTree method create a new search tree with regions and hotels
func (l *Locations) BuildSuffixTree() {

	hotels := []string{}

	// build hotels list
	for _, locs := range l.Locations {
		for _, hotel := range locs.Hotels {
			hotels = append(hotels, strings.ToLower(locs.Region)+";"+strings.ToLower(hotel))
		}
	}

	// build new suffix tree
	tree := suffixtree.NewGeneralizedSuffixTree()
	for idx, hotel := range hotels {
		tree.Put(hotel, idx)
	}

	l.Tree = tree
	l.HotelsList = hotels
}

// SearchInSuffixTree method returns filtered request with hotels
func (l *Locations) SearchInSuffixTree(text string, limit int) FilteredHotels {

	filtered := FilteredHotels{
		make(map[string][]string),
	}

	indexes := l.Tree.Search(strings.ToLower(text), -1)

	for _, index := range indexes {

		regionAndHotel := strings.Split(cases.Title(language.Und, cases.NoLower).String(l.HotelsList[index]), ";")
		region := regionAndHotel[0]
		hotel := regionAndHotel[1]

		// limit hotels per one region
		if len(filtered.Hotels[region]) < limit {
			filtered.Hotels[region] = append(filtered.Hotels[region], hotel)
		}
	}

	return filtered
}
