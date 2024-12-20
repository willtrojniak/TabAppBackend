package shop

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

const (
	shopIdParam              = "shopId"
	locationIdParam          = "locationId"
	categoryIdParam          = "categoryId"
	itemIdParam              = "itemId"
	itemVariantIdParam       = "itemVariantId"
	substitutionGroupIdParam = "substitutionGroupId"
	tabIdParam               = "tabId"
	billIdParam              = "billId"
)

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	h.logger.Info("Registering shop routes")
	router.HandleFunc("POST /shops", h.sessions.WithAuthedSession(h.handleCreateShop))
	router.HandleFunc("GET /shops", h.sessions.WithAuthedSession(h.handleGetShops))

	// Payment Methods
	router.HandleFunc("GET /payment-methods", h.sessions.WithAuthedSession(h.handleGetPaymentMethods))

	// Shops
	router.HandleFunc(fmt.Sprintf("GET /shops/{%v}", shopIdParam), h.sessions.WithAuthedSession(h.handleGetShopById))
	router.HandleFunc(fmt.Sprintf("PATCH /shops/{%v}", shopIdParam), h.sessions.WithAuthedSession(h.handleUpdateShop))
	router.HandleFunc(fmt.Sprintf("DELETE /shops/{%v}", shopIdParam), h.sessions.WithAuthedSession(h.handleDeleteShop))

	// Users & Permissions
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/users/invite", shopIdParam), h.sessions.WithAuthedSession(h.handleInviteUser))
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/users/remove", shopIdParam), h.sessions.WithAuthedSession(h.handleRemoveUser))
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/accept", shopIdParam), h.sessions.WithAuthedSession(h.handleAcceptInvite))

	// Locations
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/locations", shopIdParam), h.sessions.WithAuthedSession(h.handleCreateLocation))
	router.HandleFunc(fmt.Sprintf("PATCH /shops/{%v}/locations/{%v}", shopIdParam, locationIdParam), h.sessions.WithAuthedSession(h.handleUpdateLocation))
	router.HandleFunc(fmt.Sprintf("DELETE /shops/{%v}/locations/{%v}", shopIdParam, locationIdParam), h.sessions.WithAuthedSession(h.handleDeleteLocation))

	// Categories
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/categories", shopIdParam), h.sessions.WithAuthedSession(h.handleCreateCategory))
	router.HandleFunc(fmt.Sprintf("GET /shops/{%v}/categories", shopIdParam), h.sessions.WithAuthedSession(h.handleGetCategories))
	router.HandleFunc(fmt.Sprintf("PATCH /shops/{%v}/categories/{%v}", shopIdParam, categoryIdParam), h.sessions.WithAuthedSession(h.handleUpdateCategory))
	router.HandleFunc(fmt.Sprintf("DELETE /shops/{%v}/categories/{%v}", shopIdParam, categoryIdParam), h.sessions.WithAuthedSession(h.handleDeleteCategory))

	// Items
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/items", shopIdParam), h.sessions.WithAuthedSession(h.handleCreateItem))
	router.HandleFunc(fmt.Sprintf("GET /shops/{%v}/items", shopIdParam), h.sessions.WithAuthedSession(h.handleGetItems))
	router.HandleFunc(fmt.Sprintf("PATCH /shops/{%v}/items/{%v}", shopIdParam, itemIdParam), h.sessions.WithAuthedSession(h.handleUpdateItem))
	router.HandleFunc(fmt.Sprintf("GET /shops/{%v}/items/{%v}", shopIdParam, itemIdParam), h.sessions.WithAuthedSession(h.handleGetItem))
	router.HandleFunc(fmt.Sprintf("DELETE /shops/{%v}/items/{%v}", shopIdParam, itemIdParam), h.sessions.WithAuthedSession(h.handleDeleteItem))

	// Item Variants
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/items/{%v}/variants", shopIdParam, itemIdParam), h.sessions.WithAuthedSession(h.handleCreateItemVariant))
	router.HandleFunc(fmt.Sprintf("PATCH /shops/{%v}/items/{%v}/variants/{%v}", shopIdParam, itemIdParam, itemVariantIdParam), h.sessions.WithAuthedSession(h.handleUpdateItemVariant))
	router.HandleFunc(fmt.Sprintf("DELETE /shops/{%v}/items/{%v}/variants/{%v}", shopIdParam, itemIdParam, itemVariantIdParam), h.sessions.WithAuthedSession(h.handleDeleteItemVariant))

	// Item Substitution Groups
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/substitutions", shopIdParam), h.sessions.WithAuthedSession(h.handleCreateSubstitutionGroup))
	router.HandleFunc(fmt.Sprintf("GET /shops/{%v}/substitutions", shopIdParam), h.sessions.WithAuthedSession(h.handleGetSubstitutionGroups))
	router.HandleFunc(fmt.Sprintf("PATCH /shops/{%v}/substitutions/{%v}", shopIdParam, substitutionGroupIdParam), h.sessions.WithAuthedSession(h.handleUpdateSubstitutionGroup))
	router.HandleFunc(fmt.Sprintf("DELETE /shops/{%v}/substitutions/{%v}", shopIdParam, substitutionGroupIdParam), h.sessions.WithAuthedSession(h.handleDeleteSubstitutionGroup))

	// Tabs
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/tabs", shopIdParam), h.sessions.WithAuthedSession(h.handleCreateTab))
	router.HandleFunc(fmt.Sprintf("GET /shops/{%v}/tabs", shopIdParam), h.sessions.WithAuthedSession(h.handleGetTabs))
	router.HandleFunc(fmt.Sprintf("GET /shops/{%v}/tabs/{%v}", shopIdParam, tabIdParam), h.sessions.WithAuthedSession(h.handleGetTabById))
	router.HandleFunc(fmt.Sprintf("PATCH /shops/{%v}/tabs/{%v}", shopIdParam, tabIdParam), h.sessions.WithAuthedSession(h.handleUpdateTab))
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/tabs/{%v}/approve", shopIdParam, tabIdParam), h.sessions.WithAuthedSession(h.handleApproveTab))
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/tabs/{%v}/close", shopIdParam, tabIdParam), h.sessions.WithAuthedSession(h.handleCloseTab))
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/tabs/{%v}/bills/{%v}/close", shopIdParam, tabIdParam, billIdParam), h.sessions.WithAuthedSession(h.handleCloseTabBill))

	// Orders
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/tabs/{%v}/add-order", shopIdParam, tabIdParam), h.sessions.WithAuthedSession(h.handleAddOrderToTab))
	router.HandleFunc(fmt.Sprintf("POST /shops/{%v}/tabs/{%v}/remove-order", shopIdParam, tabIdParam), h.sessions.WithAuthedSession(h.handleRemoveOrderFromTab))

}

