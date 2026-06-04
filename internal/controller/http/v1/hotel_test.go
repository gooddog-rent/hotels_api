package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	infra "github.com/gooddog-rent/hotels_api/internal/infrastructure"
	svc "github.com/gooddog-rent/hotels_api/internal/service"
)

const URL = "http://localhost:4000/api/v1/hotels"

func TestReponseOK(t *testing.T) {

	// arrange
	searchRepoMock := infra.NewMockDB()
	searchServiceMock := svc.NewHotelService(searchRepoMock)
	hotelController := NewHotelController(searchServiceMock)

	query := "?query=bav&limit=4"
	ctx := context.WithValue(context.Background(), contextQueryKey, "bav")
	ctx = context.WithValue(ctx, contextLimitKey, 4)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, URL+query, nil)
	rr := httptest.NewRecorder()

	// act
	hotelController.GetHotels(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	// assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returns wrong status code: got %v want %v", status, http.StatusOK)
	}

	if req.Method != http.MethodGet {
		t.Errorf("Wrong request method: got %v want %v", req.Method, http.MethodGet)
	}

	output := ResponseLocations{
		Hotels: map[string][]string{
			"Bavaro":    {"Some Bavaro hotel", "Resort Bavaro Test"},
			"La Romana": {"Some La Romana bav hotel", "Resort La Romana bavaro Test"},
		},
	}

	outputJson, err := json.Marshal(output)
	if err != nil {
		t.Errorf("Cannot parse json struct! %v", err)
	}

	expected := string(outputJson)
	if rr.Body.String() != expected {
		t.Errorf("handler returns unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestInvalidParseJson(t *testing.T) {

	query := "?query=bav&limit=4"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	rr := httptest.NewRecorder()

	getHotelsHandler := func(_ http.ResponseWriter, req *http.Request) {

		x := map[string]interface{}{
			"foo": make(chan int),
		}

		// marshal json result
		_, err := json.Marshal(x)
		if err != nil {
			http.Error(rr, "500 Internal Server Error!", http.StatusInternalServerError)
			return
		}
	}

	handler := http.HandlerFunc(getHotelsHandler)
	handler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError || req.Method != http.MethodGet {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusInternalServerError, resp.StatusCode)
	}
}
