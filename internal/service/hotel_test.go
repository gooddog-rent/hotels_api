package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"hotels_api/internal/domain"
)

type MockSearcher struct {
	mock.Mock
}

func (m *MockSearcher) SearchHotels(ctx context.Context, text string, limit int) (*domain.Locations, error) {
	args := m.Called(ctx, text, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Locations), args.Error(1)
}

func TestSearchHotels(t *testing.T) {

	// Arrange
	mockSearcher := new(MockSearcher)
	mockSearcher.On("SearchHotels", mock.Anything, "test", 10).Return(&domain.Locations{}, nil)
	mockSearcher.On("SearchHotels", mock.Anything, "error", 10).Return(nil, errors.New("error"))

	service := NewHotelService(mockSearcher)

	ctx := context.Background()

	// Act
	// Test successful search
	response, err := service.SearchHotels(ctx, "test", 10)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Act
	// Test search with error
	response, err = service.SearchHotels(ctx, "error", 10)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
}
