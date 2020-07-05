package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

const url = "http://localhost:4000/hotels"

func TestReponseOK(t *testing.T) {

	l := Locations{
		Counter: Counter{
			Metrics: Metrics{
				Counter: NewMetrics().Counter,
			},
		},
	}

	buf := bytes.NewBufferString("?query=riu&limit=4")
	req := httptest.NewRequest(http.MethodGet, url, buf)
	w := httptest.NewRecorder()

	l.Hotels(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK || req.Method != http.MethodGet {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusOK, resp.StatusCode)
	}
}
