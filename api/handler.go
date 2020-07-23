package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (l *Locations) Hotels(w http.ResponseWriter, req *http.Request) {

	t := time.Now()

	l.mu.Lock()
	l.counter++
	l.mu.Unlock()

	// prometheus api conuter metric increment
	l.Metrics.Counter.Inc()

	if req.Method != http.MethodGet {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.Host, http.StatusMethodNotAllowed, time.Since(t))
		http.Error(w, "Wrong method used! Only GET method allowed.", http.StatusMethodNotAllowed)
		return
	}

	// log.Printf("%s | API Requests: %d\n", req.Host, l.counter)

	params, err := url.ParseQuery(req.URL.Query().Encode())
	if err != nil {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.Host, http.StatusBadRequest, time.Since(t))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := params["query"]; !ok {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.Host, http.StatusBadRequest, time.Since(t))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := params["limit"]; !ok {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.Host, http.StatusBadRequest, time.Since(t))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	text := params.Get("query")
	if text == "" {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.Host, http.StatusBadRequest, time.Since(t))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil {
		log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.Host, http.StatusBadRequest, time.Since(t))
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filtered := struct {
		Hotels map[string][]string `json:"hotels"`
	}{}

	filtered.Hotels = make(map[string][]string)

	// words := getAllPermutations(strings.Split(text, " "))
	words := strings.Split(text, " ")
	limitMapLength := 0

	for _, location := range l.Locations {

		for _, hotel := range location.Hotels {

			ctr := 0
			for _, word := range words {

				//TODO: filter by regions in priority
				/* if strings.Contains(strings.ToLower(location.Region), strings.ToLower(word)) {
					ctr++
					// break
				} */

				// if strings.Count(strings.ToLower(location.Region), strings.ToLower(word)) == 1 && strings.Count(strings.ToLower(hotel), strings.ToLower(word)) > 0 {
				if strings.Contains(strings.ToLower(location.Region), strings.ToLower(word)) || strings.Count(strings.ToLower(hotel), strings.ToLower(word)) > 0 {
					// if strings.Count(strings.ToLower(hotel), strings.ToLower(word)) > 0 {
					ctr++
				}
			}

			if ctr > len(words)-1 && limitMapLength < limit {
				filtered.Hotels[location.Region] = append(filtered.Hotels[location.Region], hotel)
				limitMapLength++
			}
		}
	}

	out, err := json.Marshal(filtered)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // change to host domain
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// w.Header().Set("Content-Encoding", "gzip")
	// gz := gzip.NewWriter(w)
	// json.NewEncoder(gz).Encode(out)
	// gz.Close()

	log.Printf("%s %s - %s - %d - %v\n", req.Method, req.URL.String(), req.Host, http.StatusOK, time.Since(t))

	w.Write(out)
}

// var dublicates = map[string]bool{}

/* func getAllPermutations(query []string) []string {
	res := []string{}
	Perm(query, func(a []string) {
		// joined := strings.Join(a, " ")
		// res = append(res, joined)
		res = append(res, a...)
	})
	return res
}

func Perm(a []string, f func([]string)) {
	perm(a, f, 0)
}

func perm(a []string, f func([]string), i int) {
	if i > len(a) {
		f(a)
		return
	}

	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]

		// if ok, _ := dublicates[a[j]]; !ok {
		// dublicates[a[j]] = true
		perm(a, f, i+1)
		// }
		a[i], a[j] = a[j], a[i]
	}
} */
