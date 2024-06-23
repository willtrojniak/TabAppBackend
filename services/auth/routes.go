package auth

import (
	"net/http"
)


func (h *Handler) RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("GET /auth/{provider}/callback", h.handleAuthCallback);
  router.HandleFunc("GET /auth/{provider}", h.handleAuth);
  router.HandleFunc("GET /logout", h.handleLogout);
}

func (h *Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
  provider := r.PathValue("provider");

  err := h.authorize(w, r, provider);
  if err != nil {
    h.handleError(w, err);
    return;
  }

  // http.Redirect(w, r, "http://localhost:5173", http.StatusFound);
  http.Redirect(w, r, "/api/v1/login", http.StatusFound);
  return;
  
}

func (h *Handler) handleAuth(w http.ResponseWriter, r *http.Request) {
  provider := r.PathValue("provider");
  if err := h.beginAuthorize(w, r, provider); err != nil {
    h.handleError(w, err);
    return;
  }
  
  return;
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {

  if err := h.logout(w, r); err != nil {
    h.handleError(w, err);
    return;
  }

  http.Redirect(w, r, "/api/v1/login", http.StatusFound);
  return;
}
