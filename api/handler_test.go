package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const URL = "http://localhost:4000/hotels"

func TestReponseOK(t *testing.T) {

	l := Locations{}

	query := "?query=bav&limit=4"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	w := httptest.NewRecorder()

	// build suffix search tree
	l.BuildSuffixTree()

	l.GetHotels(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK || req.Method != http.MethodGet {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusOK, resp.StatusCode)
	}
}

func TestReponseBadRequest(t *testing.T) {

	l := Locations{}

	query := "?query=&limit=4" //TODO: add more cases
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	w := httptest.NewRecorder()

	l.GetHotels(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest || req.Method != http.MethodGet {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestBadMethod(t *testing.T) {

	l := Locations{}

	query := "?query=riu&limit=4"
	req := httptest.NewRequest(http.MethodPost, URL+query, nil)
	w := httptest.NewRecorder()

	l.GetHotels(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed || req.Method != http.MethodPost {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestInvalidQuery(t *testing.T) {

	l := Locations{}

	query := "?q"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	w := httptest.NewRecorder()

	l.GetHotels(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestInvalidLimit(t *testing.T) {

	l := Locations{}

	query := "?query=bav&li"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	w := httptest.NewRecorder()

	l.GetHotels(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestLimitQueryNotNumber(t *testing.T) {

	l := Locations{}

	query := "?query=bav&limit=four"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	w := httptest.NewRecorder()

	l.GetHotels(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusBadRequest, resp.StatusCode)
	}
}
