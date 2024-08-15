package usecases

import (
	"hotels_api/cmd/app/internal/domain"
	"reflect"
	"testing"
)

type LocationUsecaseTest struct {
	domain.Locations
	HotelUsecase
}

func TestValidBuild(t *testing.T) {

	l := LocationUsecaseTest{}

	l.Build()

	if l.Tree == nil {
		t.Errorf("BuildSuffixTree must be valid tree! got: %v", l.Tree)
	}
}

func TestValidWithDataBuild(t *testing.T) {

	l := domain.Locations{
		Locations: []domain.Location{
			{
				Region: "Punta Cana",
				Hotels: []string{"Riu Bavaro", "Hotel Bavaro", "Resort Beach 5"},
			},
			{
				Region: "Macao",
				Hotels: []string{"Beach Summer 4", "Macao Beach Resort", "Star Hotel Resort"},
			},
		},
		// Tree:       nil,
		// HotelsList: nil,
		// Result:     nil,
	}

	loc := LocationUsecaseTest{
		Locations: l,
	}

	loc.Build()

	result := map[string][]string{
		"Punta Cana": {"Riu Bavaro", "Hotel Bavaro", "Resort Beach 5"},
		"Macao":      {"Beach Summer 4", "Macao Beach Resort", "Star Hotel Resort"},
	}

	// check valid regions in tree
	for _, location := range l.Locations {
		if _, ok := result[location.Region]; !ok {
			t.Errorf("BuildSuffixTree must have a valid regions! got: %v", location)
		}
	}

	// check valid hotels in tree
	for _, location := range l.Locations {

		if _, ok := result[location.Region]; ok {
			for index, h := range result[location.Region] {
				if h != location.Hotels[index] {
					t.Errorf("BuildSuffixTree must have a valid hotels! got: %v", h)
				}
			}
		}
	}

	if l.Tree == nil {
		t.Errorf("BuildSuffixTree must be valid tree! got: %v", l.Tree)
	}
}

func TestSuccessSearch(t *testing.T) {

	l := domain.Locations{
		Locations: []domain.Location{
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

	loc := LocationUsecaseTest{
		Locations: l,
	}

	text := "ea"
	limit := 10

	loc.Build()
	result := loc.Search(text, limit)

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
