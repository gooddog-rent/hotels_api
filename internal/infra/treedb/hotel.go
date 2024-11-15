package treedb

import (
	"hotels_api/internal/entity"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Search method returns filtered request with hotels
func (h *WatcherRepo) SearchHotels(text string, limit int) *entity.ResponseHotels {

	filtered := &entity.ResponseHotels{
		Hotels: make(map[string][]string),
	}

	h.mu.RLock()
	indexes := h.tree.searchTree.Tree.Search(strings.ToLower(text), -1)
	h.mu.RUnlock()

	for _, index := range indexes {

		regionAndHotel := strings.Split(cases.Title(language.Und, cases.NoLower).String(h.tree.searchTree.HotelsList[index]), ";")
		region := regionAndHotel[0]
		hotel := regionAndHotel[1]

		// limit hotels per one region
		if len(filtered.Hotels[region]) < limit {
			h.mu.Lock()
			filtered.Hotels[region] = append(filtered.Hotels[region], hotel)
			h.mu.Unlock()
		}
	}

	return filtered
}
