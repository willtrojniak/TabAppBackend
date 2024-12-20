package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services/authorization"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateLocation(ctx context.Context, session *sessions.AuthedSession, data *models.LocationCreate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return WithAuthorizeShopAction(ctx, h.store, session, data.ShopId, authorization.SHOP_ACTION_CREATE_LOCATION, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.CreateLocation(ctx, data)
	})
}

func (h *Handler) UpdateLocation(ctx context.Context, session *sessions.AuthedSession, shopId int, locationId int, data *models.LocationUpdate) error {
	h.logger.Debug("Updating location", "shopId", shopId, "locationId", locationId)
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_UPDATE_LOCATION, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.UpdateLocation(ctx, shopId, locationId, data)
	})
}

func (h *Handler) DeleteLocation(ctx context.Context, session *sessions.AuthedSession, shopId int, locationId int) error {
	h.logger.Debug("Deleting location", "shopId", shopId, "locationId", locationId)

	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_DELETE_LOCATION, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.DeleteLocation(ctx, shopId, locationId)
	})
}
