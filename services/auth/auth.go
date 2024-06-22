package auth

import (
	"fmt"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/env"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
  maxAge = 86400 * 30 // 30 days
  session_cookie = "user_session"
)

type Handler struct{
  store sessions.Store
}

func NewHandler() *Handler {
  return &Handler{
    store: gothic.Store, 
  };
}

func init() {
  config := env.GetConfig();

  store := sessions.NewCookieStore([]byte(config.SESSION_SECRET));
  store.MaxAge(maxAge);
  store.Options.Path = "/"
  store.Options.HttpOnly = true;
  store.Options.Secure = true; // BREAKS IF FALSE

  gothic.Store = store;

  goth.UseProviders(
    // TODO: Make URL dynamic
    google.New(config.OAUTH2_GOOGLE_CLIENT_ID, config.OAUTH2_GOOGLE_CLIENT_SECRET, "http://127.0.0.1:3000/auth/google/callback"),
  );
}

func (h *Handler) storeUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {

  session, _ := h.store.Get(r, session_cookie);
  session.Values["user"] = user;
  
  err := session.Save(r, w);
  return err;
}

func (h *Handler) GetUserSession(r *http.Request) (goth.User, error) {
  session, err := h.store.Get(r, session_cookie);
  if err != nil {
    return goth.User{}, err;
  }

  u := session.Values["user"];
  if u == nil {
    return goth.User{}, fmt.Errorf("user is not authenticated! %v", u);
  }

  return u.(goth.User), nil;
}

func (h *Handler) RequireAuth(next http.Handler) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    _, err := h.GetUserSession(r);
    if err != nil {
      http.Error(w, "unauthorized", http.StatusUnauthorized);
      return;
    }
    next.ServeHTTP(w, r);
  }
}
