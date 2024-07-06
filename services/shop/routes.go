package shop

import (
	"encoding/json"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/types"
)

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	h.logger.Info("Registering shop routes")

	router.HandleFunc("POST /shops", h.handleCreateShop)
	router.HandleFunc("GET /shops", h.handleGetShops)

	subrouter := http.NewServeMux()
	router.Handle("/shops/", http.StripPrefix("/shops", subrouter))

}

func (h *Handler) handleCreateShop(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessions.GetSession(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	data := &types.ShopCreate{}
	err = types.ReadRequestJson(r, data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.CreateShop(r.Context(), session, data)
	if err != nil {
		h.handleError(w, err)
		return
	}

}

func (h *Handler) handleGetShops(w http.ResponseWriter, r *http.Request) {
	// TODO: Dynamically change limit and offset
	shops, err := h.GetShops(r.Context(), 10, 0)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shops)
	return
}
