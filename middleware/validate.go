package middleware

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// ValidateRequest middleware handler check and validate user requset
func ValidateRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// check on correct request method
		if req.Method != http.MethodGet {
			http.Error(w, "Wrong method used! Only GET method allowed.", http.StatusMethodNotAllowed)
			return
		}

		// parse encoded query params to map
		params, err := url.ParseQuery(req.URL.Query().Encode())
		if err != nil {
			http.Error(w, "Invalid query parameters!", http.StatusBadRequest)
			return
		}

		// check if "query" word is exists
		if _, ok := params["query"]; !ok {
			http.Error(w, "'query' parameter does not exist or bad value!", http.StatusBadRequest)
			return
		}

		// check if "limit" word is exists
		if _, ok := params["limit"]; !ok {
			http.Error(w, "'limit' parameter does not exist or bad value!", http.StatusBadRequest)
			return
		}

		// parse search query
		text := params.Get("query")
		if text == "" {
			http.Error(w, "'query' value is empty!", http.StatusBadRequest)
			return
		}

		// convert limit value to string
		_, err = strconv.Atoi(params.Get("limit"))
		if err != nil {
			log.Println(err)
			http.Error(w, "'limit' value must be integer!", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, req)
	})
}
