package shop

import (
	"context"
	"reflect"

	"github.com/willtrojniak/TabAppBackend/db"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services"
	"github.com/willtrojniak/TabAppBackend/services/authorization"
	"github.com/willtrojniak/TabAppBackend/services/events"
	"github.com/willtrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateTab(ctx context.Context, session *sessions.AuthedSession, data *models.TabCreate) error {
	// Request data validation
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return WithAuthorizeShopAction(ctx, h.store, session, data.ShopId, authorization.SHOP_ACTION_REQUEST_TAB, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		// By default the tab status is pending, unless it is created by user with role
		status := models.TAB_STATUS_PENDING

		// Check if the user has permission to create/manage tabs
		if ok, err := authorization.AuthorizeShopAction(user, shop, authorization.SHOP_ACTION_CREATE_TAB); err == nil && ok {
			status = models.TAB_STATUS_CONFIRMED
		}

		tabId, err := pq.CreateTab(ctx, data, status)
		if err != nil {
			return err
		}

		tab, err := pq.GetTabById(ctx, int(shop.Id), tabId)
		if err != nil {
			return nil
		}

		owner, err := pq.GetUser(ctx, data.OwnerId)
		if err != nil {
			return nil
		}

		events.Dispatch(h.eventDispatcher, events.TabCreateEvent{Tab: tab, Shop: shop, TabOwner: owner})

		return nil
	})
}

func (h *Handler) UpdateTab(ctx context.Context, session *sessions.AuthedSession, shopId int, tabId int, data *models.TabUpdate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	h.logger.Debug("Shop.UpdateTab")
	return WithAuthorizeTabAction(ctx, h.store, session, shopId, tabId, authorization.TAB_ACTION_REQUEST_UPDATE, func(pq *db.PgxQueries, user *models.User, shop *models.Shop, tab *models.Tab) error {
		tabLocationIds := make([]uint, 0)
		for _, location := range tab.Locations {
			tabLocationIds = append(tabLocationIds, location.Id)
		}
		h.logger.Debug("Shop.UpdateTab Authorized Request")

		// Check if the 'updates' are unchanged from current tab
		if reflect.DeepEqual(tab.TabBase, data.TabBase) && reflect.DeepEqual(tab.VerificationList, data.VerificationList) && reflect.DeepEqual(tabLocationIds, data.LocationIds) {
			return nil
		}

		// Next, check if have permission to update the tab directly
		if ok, err := authorization.AuthorizeTabAction(user, &authorization.TabTarget{Tab: tab, Shop: shop}, authorization.TAB_ACTION_UPDATE); err == nil && ok {
			h.logger.Debug("Shop.UpdateTab Authorized Direct Update")
			return pq.UpdateTab(ctx, shopId, tabId, data)
		}

		// Otherwise:
		// Check if part of the tab data has changed and request updates
		if !(reflect.DeepEqual(tab.TabBase, data.TabBase)) || !reflect.DeepEqual(tabLocationIds, data.LocationIds) {
			h.logger.Debug("Shop.UpdateTab One")
			err = pq.SetTabUpdates(ctx, shopId, tabId, data)
		} else {
			// If not, only update the people on the tab
			h.logger.Debug("Shop.UpdateTab Two")
			err = pq.SetTabUsers(ctx, shopId, tabId, data.VerificationList)
		}

		if err != nil {
			return err
		}

		owner, err := pq.GetUser(ctx, tab.OwnerId)
		if err != nil {
			return err
		}

		tab, err = pq.GetTabById(ctx, shopId, tabId)
		if err != nil {
			return err
		}

		events.Dispatch(h.eventDispatcher, events.TabUpdateEvent{Shop: shop, Tab: tab, TabOwner: owner})
		return nil
	})
}

func (h *Handler) ApproveTab(ctx context.Context, session *sessions.AuthedSession, shopId int, tabId int) error {
	return WithAuthorizeTabAction(ctx, h.store, session, shopId, tabId, authorization.TAB_ACTION_APPROVE, func(pq *db.PgxQueries, user *models.User, shop *models.Shop, tab *models.Tab) error {
		if tab.Status != models.TAB_STATUS_PENDING.String() && tab.Status != models.TAB_STATUS_CONFIRMED.String() {
			return services.NewDataConflictServiceError(nil)
		}

		err := pq.ApproveTab(ctx, shopId, tabId)
		if err != nil {
			return err
		}

		owner, err := pq.GetUser(ctx, tab.OwnerId)
		if err != nil {
			return err
		}
		events.Dispatch(h.eventDispatcher, events.TabApproveEvent{Shop: shop, Tab: tab, TabOwner: owner})

		return nil
	})
}

