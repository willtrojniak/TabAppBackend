package api

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/willtrojniak/TabAppBackend/cache"
	"github.com/willtrojniak/TabAppBackend/db"
	"github.com/willtrojniak/TabAppBackend/services"
	"github.com/willtrojniak/TabAppBackend/services/auth"
	"github.com/willtrojniak/TabAppBackend/services/events"
	"github.com/willtrojniak/TabAppBackend/services/sessions"
	"github.com/willtrojniak/TabAppBackend/services/shop"
	"github.com/willtrojniak/TabAppBackend/services/user"
)

type APIServer struct {
	addr   string
	store  *db.PgxStore
	cache  *redis.Client
	events *events.EventDispatcher
}

type Handler interface {
	RegisterRoutes(http.ServeMux)
}

func NewAPIServer(
	addr string,
	store *db.PgxStore,
	cache *redis.Client,
	events *events.EventDispatcher) *APIServer {
	return &APIServer{
		addr:   addr,
		store:  store,
		cache:  cache,
		events: events,
	}
}

func (s *APIServer) Run() error {

	sessionStore := cache.NewRedisCache(s.cache)
	sessionManager := sessions.New(sessionStore, time.Hour*24*30, time.Hour*1, services.HandleHttpError, slog.Default())
	userHandler := user.NewHandler(s.store, sessionManager, services.HandleHttpError, slog.Default())

	authHandler, err := auth.NewHandler(services.HandleHttpError, sessionManager, userHandler, slog.Default())
	if err != nil {
		log.Fatal("Failed to initialize auth handler")
	}

	shopHandler := shop.NewHandler(s.store, sessionManager, s.events, services.HandleHttpError, slog.Default())

	router := http.NewServeMux()
	v1 := http.NewServeMux()

	authHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(v1)
	shopHandler.RegisterRoutes(v1)

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", WithMiddleware(
		sessionManager.RequireAuth)(v1)))

	return http.ListenAndServe(s.addr, WithMiddleware(RequestLoggerMiddleware, CORSMiddleware, sessionManager.RequireCSRFToken)(router))
}
