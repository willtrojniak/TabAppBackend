package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
)


func (h *Handler) RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("GET /auth/{provider}/callback", h.handleAuthCallback);
  router.HandleFunc("GET /auth/{provider}", h.handleAuth);
}

func (h *Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
  provider := r.PathValue("provider");
  r = r.WithContext(context.WithValue(context.Background(), "provider", provider));

  user, err := gothic.CompleteUserAuth(w, r);
  if err != nil {
    fmt.Fprintln(w, err);
    return
  }
  
  err = h.storeUserSession(w, r, user);
  if err != nil {
    http.Error(w, "Internal server error", http.StatusInternalServerError);
  }

  fmt.Println(user);
  http.Redirect(w, r, "http://localhost:5173", http.StatusFound);
  return;
  
}

func (h *Handler) handleAuth(w http.ResponseWriter, r *http.Request) {
  provider := r.PathValue("provider");
  r = r.WithContext(context.WithValue(context.Background(), "provider", provider));

  if user, err := gothic.CompleteUserAuth(w, r); err == nil {
    fmt.Println(user);
  } else {
    gothic.BeginAuthHandler(w, r);
  }
}
