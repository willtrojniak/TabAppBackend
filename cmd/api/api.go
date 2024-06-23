package api

import (
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/services/auth"
	"github.com/WilliamTrojniak/TabAppBackend/services/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
  addr string
  pool *pgxpool.Pool
}

type Handler interface {
  RegisterRoutes(http.ServeMux);
}

func NewAPIServer(addr string, pool *pgxpool.Pool) *APIServer {
  return &APIServer{
    addr: addr,
    pool: pool,
  };
}

func (s *APIServer) Run() error {

  authHandler := auth.NewHandler(user.NewPGXStore(s.pool).CreateUser);
  userHandler := user.NewHandler(user.NewPGXStore(s.pool), authHandler);

  router := http.NewServeMux()
  v1 := http.NewServeMux();
  
  authHandler.RegisterRoutes(router);
  userHandler.RegisterRoutes(v1);

  router.Handle("/api/v1/", http.StripPrefix("/api/v1", authHandler.RequireAuth(v1)));
  return http.ListenAndServe(s.addr, WithMiddleware(RequestLoggerMiddleware)(router));
}