func (h *Handler) handleCreateShop(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {

	data := &models.ShopCreate{}
	err := models.ReadRequestJson(r, data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	_, err = h.CreateShop(r.Context(), session, data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleGetShops(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	// Query params
	const memberKey = "member"
	const pendingKey = "pending"

	var params models.GetShopsQueryParams

	// TODO: Dynamically change limit and offset
	params.Limit = 10
	params.Offset = 0

	rawParams := r.URL.Query()
	if rawParams.Has(memberKey) {
		if isMember, err := strconv.ParseBool(rawParams.Get(memberKey)); err == nil {
			params.IsMember = &isMember
			params.UserId = &session.UserId
		}
	}

	if rawParams.Has(pendingKey) {
		if pending, err := strconv.ParseBool(rawParams.Get(pendingKey)); err == nil {
			params.IsPending = &pending
			if params.IsMember == nil {
				isMember := true
				params.IsMember = &isMember
			}
			params.UserId = &session.UserId
		}
	}

	shops, err := h.GetShops(r.Context(), &params)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shops)
	return
}

func (h *Handler) handleGetShopById(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shopId"))
		return
	}

	shop, err := h.GetShopById(r.Context(), session, shopId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shop)

}

func (h *Handler) handleUpdateShop(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shopId"))
		return
	}

	data := models.ShopUpdate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.UpdateShop(r.Context(), session, shopId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleDeleteShop(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shopId"))
		return
	}

	err = h.DeleteShop(r.Context(), session, shopId)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleInviteUser(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shopId"))
		return
	}

	data := models.ShopUserCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.InviteUserToShop(r.Context(), session, shopId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleRemoveUser(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shopId"))
		return
	}

	data := models.ShopUserCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.RemoveUserFromShop(r.Context(), session, shopId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleAcceptInvite(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shopId"))
		return
	}

	err = h.AcceptInviteToShop(r.Context(), session, shopId)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleGetPaymentMethods(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	methods := make([]models.PaymentMethod, 0)
	methods = append(methods, models.PaymentMethodInPerson, models.PaymentMethodChartstring)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methods)
}

func (h *Handler) handleCreateLocation(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	data := models.LocationCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
	data.ShopId = shopId

	err = h.CreateLocation(r.Context(), session, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleUpdateLocation(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	locationId, err := strconv.Atoi(r.PathValue(locationIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid location id"))
		return
	}

	data := models.LocationUpdate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.UpdateLocation(r.Context(), session, shopId, locationId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleDeleteLocation(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	locationId, err := strconv.Atoi(r.PathValue(locationIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid location id"))
		return
	}

	err = h.DeleteLocation(r.Context(), session, shopId, locationId)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleCreateCategory(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	data := models.CategoryCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
	data.ShopId = shopId

	err = h.CreateCategory(r.Context(), session, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleGetCategories(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	categories, err := h.GetCategories(r.Context(), shopId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *Handler) handleUpdateCategory(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	categoryId, err := strconv.Atoi(r.PathValue(categoryIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	data := models.CategoryUpdate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.UpdateCategory(r.Context(), session, shopId, categoryId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleDeleteCategory(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	categoryId, err := strconv.Atoi(r.PathValue(categoryIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	err = h.DeleteCategory(r.Context(), session, shopId, categoryId)
	if err != nil {
		h.handleError(w, err)
		return
	}

}

func (h *Handler) handleCreateItem(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	data := models.ItemCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
	data.ShopId = shopId

	err = h.CreateItem(r.Context(), session, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleGetItems(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	items, err := h.GetItems(r.Context(), session, shopId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *Handler) handleUpdateItem(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	itemId, err := strconv.Atoi(r.PathValue(itemIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid item id"))
		return
	}

	data := models.ItemUpdate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.UpdateItem(r.Context(), session, shopId, itemId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

}

func (h *Handler) handleGetItem(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	itemId, err := strconv.Atoi(r.PathValue(itemIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid item id"))
		return
	}

	item, err := h.GetItem(r.Context(), session, shopId, itemId)
	if err != nil {
		h.handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *Handler) handleDeleteItem(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	itemId, err := strconv.Atoi(r.PathValue(itemIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid item id"))
		return
	}

	err = h.DeleteItem(r.Context(), session, shopId, itemId)
	if err != nil {
		h.handleError(w, err)
		return
	}

}

func (h *Handler) handleCreateItemVariant(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	itemId, err := strconv.Atoi(r.PathValue(itemIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid item id"))
		return
	}

	data := models.ItemVariantCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
	data.ShopId = shopId
	data.ItemId = itemId

	err = h.CreateItemVariant(r.Context(), session, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	return
}

func (h *Handler) handleUpdateItemVariant(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	itemId, err := strconv.Atoi(r.PathValue(itemIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid item id"))
		return
	}
	variantId, err := strconv.Atoi(r.PathValue(itemVariantIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid item variant id"))
		return
	}

	data := models.ItemVariantUpdate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.UpdateItemVariant(r.Context(), session, shopId, itemId, variantId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleDeleteItemVariant(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	itemId, err := strconv.Atoi(r.PathValue(itemIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid item id"))
		return
	}
	variantId, err := strconv.Atoi(r.PathValue(itemVariantIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid item variant id"))
		return
	}

	err = h.DeleteItemVariant(r.Context(), session, shopId, itemId, variantId)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleCreateSubstitutionGroup(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	data := models.SubstitutionGroupCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
	data.ShopId = shopId

	err = h.CreateSubstitutionGroup(r.Context(), session, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleUpdateSubstitutionGroup(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	substitutionGroupId, err := strconv.Atoi(r.PathValue(substitutionGroupIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid substitution id"))
		return
	}

	data := models.SubstitutionGroupUpdate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.UpdateSubstitutionGroup(r.Context(), session, shopId, substitutionGroupId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleGetSubstitutionGroups(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	substitutionGroups, err := h.GetSubstitutionGroups(r.Context(), session, shopId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(substitutionGroups)
}

func (h *Handler) handleDeleteSubstitutionGroup(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	substitutionGroupId, err := strconv.Atoi(r.PathValue(substitutionGroupIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid substitution id"))
		return
	}

	err = h.DeleteSubstitutionGroup(r.Context(), session, shopId, substitutionGroupId)
	if err != nil {
		h.handleError(w, err)
		return
	}

}

func (h *Handler) handleCreateTab(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	data := models.TabCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
	data.ShopId = shopId
	data.OwnerId = session.UserId

	err = h.CreateTab(r.Context(), session, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleGetTabs(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	tabs, err := h.GetTabs(r.Context(), session, shopId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tabs)

}

func (h *Handler) handleGetTabById(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}
	tabId, err := strconv.Atoi(r.PathValue(tabIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid tab id"))
		return
	}

	tab, err := h.GetTabById(r.Context(), session, shopId, tabId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tab)

}

func (h *Handler) handleUpdateTab(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	tabId, err := strconv.Atoi(r.PathValue(tabIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid tab id"))
		return
	}

	data := models.TabUpdate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.UpdateTab(r.Context(), session, shopId, tabId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

}

func (h *Handler) handleApproveTab(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	tabId, err := strconv.Atoi(r.PathValue(tabIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid tab id"))
		return
	}

	err = h.ApproveTab(r.Context(), session, shopId, tabId)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleCloseTab(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	tabId, err := strconv.Atoi(r.PathValue(tabIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid tab id"))
		return
	}

	err = h.CloseTab(r.Context(), session, shopId, tabId)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *Handler) handleAddOrderToTab(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	tabId, err := strconv.Atoi(r.PathValue(tabIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid tab id"))
		return
	}

	data := models.BillOrderCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.AddOrderToTab(r.Context(), session, shopId, tabId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

}

func (h *Handler) handleRemoveOrderFromTab(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	tabId, err := strconv.Atoi(r.PathValue(tabIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid tab id"))
		return
	}

	data := models.BillOrderCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

	err = h.RemoveOrderFromTab(r.Context(), session, shopId, tabId, &data)
	if err != nil {
		h.handleError(w, err)
		return
	}

}

func (h *Handler) handleCloseTabBill(w http.ResponseWriter, r *http.Request, session *sessions.AuthedSession) {
	shopId, err := strconv.Atoi(r.PathValue(shopIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid shop id"))
		return
	}

	tabId, err := strconv.Atoi(r.PathValue(tabIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid tab id"))
		return
	}

	billId, err := strconv.Atoi(r.PathValue(billIdParam))
	if err != nil {
		h.handleError(w, services.NewValidationServiceError(err, "Invalid bill id"))
		return
	}

	err = h.MarkTabBillPaid(r.Context(), session, shopId, tabId, billId)
	if err != nil {
		h.handleError(w, err)
		return
	}
}
