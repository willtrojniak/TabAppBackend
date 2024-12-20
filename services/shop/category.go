package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services/authorization"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateCategory(ctx context.Context, session *sessions.AuthedSession, data *models.CategoryCreate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}
	return WithAuthorizeShopAction(ctx, h.store, session, data.ShopId, authorization.SHOP_ACTION_CREATE_CATEGORY, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.CreateCategory(ctx, data)
	})
}

func (h *Handler) GetCategories(ctx context.Context, shopId int) ([]models.Category, error) {
	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.Category, error) {
		return pq.GetCategories(ctx, shopId)
	})
}

func (h *Handler) UpdateCategory(ctx context.Context, session *sessions.AuthedSession, shopId int, categoryId int, data *models.CategoryUpdate) error {
	h.logger.Debug("Updating category", "shopId", shopId, "categoryId", categoryId)
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}
	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_UPDATE_CATEGORY, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.UpdateCategory(ctx, shopId, categoryId, data)
	})
}

func (h *Handler) DeleteCategory(ctx context.Context, session *sessions.AuthedSession, shopId int, categoryId int) error {
	h.logger.Debug("Deleting category", "shopId", shopId, "categoryId", categoryId)
	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_DELETE_CATEGORY, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.DeleteCategory(ctx, shopId, categoryId)
	})
}
