package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/env"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
  maxAge = 86400 * 30 // 30 days
  session_cookie = "user_session"
)

type CreateUserFn func(context context.Context, user *types.UserCreate) (*uuid.UUID, error)

type Handler struct{
  logger *slog.Logger;
  store sessions.Store
  handleError services.HTTPErrorHandler;
  createUser CreateUserFn;
}

func NewHandler(handleError services.HTTPErrorHandler, logger *slog.Logger) *Handler {
  return &Handler{
    store: gothic.Store, 
    handleError: handleError,
    logger: logger,
  };
}

func (h *Handler) SetCreateUserFn(fn CreateUserFn) {
  h.createUser = fn;
}

func init() {

  store := sessions.NewCookieStore([]byte(env.Envs.SESSION_SECRET));
  store.MaxAge(maxAge);
  store.Options.Path = "/"
  store.Options.HttpOnly = true;
  store.Options.Secure = true; // BREAKS IF FALSE

  gothic.Store = store;

  goth.UseProviders(
    // TODO: Make URL dynamic
    google.New(env.Envs.OAUTH2_GOOGLE_CLIENT_ID, env.Envs.OAUTH2_GOOGLE_CLIENT_SECRET, "http://127.0.0.1:3000/auth/google/callback",
      "openid", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"),
  );
}

func (h *Handler) beginAuthorize(w http.ResponseWriter, r *http.Request, provider string) error {
  r = r.WithContext(context.WithValue(context.Background(), "provider", provider));
  gothic.BeginAuthHandler(w, r);

  return nil;
}

func (h *Handler) authorize(w http.ResponseWriter, r *http.Request, provider string) error {
  r = r.WithContext(context.WithValue(context.Background(), "provider", provider));
  user, err := gothic.CompleteUserAuth(w, r);
  if err != nil {
    h.logger.Error("Error while authenticating with goth");
    return services.NewInternalServiceError(err);
  }

  userUUID, err := h.createUser(context.Background(), &types.UserCreate{Email: user.Email, Name: user.Name});
  if err != nil {
    return services.NewInternalServiceError(err);
  }
  
  if err := h.storeUserUUID(w, r, userUUID); err != nil {
    return services.NewInternalServiceError(err);
  }

  return nil;
}

func (h *Handler) storeUserUUID(w http.ResponseWriter, r *http.Request, id *uuid.UUID) error {

  session, _ := h.store.Get(r, session_cookie);
  session.Values["user"] = id;
  
  err := session.Save(r, w);
  if err != nil {
    h.logger.Error("Error while storing userUUID");
    return services.NewInternalServiceError(err);
  }
  return nil;
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) error {

  session, err := h.store.Get(r, session_cookie);
  if err != nil {
    return services.NewInternalServiceError(err);
  }
  session.Options.MaxAge = -1;
  if err := session.Save(r, w); err != nil {
    return services.NewInternalServiceError(err);
  }

  return nil;
}

func (h *Handler) GetUserSession(r *http.Request) (uuid.UUID, error) {
  session, err := h.store.Get(r, session_cookie);
  if err != nil {
    return uuid.UUID{}, services.NewUnauthorizedServiceError(err);
  }

  u := session.Values["user"];
  if u == nil {
    return uuid.UUID{}, services.NewUnauthorizedServiceError(err);
  }

  return u.(uuid.UUID), nil;
}

func (h *Handler) RequireAuth(next http.Handler) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    _, err := h.GetUserSession(r);
    if err != nil {
      h.handleError(w, err);
      return;
    }
    next.ServeHTTP(w, r);
  }
}

