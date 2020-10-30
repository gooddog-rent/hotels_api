package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestShowLogMiddleware(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	showLogHandler := func(w http.ResponseWriter, req *http.Request) {
		tn := time.Now()
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.Host, http.StatusBadRequest, time.Since(tn))
	}

	rr := httptest.NewRecorder()

	handler := ShowLog(http.HandlerFunc(showLogHandler))
	handler(rr, req)

	if buf.Len() == 0 {
		t.Error("Empty log output!")
	}
}
