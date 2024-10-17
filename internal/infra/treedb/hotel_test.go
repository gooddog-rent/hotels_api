package treedb

import (
	"hotels_api/internal/entity"
	"hotels_api/pkg/suffixtree"
	"reflect"
	"testing"
)

func TestSuccessSearch(t *testing.T) {

	l := &entity.LocationsStore{
		Locations: []entity.Location{
			{
				Region: "Punta Cana",
				Hotels: []string{"Riu Bavaro", "Hotel Bavaro", "Resort Beach 5"},
			},
			{
				Region: "Macao",
				Hotels: []string{"Beach Summer 4", "Macao Beach Resort", "Star Hotel Resort"},
			},
		},
	}

	treeRepo := &TreeRepo{
		searchTree: &entity.SearchTree{
			Tree:       suffixtree.NewGeneralizedSuffixTree(),
			HotelsList: []string{},
		},
	}

	text := "ea"
	limit := 10

	treeRepo.searchTree = treeRepo.BuildTree(l)

	loc := &WatcherRepo{
		tree: treeRepo,
	}

	result := loc.SearchHotels(text, limit)

	type ResultHotels struct {
		Hotels map[string][]string
	}

	validResult := ResultHotels{
		Hotels: map[string][]string{
			"Punta Cana": {"Resort Beach 5"},
			"Macao":      {"Macao Beach Resort", "Beach Summer 4"},
		},
	}

	if !reflect.DeepEqual(result.Hotels, validResult.Hotels) {
		t.Errorf("Search results must have a valid regions and hotels! got: %v", result)
	}
}
