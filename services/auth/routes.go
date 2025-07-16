package auth

import (
	"fmt"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/willtrojniak/TabAppBackend/env"
	"github.com/willtrojniak/TabAppBackend/services"
)

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /auth/{provider}/callback", h.handleAuthCallback)
	router.HandleFunc("GET /auth/{provider}", h.handleAuth)
	router.HandleFunc("POST /logout", h.handleLogout)
	h.logger.Info("Registered auth routes")
}

func (h *Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	session, _ := h.sessionManager.GetSession(r)
	if provider != "google" && (session == nil || !session.IsAuthed()) {
		h.handleError(w, services.NewUnauthenticatedServiceError(nil))
		return
	}

	oauth2token, err := h.authorize(r, provider)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if provider == "google" {
		user, err := h.createUserFromToken(r, oauth2token, provider)
		if err != nil {
			h.handleError(w, err)
			return
		}

		_, err = h.sessionManager.SetNewSession(w, r, user)
		if err != nil {
			h.handleError(w, err)
			return
		}
	} else if provider == "slack" {
		api := slack.New(oauth2token.AccessToken)
		one, two, three, err := api.SendMessage("#test", slack.MsgOptionBlocks(slack.NewMarkdownBlock("md", "Success")))
		h.logger.Debug("Result", "one", one, "two", two, "three", three, "err", err)
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
	provider := r.PathValue("provider")
	if err := h.beginAuthorize(w, r, provider); err != nil {
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
