package user

import (
	"encoding/json"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/models"
)

const userIdPath = "userId"

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	h.logger.Info("Registering user routes")

	router.HandleFunc("GET /users", h.handleGetUser)
	router.HandleFunc("PATCH /users", h.handleUpdateUser)

}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	h.sessions.WithAuthedSessionUserId(func(w http.ResponseWriter, r *http.Request, sUserId string) {
		h.logger.Debug("Handling get user")

		user, err := h.GetUser(r.Context(), sUserId)
		if err != nil {
			h.handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return
	})
}

func (h *Handler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	h.sessions.WithAuthedSessionUserId(func(w http.ResponseWriter, r *http.Request, sUserId string) {
		data := models.UserUpdate{}
		err := models.ReadRequestJson(r, &data)
		if err != nil {
			h.handleError(w, err)
			return
		}

		err = h.UpdateUser(r.Context(), sUserId, sUserId, &data)
		if err != nil {
			h.handleError(w, err)
			return
		}
	})
}
