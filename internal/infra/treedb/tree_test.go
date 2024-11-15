package treedb

import (
	"hotels_api/internal/entity"
	"hotels_api/pkg/suffixtree"
	"testing"
)

func TestValidBuild(t *testing.T) {

	l := &entity.LocationsStore{}
	st := &entity.SearchTree{
		Tree: suffixtree.NewGeneralizedSuffixTree(),
	}

	treeRepo := &TreeRepo{
		searchTree: st,
	}
	searchTree := treeRepo.BuildTree(l)

	if searchTree.Tree == nil {
		t.Errorf("BuildTree must be valid tree! got: %v", searchTree.Tree)
	}
}

func TestValidWithDataBuild(t *testing.T) {

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

	loc := &entity.LocationsStore{
		Locations: l.Locations,
	}

	st := &entity.SearchTree{
		Tree: suffixtree.NewGeneralizedSuffixTree(),
	}

	treeRepo := &TreeRepo{
		searchTree: st,
	}
	searchTree := treeRepo.BuildTree(loc)

	result := map[string][]string{
		"Punta Cana": {"Riu Bavaro", "Hotel Bavaro", "Resort Beach 5"},
		"Macao":      {"Beach Summer 4", "Macao Beach Resort", "Star Hotel Resort"},
	}

	// check valid regions in tree
	for _, location := range l.Locations {
		if _, ok := result[location.Region]; !ok {
			t.Errorf("BuildTree must have a valid regions! got: %v", location)
		}
	}

	// check valid hotels in tree
	for _, location := range l.Locations {

		if _, ok := result[location.Region]; ok {
			for index, h := range result[location.Region] {
				if h != location.Hotels[index] {
					t.Errorf("BuildTree must have a valid hotels! got: %v", h)
				}
			}
		}
	}

	if searchTree.Tree == nil {
		t.Errorf("BuildTree must be valid tree! got: %v", searchTree.Tree)
	}
}
