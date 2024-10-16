package api

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/WilliamTrojniak/TabAppBackend/cache"
	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/auth"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/services/shop"
	"github.com/WilliamTrojniak/TabAppBackend/services/user"
	"github.com/redis/go-redis/v9"
)

type APIServer struct {
	addr  string
	store *db.PgxStore
	cache *redis.Client
}

type Handler interface {
	RegisterRoutes(http.ServeMux)
}

func NewAPIServer(addr string, store *db.PgxStore, cache *redis.Client) *APIServer {
	return &APIServer{
		addr:  addr,
		store: store,
		cache: cache,
	}
}

func (s *APIServer) Run() error {
	sessionStore := cache.NewRedisCache(s.cache)
	sessionManager := sessions.New(sessionStore, time.Hour*24*30, time.Hour*1, services.HandleHttpError, slog.Default())

	authHandler, err := auth.NewHandler(services.HandleHttpError, sessionManager, slog.Default())
	if err != nil {
		log.Fatal("Failed to initialize auth handler")
	}
	userHandler := user.NewHandler(s.store, sessionManager, services.HandleHttpError, slog.Default())
	authHandler.SetCreateUserFn(userHandler.CreateUser)

	shopHandler := shop.NewHandler(s.store, sessionManager, services.HandleHttpError, slog.Default())

	router := http.NewServeMux()
	v1 := http.NewServeMux()

	authHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(v1)
	shopHandler.RegisterRoutes(v1)

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", WithMiddleware(
		sessionManager.RequireAuth)(v1)))

	return http.ListenAndServe(s.addr, WithMiddleware(RequestLoggerMiddleware, CORSMiddleware, sessionManager.RequireCSRFToken)(router))
}
