package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const URL = "http://localhost:4000/hotels"

func TestReponseOK(t *testing.T) {

	query := "?query=bav&limit=4"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	rr := httptest.NewRecorder()

	validateRequestHandler := func(w http.ResponseWriter, req *http.Request) {

		// check on correct request method
		if req.Method != http.MethodGet {
			http.Error(rr, "Wrong method used! Only GET method allowed.", http.StatusMethodNotAllowed)
			return
		}

		// parse encoded query params to map
		params, err := url.ParseQuery(req.URL.Query().Encode())
		if err != nil {
			http.Error(rr, "Invalid query parameters!", http.StatusBadRequest)
			return
		}

		// check if "query" word is exists
		if _, ok := params["query"]; !ok {
			http.Error(rr, "'query' parameter does not exist or bad value!", http.StatusBadRequest)
			return
		}

		// check if "limit" word is exists
		if _, ok := params["limit"]; !ok {
			http.Error(rr, "'limit' parameter does not exist or bad value!", http.StatusBadRequest)
			return
		}

		// parse search query
		text := params.Get("query")
		if text == "" {
			http.Error(rr, "'query' value is empty!", http.StatusBadRequest)
			return
		}

		// convert limit value to string
		_, err = strconv.Atoi(params.Get("limit"))
		if err != nil {
			log.Println(err)
			http.Error(rr, "'limit' value must be integer!", http.StatusBadRequest)
			return
		}
	}

	handler := ValidateRequest(http.HandlerFunc(validateRequestHandler))
	handler(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusOK || req.Method != http.MethodGet {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusOK, resp.StatusCode)
	}
}

func TestReponseBadRequest(t *testing.T) {

	query := "?query=&limit=4" //TODO: add more cases
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	rr := httptest.NewRecorder()

	validateRequestHandler := func(w http.ResponseWriter, req *http.Request) {

		// parse encoded query params to map
		params, _ := url.ParseQuery(req.URL.Query().Encode())

		// check if "query" word is exists
		if _, ok := params["query"]; !ok {
			http.Error(rr, "'query' parameter does not exist or bad value!", http.StatusBadRequest)
			return
		}
	}

	handler := ValidateRequest(http.HandlerFunc(validateRequestHandler))
	handler(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusBadRequest || req.Method != http.MethodGet {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestBadMethod(t *testing.T) {

	query := "?query=riu&limit=4"
	req := httptest.NewRequest(http.MethodPost, URL+query, nil)
	rr := httptest.NewRecorder()

	validateRequestHandler := func(w http.ResponseWriter, req *http.Request) {

		// check on correct request method
		if req.Method != http.MethodGet {
			http.Error(rr, "Wrong method used! Only GET method allowed.", http.StatusMethodNotAllowed)
			return
		}
	}

	handler := ValidateRequest(http.HandlerFunc(validateRequestHandler))
	handler(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed || req.Method != http.MethodPost {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestInvalidQuery(t *testing.T) {

	query := "?q"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	rr := httptest.NewRecorder()

	validateRequestHandler := func(w http.ResponseWriter, req *http.Request) {

		// parse encoded query params to map
		_, err := url.ParseQuery(req.URL.Query().Encode())
		if err != nil {
			http.Error(rr, "Invalid query parameters!", http.StatusBadRequest)
			return
		}
	}

	handler := ValidateRequest(http.HandlerFunc(validateRequestHandler))
	handler(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestInvalidLimit(t *testing.T) {

	query := "?query=bav&li"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	rr := httptest.NewRecorder()

	validateRequestHandler := func(w http.ResponseWriter, req *http.Request) {

		// parse encoded query params to map
		params, _ := url.ParseQuery(req.URL.Query().Encode())

		// check if "limit" word is exists
		if _, ok := params["limit"]; !ok {
			http.Error(rr, "'limit' parameter does not exist or bad value!", http.StatusBadRequest)
			return
		}
	}

	handler := ValidateRequest(http.HandlerFunc(validateRequestHandler))
	handler(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestLimitQueryNotNumber(t *testing.T) {

	query := "?query=bav&limit=four"
	req := httptest.NewRequest(http.MethodGet, URL+query, nil)
	rr := httptest.NewRecorder()

	validateRequestHandler := func(w http.ResponseWriter, req *http.Request) {

		// parse encoded query params to map
		params, _ := url.ParseQuery(req.URL.Query().Encode())

		// convert limit value to string
		_, err := strconv.Atoi(params.Get("limit"))
		if err != nil {
			http.Error(rr, "'limit' value must be integer!", http.StatusBadRequest)
			return
		}
	}

	handler := ValidateRequest(http.HandlerFunc(validateRequestHandler))
	handler(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Bad response status code! Excpect: %d Have: %d", http.StatusBadRequest, resp.StatusCode)
	}
}
