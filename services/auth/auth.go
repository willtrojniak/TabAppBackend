package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/willtrojniak/TabAppBackend/env"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services"
	"github.com/willtrojniak/TabAppBackend/services/sessions"
	"github.com/willtrojniak/TabAppBackend/services/user"
	"golang.org/x/oauth2"
)

type Handler struct {
	logger         *slog.Logger
	handleError    services.HTTPErrorHandler
	userHandler    *user.Handler
	googleCfg      oauth2.Config
	googleProvider *oidc.Provider
	sessionManager *sessions.Handler
}

func NewHandler(handleError services.HTTPErrorHandler, sessionManager *sessions.Handler, userHandler *user.Handler, logger *slog.Logger) (*Handler, error) {
	googleProvider, err := oidc.NewProvider(context.TODO(), "https://accounts.google.com")
	if err != nil {
		return nil, err
	}
	googleCfg := oauth2.Config{
		ClientID:     env.Envs.OAUTH2_GOOGLE_CLIENT_ID,
		ClientSecret: env.Envs.OAUTH2_GOOGLE_CLIENT_SECRET,
		Endpoint:     googleProvider.Endpoint(),
		RedirectURL:  fmt.Sprintf("%v/auth/google/callback", env.Envs.BASE_URI),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Handler{
		handleError:    handleError,
		logger:         logger,
		userHandler:    userHandler,
		googleCfg:      googleCfg,
		googleProvider: googleProvider,
		sessionManager: sessionManager,
	}, nil
}

func (h *Handler) BeginAuthorize(w http.ResponseWriter, r *http.Request, cfg *oauth2.Config) error {
	state, err := randString(16)
	if err != nil {
		return err
	}
	nonce, err := randString(16)
	if err != nil {
		return err
	}

	redirect := r.URL.Query().Get("redirect")
	h.logger.Debug("Redirect query param", "value", redirect)

	setCallbackCookie(w, r, "state", state)
	setCallbackCookie(w, r, "nonce", nonce)
	setCallbackCookie(w, r, "redirect", redirect)
	http.Redirect(w, r, cfg.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)

	return nil
}

func (h *Handler) Authorize(r *http.Request, cfg *oauth2.Config) (*oauth2.Token, error) {
	// Check that the CSRF token matches
	state, err := r.Cookie("state")
	if err != nil {
		h.logger.Warn("Missing State Cookie")
		return nil, services.NewInternalServiceError(err)
	}

	if r.URL.Query().Get("state") != state.Value {
		return nil, fmt.Errorf("State did not match")
	}

	// Exchange the code for a token
	oauth2Token, err := cfg.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}
	return oauth2Token, nil
}

func (h *Handler) createUserFromToken(r *http.Request, token *oauth2.Token, cfg *oauth2.Config, provider *oidc.Provider) (*models.User, error) {
	// Get the id token from the JWT
	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("No id_token field in oauth2 token")
	}

	// Verify the id token
	oidcConfig := &oidc.Config{ClientID: cfg.ClientID}
	verifier := provider.Verifier(oidcConfig)

	idToken, err := verifier.Verify(r.Context(), rawIdToken)
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}

	nonce, err := r.Cookie("nonce")
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}

	if idToken.Nonce != nonce.Value {
		return nil, fmt.Errorf("nonce did not match")
	}

	// Get user data from the id token
	var claims struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		Sub           string `json:"sub"`
		EmailVerified bool   `json:"email_verified"`
	}

	userInfo, err := provider.UserInfo(r.Context(), oauth2.StaticTokenSource(token))
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}

	if err := userInfo.Claims(&claims); err != nil {
		return nil, services.NewInternalServiceError(err)
	}

	// Add the user to the database if not already
	user, err := h.userHandler.CreateUser(r.Context(), &models.UserCreate{Id: claims.Sub, Email: claims.Email, Name: claims.Name})
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}

	return user, nil
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) error {

	_, err := h.sessionManager.SetNewSession(w, r, nil)
	if err != nil {
		return err
	}

	return nil
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name string, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
