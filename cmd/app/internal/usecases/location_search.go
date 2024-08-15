package usecases

import (
	"hotels_api/cmd/app/internal/domain"
	"hotels_api/pkg/suffixtree"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// static interface implementation check for convinience
var _ HotelUsecase = (*LocationUsecase)(nil)

type LocationUsecase struct {
	domain.Locations
}

// Build method create a new search tree with regions and hotels
func (l LocationUsecase) Build() {

	hotels := []string{}

	// build hotels list
	for _, locs := range l.Locations.Locations {
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

// Search method returns filtered request with hotels
func (l LocationUsecase) Search(text string, limit int) domain.ResponseHotels {

	filtered := domain.ResponseHotels{
		Hotels: make(map[string][]string),
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
