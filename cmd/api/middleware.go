package api

import (
	"log/slog"
	"net/http"
)

type Middleware func(http.Handler) http.HandlerFunc

func WithMiddleware(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("API Endpoint", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
