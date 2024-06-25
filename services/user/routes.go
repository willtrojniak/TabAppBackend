package user

import (
	"encoding/json"
	"net/http"
)


func (h *Handler) RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("GET /login", h.handleLogin);
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode("Hello");
  return;
}
