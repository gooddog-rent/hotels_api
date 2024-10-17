package entity

import (
	"hotels_api/pkg/suffixtree"
)

type SearchTree struct {
	Tree       *suffixtree.SuffixTree
	HotelsList []string
}
