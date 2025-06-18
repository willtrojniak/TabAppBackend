package auth

import (
	"fmt"
	"net/http"

	"github.com/willtrojniak/TabAppBackend/env"
)

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /auth/{provider}/callback", h.handleAuthCallback)
	router.HandleFunc("GET /auth/{provider}", h.handleAuth)
	router.HandleFunc("POST /logout", h.handleLogout)
	h.logger.Info("Registered auth routes")
}

func (h *Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	err := h.authorize(w, r)
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
	if err := h.beginAuthorize(w, r); err != nil {
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
