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

func CORSMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	}
}
