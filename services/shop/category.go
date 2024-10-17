package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateCategory(ctx context.Context, session *sessions.Session, data *models.CategoryCreate) error {
	return h.WithAuthorizeModifyShop(ctx, session, data.ShopId, func(pq *db.PgxQueries) error {
		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}

		err = pq.CreateCategory(ctx, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func (h *Handler) GetCategories(ctx context.Context, shopId int) ([]models.Category, error) {
	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.Category, error) {
		return pq.GetCategories(ctx, shopId)
	})
}

func (h *Handler) UpdateCategory(ctx context.Context, session *sessions.Session, shopId int, categoryId int, data *models.CategoryUpdate) error {
	return h.WithAuthorizeModifyShop(ctx, session, shopId, func(pq *db.PgxQueries) error {
		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}
		h.logger.Debug("Updating category", "shopId", shopId, "categoryId", categoryId)

		err = pq.UpdateCategory(ctx, shopId, categoryId, data)
		if err != nil {
			return err
		}
		h.logger.Debug("Updated category", "shopId", shopId, "categoryId", categoryId)

		return nil
	})
}

func (h *Handler) DeleteCategory(ctx context.Context, session *sessions.Session, shopId int, categoryId int) error {
	return h.WithAuthorizeModifyShop(ctx, session, shopId, func(pq *db.PgxQueries) error {
		h.logger.Debug("Deleting category", "shopId", shopId, "categoryId", categoryId)

		err := pq.DeleteCategory(ctx, shopId, categoryId)
		if err != nil {
			return err
		}
		h.logger.Debug("Deleted category", "shopId", shopId, "categoryId", categoryId)

		return nil
	})
}
