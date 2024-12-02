package http

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// ValidateRequest middleware handler check and validate user requset
func ValidateRequest(next http.Handler) http.Handler {
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
		if _, ok := params[string(contextQueryKey)]; !ok {
			http.Error(w, "'query' parameter does not exist or bad value!", http.StatusBadRequest)
			return
		}

		// check if "limit" word is exists
		if _, ok := params[string(contextLimitKey)]; !ok {
			http.Error(w, "'limit' parameter does not exist or bad value!", http.StatusBadRequest)
			return
		}

		// parse search query
		text := params.Get(string(contextQueryKey))
		if text == "" {
			http.Error(w, "'query' value is empty!", http.StatusBadRequest)
			return
		}

		// convert limit value to string
		limit, err := strconv.Atoi(params.Get(string(contextLimitKey)))
		if err != nil {
			log.Println(err)
			http.Error(w, "'limit' value must be integer!", http.StatusBadRequest)
			return
		}

		// add valid query and limit values to context
		ctx := context.WithValue(req.Context(), contextQueryKey, text)
		ctx = context.WithValue(ctx, contextLimitKey, limit)

		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
