package location

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const URL = "http://localhost:4000/hotels"

func TestReponseOK(t *testing.T) {

	l := LocationService{}

	query := "?query=bav&limit=4"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	rr := httptest.NewRecorder()

	// build suffix search tree
	l.Build()

	l.GetHotels(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusOK || req.Method != http.MethodGet {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusOK, resp.StatusCode)
	}
}

func TestInvalidParseJson(t *testing.T) {

	query := "?query=bav&limit=4"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	rr := httptest.NewRecorder()

	getHotelsHandler := func(w http.ResponseWriter, req *http.Request) {

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

	if resp.StatusCode != http.StatusInternalServerError || req.Method != http.MethodGet {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusInternalServerError, resp.StatusCode)
	}
}
