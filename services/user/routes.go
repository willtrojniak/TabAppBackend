package user

import (
	"encoding/json"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/services"
)

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	h.logger.Info("Registering user routes")

	router.HandleFunc("/users", h.handleGetUser)

}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessions.GetSession(r)
	if err != nil || session.UserId == "" {
		h.handleError(w, services.NewUnauthorizedServiceError(err))
		return
	}

	user, err := h.GetUser(r.Context(), session.UserId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	return
}
