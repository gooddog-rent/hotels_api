package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (l *Locations) Hotels(w http.ResponseWriter, req *http.Request) {

	l.mu.Lock()
	l.counter++
	l.mu.Unlock()

	l.Metrics.Counter.Inc()

	if req.Method != http.MethodGet {
		http.Error(w, "Wrong method used! Only GET method allowed.", http.StatusMethodNotAllowed)
	}

	log.Printf("%s | API Requests: %d\n", req.Host, l.counter)

	text := req.URL.Query().Get("query")
	limit, err := strconv.Atoi(req.URL.Query().Get("limit"))
	if err != nil {
		log.Println(err)
	}

	filtered := struct {
		Hotels map[string][]string `json:"hotels"`
	}{}

	filtered.Hotels = make(map[string][]string)

	for _, location := range l.Locations {
		for _, hotel := range location.Hotels {
			if (strings.Contains(strings.ToLower(hotel), strings.ToLower(text)) || strings.Contains(strings.ToLower(location.Region), strings.ToLower(text))) && len(filtered.Hotels[location.Region]) < limit/2 && len(filtered.Hotels) < limit {
				filtered.Hotels[location.Region] = append(filtered.Hotels[location.Region], hotel)
			}
		}
	}

	out, err := json.Marshal(filtered)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://192.168.1.108:5500")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Write(out)
}
