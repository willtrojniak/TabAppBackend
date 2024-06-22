package api

import (
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/services/auth"
	"github.com/WilliamTrojniak/TabAppBackend/services/user"
)

type APIServer struct {
  addr string
}

type Handler interface {
  RegisterRoutes(http.ServeMux);
}

func NewAPIServer(addr string) *APIServer {
  return &APIServer{
    addr: addr,
  };
}

func (s *APIServer) Run() error {
  router := http.NewServeMux()

  authHandler := auth.NewHandler();
  authHandler.RegisterRoutes(router);

  v1 := http.NewServeMux();
  
  userHandler := user.NewHandler(authHandler);
  userHandler.RegisterRoutes(v1);


  router.Handle("/api/v1/", http.StripPrefix("/api/v1", authHandler.RequireAuth(v1)));
  return http.ListenAndServe(s.addr, WithMiddleware(RequestLoggerMiddleware)(router));
}
