package infrastructure

import (
	"hotels_api/internal/entity"
	"hotels_api/pkg/suffixtree"
	"strings"
	"sync"
)

// static interface implementation check for convinience
var _ Tree = (*TreeRepo)(nil)

type TreeRepo struct {
	searchTree *entity.SearchTree
	mu         sync.Mutex
}

// BuildTree method create a new search tree with regions and hotels
func (t *TreeRepo) BuildTree(l *entity.LocationsStore) *entity.SearchTree {

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
		t.mu.Lock()
		tree.Put(hotel, idx)
		t.mu.Unlock()
	}

	return &entity.SearchTree{
		Tree:       tree,
		HotelsList: hotels,
	}
}
