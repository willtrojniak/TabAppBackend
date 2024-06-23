package user

import (
	"encoding/json"
	"net/http"
)


func (h *Handler) RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("GET /login", h.handleLogin);
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
  id, err := h.auth.GetUserSession(r);
  if err != nil {
    h.handleError(w, err);
    return;
  }

  json.NewEncoder(w).Encode(id);
  return;
}
