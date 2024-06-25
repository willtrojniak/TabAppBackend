package api

import (
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

  userHandler := user.NewHandler(s.store, services.HandleHttpError);
  authHandler := auth.NewHandler(userHandler, services.HandleHttpError);

  router := http.NewServeMux()
  v1 := http.NewServeMux();
  
  authHandler.RegisterRoutes(router);
  userHandler.RegisterRoutes(v1);

  router.Handle("/api/v1/", http.StripPrefix("/api/v1", authHandler.RequireAuth(v1)));
  return http.ListenAndServe(s.addr, WithMiddleware(RequestLoggerMiddleware)(router));
}
