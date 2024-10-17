package middleware

import (
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

func Middlewares(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func Bind(mid Middleware, route http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		mid(http.HandlerFunc(route)).ServeHTTP(w, req)
	}
}
