package shop

import (
	"context"
	"errors"
	"reflect"

	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateTab(ctx context.Context, session *sessions.Session, data *models.TabCreate) error {
	// Check that the client is authenticated
	userId, err := session.GetUserId()
	if err != nil {
		return err
	}

	// Request data validation
	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	// TODO: Wrap in transaction

	// Get the target shop
	shop, err := h.store.GetShopById(ctx, data.ShopId)
	if err != nil {
		return err
	}

	// By default the tab status is pending, unless it is created by the shop owner
	status := models.TAB_STATUS_PENDING
	if shop.OwnerId == userId {
		status = models.TAB_STATUS_CONFIRMED
	}

	err = h.store.CreateTab(ctx, data, status)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) UpdateTab(ctx context.Context, session *sessions.Session, shopId int, tabId int, data *models.TabUpdate) error {
	userId, err := session.GetUserId()
	if err != nil {
		return err
	}

	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	// TODO: Perform within a transaction
	tab, err := h.store.GetTabById(ctx, shopId, tabId)
	if err != nil {
		return err
	}

	authErr := h.AuthorizeModifyShop(ctx, session, shopId)
	if authErr != nil && userId != tab.OwnerId {
		return authErr
	}

	if reflect.DeepEqual(tab.TabBase, data.TabBase) && reflect.DeepEqual(tab.VerificationList, data.VerificationList) {
		return nil
	}

	if authErr == nil || (userId == tab.OwnerId && tab.Status == models.TAB_STATUS_PENDING.String()) {
		// If the request is made by the shop owner, or if the tab has not yet been confirmed, update the tab directly
		err = h.store.UpdateTab(ctx, shopId, tabId, data)
		if err != nil {
			return err
		}
	} else if tab.Status == models.TAB_STATUS_CONFIRMED.String() {
		tabLocationIds := make([]uint, 0)
		for _, location := range tab.Locations {
			tabLocationIds = append(tabLocationIds, location.Id)
		}
		// Here it must be the case that the user is the tab owner, so only request an update so long as the tab is in the confirmed state
		if !reflect.DeepEqual(tab.TabBase, data.TabBase) || !reflect.DeepEqual(tabLocationIds, data.LocationIds) {
			err = h.store.SetTabUpdates(ctx, shopId, tabId, data)
		} else {
			err = h.store.SetTabUsers(ctx, shopId, tabId, data.VerificationList)
		}
		if err != nil {
			return err
		}
	} else if tab.Status == models.TAB_STATUS_CLOSED.String() {
		return services.NewDataConflictServiceError(errors.New("Cannot update closed tab"))
	} else {
		h.logger.Error("Unknown tab state in update tab", "userId", userId, "tab", tab)
		return services.NewInternalServiceError(errors.New("Unknown tab state"))
	}

	return nil
}

func (h *Handler) ApproveTab(ctx context.Context, session *sessions.Session, shopId int, tabId int) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	tab, err := h.store.GetTabById(ctx, shopId, tabId)
	if err != nil {
		return err
	}
	if tab.Status != models.TAB_STATUS_PENDING.String() && tab.Status != models.TAB_STATUS_CONFIRMED.String() {
		return services.NewDataConflictServiceError(nil)
	}

	err = h.store.ApproveTab(ctx, shopId, tabId)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) CloseTab(ctx context.Context, session *sessions.Session, shopId int, tabId int) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	tab, err := h.store.GetTabById(ctx, shopId, tabId)
	if err != nil {
		return err
	}
	if !(tab.Status == models.TAB_STATUS_PENDING.String() || tab.Status == models.TAB_STATUS_CONFIRMED.String()) {
		return services.NewDataConflictServiceError(nil)
	}

	err = h.store.CloseTab(ctx, shopId, tabId)
	if err != nil {
		return err
	}

	return nil

}

func (h *Handler) MarkTabBillPaid(ctx context.Context, session *sessions.Session, shopId int, tabId int, billId int) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = h.store.MarkTabBillPaid(ctx, shopId, tabId, billId)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetTabs(ctx context.Context, session *sessions.Session, shopId int) ([]models.TabOverview, error) {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return nil, err
	}

	tabs, err := h.store.GetTabs(ctx, shopId)
	if err != nil {
		return nil, err
	}

	return tabs, nil
}

func (h *Handler) GetTabById(ctx context.Context, session *sessions.Session, shopId int, tabId int) (models.Tab, error) {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return models.Tab{}, err
	}

	tab, err := h.store.GetTabById(ctx, shopId, tabId)
	if err != nil {
		return models.Tab{}, err
	}

	return tab, nil
}

func (h *Handler) AddOrderToTab(ctx context.Context, session *sessions.Session, shopId int, tabId int, data *models.BillOrderCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	tab, err := h.store.GetTabById(ctx, shopId, tabId)
	if err != nil {
		return err
	}

	if !tab.IsActive() {
		return services.NewDataConflictServiceError(nil)
	}

	err = h.store.AddOrderToTab(ctx, shopId, tabId, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) RemoveOrderFromTab(ctx context.Context, session *sessions.Session, shopId int, tabId int, data *models.BillOrderCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	tab, err := h.store.GetTabById(ctx, shopId, tabId)
	if err != nil {
		return err
	}

	if !tab.IsActive() {
		return services.NewDataConflictServiceError(nil)
	}

	err = h.store.RemoveOrderFromTab(ctx, shopId, tabId, data)
	if err != nil {
		return err
	}

	return nil
}
