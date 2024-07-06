package shop

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
)

const shopIdParam = "shopId"

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	h.logger.Info("Registering shop routes")

	router.HandleFunc("POST /shops", h.handleCreateShop)
	router.HandleFunc("GET /shops", h.handleGetShops)

	subrouter := http.NewServeMux()
	router.Handle("/shops/", http.StripPrefix("/shops", subrouter))

	subrouter.HandleFunc(fmt.Sprintf("GET /{%v}", shopIdParam), h.handleGetShopById)
	subrouter.HandleFunc(fmt.Sprintf("PATCH /{%v}", shopIdParam), h.handleUpdateShop)

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

func (h *Handler) handleGetShopById(w http.ResponseWriter, r *http.Request) {
	shopId, err := uuid.Parse(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shopId"))
		return
	}

	shop, err := h.GetShopById(r.Context(), shopId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shop)

}

func (h *Handler) handleUpdateShop(w http.ResponseWriter, r *http.Request) {
	shopId, err := uuid.Parse(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shopId"))
		return
	}

	session, err := h.sessions.GetSession(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	data := types.ShopUpdate{}
	err = types.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
	data.Id = shopId

	err = h.UpdateShop(r.Context(), session, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}
