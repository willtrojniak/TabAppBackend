package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/types"
)

const userIdPath = "userId"

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	h.logger.Info("Registering user routes")

	router.HandleFunc("GET /users", h.handleGetUser)
	router.HandleFunc(fmt.Sprintf("PATCH /users/{%v}", userIdPath), h.handleUpdateUser)

}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessions.GetSession(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	user, err := h.GetUser(r.Context(), session)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	return
}

func (h *Handler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessions.GetSession(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	userId := r.PathValue(userIdPath)

	data := types.UserUpdate{}
	err = types.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.UpdateUser(r.Context(), session, userId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}
