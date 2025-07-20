package auth

import (
	"fmt"
	"net/http"

	"github.com/willtrojniak/TabAppBackend/env"
)

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /auth/google/callback", h.handleAuthCallback)
	router.HandleFunc("GET /auth/google", h.handleAuth)
	router.HandleFunc("POST /logout", h.handleLogout)
	h.logger.Info("Registered auth routes")
}

func (h *Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauth2token, err := h.Authorize(r, &h.googleCfg)
	if err != nil {
		h.handleError(w, err)
		return
	}

	user, err := h.createUserFromToken(r, oauth2token, &h.googleCfg, h.googleProvider)
	if err != nil {
		h.handleError(w, err)
		return
	}

	_, err = h.sessionManager.SetNewSession(w, r, user)
	if err != nil {
		h.handleError(w, err)
		return
	}

	redirectCookie, err := r.Cookie("redirect")
	redirect := ""
	if err == nil {
		redirect = redirectCookie.Value
	}

	http.Redirect(w, r, fmt.Sprintf("%v/%v", env.Envs.UI_URI, redirect), http.StatusFound)
	return

}

func (h *Handler) handleAuth(w http.ResponseWriter, r *http.Request) {
	if err := h.BeginAuthorize(w, r, &h.googleCfg); err != nil {
		h.handleError(w, err)
		return
	}

	return
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {

	if err := h.logout(w, r); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
