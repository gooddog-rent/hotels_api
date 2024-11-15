package http

import (
	"context"
	"encoding/json"
	"hotels_api/internal/infra"
	"hotels_api/internal/service"
	"hotels_api/pkg/suffixtree"
	"net/http"
	"net/http/httptest"
	"testing"
)

const URL = "http://localhost:4000/hotels"

func TestReponseOK(t *testing.T) {

	// TODO: mock
	tree := suffixtree.NewGeneralizedSuffixTree()

	repositories := infra.NewRepositories(tree, "")

	deps := service.ServicesDependencies{
		Repos: repositories,
	}

	services := service.NewServices(deps)

	h := &hotelRoutes{
		hotelService: services.Hotel,
	}

	query := "?query=bav&limit=4"
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, URL+query, nil)
	rr := httptest.NewRecorder()

	h.getHotels(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK || req.Method != http.MethodGet {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusOK, resp.StatusCode)
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
