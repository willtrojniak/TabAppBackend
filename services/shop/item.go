package shop

import (
	"context"
	"math"

	"github.com/willtrojniak/TabAppBackend/db"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services/authorization"
	"github.com/willtrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateItem(ctx context.Context, session *sessions.AuthedSession, data *models.ItemCreate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}
	rounded_price := float32(math.Round(float64(*data.BasePrice)*100) / 100)
	data.BasePrice = &rounded_price
	return WithAuthorizeShopAction(ctx, h.store, session, data.ShopId, authorization.SHOP_ACTION_CREATE_ITEM, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.CreateItem(ctx, data)
	})
}

func (h *Handler) GetItems(ctx context.Context, session *sessions.AuthedSession, shopId int) (items []models.ItemOverview, err error) {
	err = WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_READ_ITEMS, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		items, err = pq.GetItems(ctx, shopId)
		return err
	})
	return items, err
}

func (h *Handler) UpdateItem(ctx context.Context, session *sessions.AuthedSession, shopId int, itemId int, data *models.ItemUpdate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}
	rounded_price := float32(math.Round(float64(*data.BasePrice)*100) / 100)
	data.BasePrice = &rounded_price

	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_UPDATE_ITEM, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.UpdateItem(ctx, shopId, itemId, data)
	})
}

func (h *Handler) GetItem(ctx context.Context, session *sessions.AuthedSession, shopId int, itemId int) (item *models.Item, err error) {
	err = WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_READ_ITEM, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		item, err = pq.GetItem(ctx, shopId, itemId)
		return err
	})
	return item, err
}

func (h *Handler) DeleteItem(ctx context.Context, session *sessions.AuthedSession, shopId int, itemId int) error {
	h.logger.Debug("Deleting item", "id", itemId)
	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_DELETE_ITEM, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.DeleteItem(ctx, shopId, itemId)
	})
}

func (h *Handler) CreateItemVariant(ctx context.Context, session *sessions.AuthedSession, data *models.ItemVariantCreate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return WithAuthorizeShopAction(ctx, h.store, session, data.ShopId, authorization.SHOP_ACTION_CREATE_VARIANT, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.CreateItemVariant(ctx, data)
	})
}

func (h *Handler) UpdateItemVariant(ctx context.Context, session *sessions.AuthedSession, shopId int, itemId int, variantId int, data *models.ItemVariantUpdate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_UPDATE_VARIANT, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.UpdateItemVariant(ctx, shopId, itemId, variantId, data)
	})
}

func (h *Handler) DeleteItemVariant(ctx context.Context, session *sessions.AuthedSession, shopId int, itemId int, variantId int) error {
	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_DELETE_VARIANT, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.DeleteItemVariant(ctx, shopId, itemId, variantId)
	})
}
