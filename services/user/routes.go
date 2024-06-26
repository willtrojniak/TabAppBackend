package user

import (
	"encoding/json"
	"net/http"
)


func (h *Handler) RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("GET /login", h.handleLogin);
  h.logger.Info("Registered user routes");
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode("Hello");
  return;
}
