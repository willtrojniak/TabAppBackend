package api

import (
	"log/slog"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/auth"
	"github.com/WilliamTrojniak/TabAppBackend/services/user"
)

type APIServer struct {
  addr string
  store *db.PgxStore
}

type Handler interface {
  RegisterRoutes(http.ServeMux);
}

func NewAPIServer(addr string, store *db.PgxStore) *APIServer {
  return &APIServer{
    addr: addr,
    store: store,
  };
}

func (s *APIServer) Run() error {

  authHandler := auth.NewHandler(services.HandleHttpError, slog.Default());
  userHandler := user.NewHandler(s.store, authHandler, services.HandleHttpError, slog.Default());
  authHandler.SetCreateUserFn(userHandler.CreateUser);

  router := http.NewServeMux()
  v1 := http.NewServeMux();
  
  authHandler.RegisterRoutes(router);
  userHandler.RegisterRoutes(v1);

  router.Handle("/api/v1/", http.StripPrefix("/api/v1", authHandler.RequireAuth(v1)));
  return http.ListenAndServe(s.addr, WithMiddleware(RequestLoggerMiddleware)(router));
}
