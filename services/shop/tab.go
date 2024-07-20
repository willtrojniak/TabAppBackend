package shop

import (
	"context"
	"errors"
	"reflect"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/types"
)

func (h *Handler) CreateTab(ctx context.Context, session *sessions.Session, data *types.TabCreate) error {

	err := types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.CreateTab(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) UpdateTab(ctx context.Context, session *sessions.Session, shopId int, tabId int, data *types.TabUpdate) error {
	userId, err := session.GetUserId()
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
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

	if reflect.DeepEqual(tab.TabCreate.TabUpdate.TabBase, data.TabBase) && reflect.DeepEqual(tab.VerificationList, data.VerificationList) {
		return nil
	}

	if authErr == nil || (userId == tab.OwnerId && tab.Status == types.TAB_STATUS_PENDING.String()) {
		// If the request is made by the shop owner, or if the tab has not yet been confirmed, update the tab directly
		err = h.store.UpdateTab(ctx, shopId, tabId, data)
		if err != nil {
			return err
		}
	} else if tab.Status == types.TAB_STATUS_CONFIRMED.String() {
		// Here it must be the case that the user is the tab owner, so only request an update so long as the tab is in the confirmed state
		if !reflect.DeepEqual(tab.TabCreate.TabUpdate.TabBase, data.TabBase) {
			err = h.store.SetTabUpdates(ctx, shopId, tabId, data)
		} else {
			err = h.store.SetTabUsers(ctx, shopId, tabId, data.VerificationList)
		}
		if err != nil {
			return err
		}
	} else if tab.Status == types.TAB_STATUS_CLOSED.String() {
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
	if tab.Status != types.TAB_STATUS_PENDING.String() && tab.Status != types.TAB_STATUS_CONFIRMED.String() {
		return services.NewDataConflictServiceError(nil)
	}

	err = h.store.ApproveTab(ctx, shopId, tabId)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetTabs(ctx context.Context, session *sessions.Session, shopId int) ([]types.TabOverview, error) {
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

func (h *Handler) GetTabById(ctx context.Context, session *sessions.Session, shopId int, tabId int) (types.Tab, error) {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return types.Tab{}, err
	}

	tab, err := h.store.GetTabById(ctx, shopId, tabId)
	if err != nil {
		return types.Tab{}, err
	}

	return tab, nil
}

func (h *Handler) AddOrderToTab(ctx context.Context, session *sessions.Session, shopId int, tabId int, data *types.BillOrderCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
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

func (h *Handler) RemoveOrderFromTab(ctx context.Context, session *sessions.Session, shopId int, tabId int, data *types.BillOrderCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
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
