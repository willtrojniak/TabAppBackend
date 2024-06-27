package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

type Middleware func(http.Handler) http.HandlerFunc

func WithMiddleware(middlewares ...Middleware) Middleware {
  return func(next http.Handler) http.HandlerFunc {
    for i := len(middlewares) - 1; i >= 0; i-- {
      next = middlewares[i](next)
    }
    return next.ServeHTTP;
  }
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    slog.Info("API Endpoint", "method", r.Method, "path", r.URL.Path);
    next.ServeHTTP(w, r);
  }
}

func RequireSession(s *sessions.SessionManager, h services.HTTPErrorHandler) Middleware {
  return func (next http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      _, err := s.GetSession(context.Background(), r);
      if err != nil {
        h(w, err);
        return;
      }
      next.ServeHTTP(w, r);
    }
  }
}
