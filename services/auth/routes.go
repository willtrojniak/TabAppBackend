package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/markbates/goth/gothic"
)


func (h *Handler) RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("GET /auth/{provider}/callback", h.handleAuthCallback);
  router.HandleFunc("GET /auth/{provider}", h.handleAuth);
  router.HandleFunc("GET /logout/{provider}", h.handleLogout);
}

func (h *Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
  provider := r.PathValue("provider");
  r = r.WithContext(context.WithValue(context.Background(), "provider", provider));

  user, err := gothic.CompleteUserAuth(w, r);
  if err != nil {
    fmt.Fprintln(w, err);
    return;
  }

  userUUID, err := h.userStore.CreateUser(context.TODO(), &types.UserCreate{Email: user.Email, Name: user.Name});
  if err != nil {
    // TODO: Better logging
    fmt.Printf("%v\n", err);
    http.Error(w, "Internal server error", http.StatusInternalServerError);
    return;
  }
  
  if err := h.storeUserUUID(w, r, userUUID); err != nil {
    // TODO: Better logging
    fmt.Printf("%v\n", err);
    http.Error(w, "Internal server error", http.StatusInternalServerError);
    return;
  }
  

  fmt.Println("=====NAME=====")
  // fmt.Println(user.Name);
  fmt.Println("=====EMAIL=====")
  // fmt.Println(user.Email);
  fmt.Println("=====UUID=====")
  fmt.Println(userUUID.String());
  // http.Redirect(w, r, "http://localhost:5173", http.StatusFound);
  http.Redirect(w, r, "/api/v1/login", http.StatusFound);
  return;
  
}

func (h *Handler) handleAuth(w http.ResponseWriter, r *http.Request) {
  provider := r.PathValue("provider");
  r = r.WithContext(context.WithValue(context.Background(), "provider", provider));
  gothic.BeginAuthHandler(w, r);
  return;
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
  provider := r.PathValue("provider");
  r = r.WithContext(context.WithValue(context.Background(), "provider", provider));

  if err := h.logout(w, r); err != nil {
    http.Error(w, "Internal service error.", http.StatusInternalServerError);
    return;
  }
  gothic.Logout(w, r);
  http.Redirect(w, r, "/api/v1/login", http.StatusFound);
  return;
}
