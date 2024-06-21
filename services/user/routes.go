package user

import (
	"encoding/json"
	"net/http"
)

type Handler struct{
}

type Test struct {
  Name string `json:"name"`
}

func NewHandler() *Handler {
  return &Handler{};
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("GET /login", h.handleLogin);
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(Test{Name: "Will"});
  return;
}
