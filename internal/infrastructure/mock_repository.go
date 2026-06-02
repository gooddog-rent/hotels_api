package infrastructure

import (
	"context"
	"hotels_api/internal/domain"
)

// static check interface implementation
var _ Searcher = (*MockDB)(nil)

// mock repository struct
type MockDB struct{}

// mock implementation without db, only mock title
func (_ *MockDB) SearchHotels(_ context.Context, text string, limit int) (*domain.Locations, error) {
	return &domain.Locations{
		Hotels: map[string][]string{

			// mocked test data. simulated db response
			"Bavaro":    {"Some Bavaro hotel", "Resort Bavaro Test"},
			"La Romana": {"Some La Romana bav hotel", "Resort La Romana bavaro Test"},
		},
	}, nil
}

func NewMockDB() *MockDB {
	return &MockDB{}
}
