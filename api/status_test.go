package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteHeaderValid(t *testing.T) {

	sr := StatusHTTP{
		httptest.NewRecorder(),
		http.StatusOK,
	}

	testStatusCode := 400
	sr.WriteHeader(testStatusCode)

	if sr.StatusCode != testStatusCode {
		t.Errorf("Wrong Status Code in header!: got %d", sr.StatusCode)
	}
}
