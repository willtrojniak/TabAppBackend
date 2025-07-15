package user

import (
	"encoding/json"
	"net/http"

	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services/sessions"
)

const userIdPath = "userId"

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	h.logger.Info("Registering user routes")

	router.HandleFunc("GET /users", h.handleGetUser)
	router.HandleFunc("PATCH /users", h.handleUpdateUser)

}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling get user")
	h.sessions.WithAuthedSession(func(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {

		user, err := h.GetUser(r.Context(), session)
		if err != nil {
			h.handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return
	})(w, r)
}

func (h *Handler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling Update User")
	h.sessions.WithAuthedSession(func(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
		data := models.UserUpdate{}
		err := models.ReadRequestJson(r, &data)
		if err != nil {
			h.handleError(w, err)
			return
		}
		if *data.PreferredName == "" {
			data.PreferredName = nil
		}

		err = h.UpdateUser(r.Context(), session, session.UserId, &data)
		if err != nil {
			h.handleError(w, err)
			return
		}
	})(w, r)
}
