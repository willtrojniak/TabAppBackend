package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/services/auth"
)

type Handler struct{
  authHandler *auth.Handler
  
}

func NewHandler(authHandler *auth.Handler) *Handler {
  return &Handler{
    authHandler: authHandler,
  };
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("GET /login", h.handleLogin);
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
  user, err := h.authHandler.GetUserSession(r);
  if err != nil {
    http.Error(w, "unauthorized", http.StatusUnauthorized);
    return;
  }
  fmt.Println(user);

  json.NewEncoder(w).Encode(user);
  return;
}
