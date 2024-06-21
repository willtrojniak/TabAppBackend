package api

import (
	"net/http"

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
  v1 := http.NewServeMux();
  
  userHandler := user.NewHandler();
  userHandler.RegisterRoutes(v1);


  router := http.NewServeMux()
  router.Handle("/api/v1/", http.StripPrefix("/api/v1", WithMiddleware(RequestLoggerMiddleware)(v1)));
  return http.ListenAndServe(s.addr, router);
}