func (h *Handler) CloseTab(ctx context.Context, session *sessions.AuthedSession, shopId int, tabId int) error {
	return WithAuthorizeTabAction(ctx, h.store, session, shopId, tabId, authorization.TAB_ACTION_CLOSE, func(pq *db.PgxQueries, user *models.User, shop *models.Shop, tab *models.Tab) error {
		if !(tab.Status == models.TAB_STATUS_PENDING.String() || tab.Status == models.TAB_STATUS_CONFIRMED.String()) {
			return services.NewDataConflictServiceError(nil)
		}

		err := pq.CloseTab(ctx, shopId, tabId)
		if err != nil {
			return err
		}

		owner, err := pq.GetUser(ctx, tab.OwnerId)
		if err != nil {
			return err
		}

		events.Dispatch(h.eventDispatcher, events.TabCloseEvent{Shop: shop, Tab: tab, TabOwner: owner})
		return nil
	})
}

func (h *Handler) MarkTabBillPaid(ctx context.Context, session *sessions.AuthedSession, shopId int, tabId int, billId int) error {
	return WithAuthorizeTabAction(ctx, h.store, session, shopId, tabId, authorization.TAB_ACTION_CLOSE_BILL, func(pq *db.PgxQueries, user *models.User, shop *models.Shop, tab *models.Tab) error {
		err := pq.MarkTabBillPaid(ctx, shopId, tabId, billId)
		if err != nil {
			return err
		}

		tab, err = pq.GetTabById(ctx, shopId, tabId)
		if err != nil {
			return err
		}

		owner, err := pq.GetUser(ctx, tab.OwnerId)
		if err != nil {
			return err
		}

		var bill *models.Bill
		for _, b := range tab.Bills {
			if b.Id == billId {
				bill = &b
			}
		}

		events.Dispatch(h.eventDispatcher, events.TabBillPaidEvent{Shop: shop, Tab: tab, Bill: bill, TabOwner: owner})
		return nil
	})
}

func (h *Handler) GetTabsForUser(ctx context.Context, session *sessions.AuthedSession, userId string) (tabs []models.TabOverview, err error) {
	if session.UserId != userId {
		return nil, services.NewUnauthorizedServiceError(nil)
	}

	return db.WithTxRet(ctx, db.PgxConn(h.store), func(pq *db.PgxQueries) ([]models.TabOverview, error) {
		query := models.GetTabsQueryParams{
			OwnerId: &userId,
		}
		return pq.GetTabs(ctx, &query)
	})
}
func (h *Handler) GetTabsForShop(ctx context.Context, session *sessions.AuthedSession, shopId int) (tabs []models.TabOverview, err error) {
	err = WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_READ_TABS, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		query := models.GetTabsQueryParams{
			ShopId: &shopId,
		}
		tabs, err = pq.GetTabs(ctx, &query)
		return err
	})
	return tabs, err
}

func (h *Handler) GetTabById(ctx context.Context, session *sessions.AuthedSession, shopId int, tabId int) (t *models.Tab, err error) {
	err = WithAuthorizeTabAction(ctx, h.store, session, shopId, tabId, authorization.TAB_ACTION_READ, func(pq *db.PgxQueries, user *models.User, shop *models.Shop, tab *models.Tab) error {
		t = tab
		return nil
	})
	return t, err
}

func (h *Handler) AddOrderToTab(ctx context.Context, session *sessions.AuthedSession, shopId int, tabId int, data *models.BillOrderCreate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return WithAuthorizeTabAction(ctx, h.store, session, shopId, tabId, authorization.TAB_ACTION_ADD_ORDER, func(pq *db.PgxQueries, user *models.User, shop *models.Shop, tab *models.Tab) error {
		return pq.AddOrderToTab(ctx, shopId, tabId, data)
	})
}

func (h *Handler) RemoveOrderFromTab(ctx context.Context, session *sessions.AuthedSession, shopId int, tabId int, data *models.BillOrderCreate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return WithAuthorizeTabAction(ctx, h.store, session, shopId, tabId, authorization.TAB_ACTION_REMOVE_ORDER, func(pq *db.PgxQueries, user *models.User, shop *models.Shop, tab *models.Tab) error {
		return pq.RemoveOrderFromTab(ctx, shopId, tabId, data)
	})
}
func WithAuthorizeTabAction(ctx context.Context, conn db.PgxConn, session *sessions.AuthedSession, shopId int, tabId int, action authorization.Action, fn func(pq *db.PgxQueries, user *models.User, shop *models.Shop, tab *models.Tab) error) error {
	return db.WithTx(ctx, conn, func(pq *db.PgxQueries) error {
		user, err := pq.GetUser(ctx, session.UserId)
		if err != nil {
			return err
		}

		shop, err := pq.GetShopById(ctx, shopId)
		if err != nil {
			return err
		}

		tab, err := pq.GetTabById(ctx, shopId, tabId)
		if err != nil {
			return err
		}

		if ok, err := authorization.AuthorizeTabAction(user, &authorization.TabTarget{Shop: shop, Tab: tab}, action); err != nil {
			return err
		} else if !ok {
			return services.NewUnauthorizedServiceError(nil)
		}
		return fn(pq, user, shop, tab)
	})
}
