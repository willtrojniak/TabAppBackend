package shop

import (
	"context"
	"math"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateItem(ctx context.Context, session *sessions.Session, data *models.ItemCreate) error {
	return h.WithAuthorizeModifyShop(ctx, session, data.ShopId, func(pq *db.PgxQueries) error {
		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}
		rounded_price := float32(math.Round(float64(*data.BasePrice)*100) / 100)
		data.BasePrice = &rounded_price

		err = pq.CreateItem(ctx, data)
		if err != nil {
			return err
		}
		return nil
	})
}

func (h *Handler) GetItems(ctx context.Context, shopId int) ([]models.ItemOverview, error) {
	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.ItemOverview, error) {
		return pq.GetItems(ctx, shopId)
	})
}

func (h *Handler) UpdateItem(ctx context.Context, session *sessions.Session, shopId int, itemId int, data *models.ItemUpdate) error {
	return h.WithAuthorizeModifyShop(ctx, session, shopId, func(pq *db.PgxQueries) error {
		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}
		rounded_price := float32(math.Round(float64(*data.BasePrice)*100) / 100)
		data.BasePrice = &rounded_price

		err = pq.UpdateItem(ctx, shopId, itemId, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func (h *Handler) GetItem(ctx context.Context, shopId int, itemId int) (models.Item, error) {
	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (models.Item, error) {
		return pq.GetItem(ctx, shopId, itemId)
	})
}

func (h *Handler) DeleteItem(ctx context.Context, session *sessions.Session, shopId int, itemId int) error {
	return h.WithAuthorizeModifyShop(ctx, session, shopId, func(pq *db.PgxQueries) error {
		h.logger.Debug("Deleting item", "id", itemId)
		err := pq.DeleteItem(ctx, shopId, itemId)
		if err != nil {
			return err
		}
		h.logger.Debug("Deleted item", "id", itemId)

		return nil
	})
}

func (h *Handler) CreateItemVariant(ctx context.Context, session *sessions.Session, data *models.ItemVariantCreate) error {
	return h.WithAuthorizeModifyShop(ctx, session, data.ShopId, func(pq *db.PgxQueries) error {
		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}

		err = pq.CreateItemVariant(ctx, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func (h *Handler) UpdateItemVariant(ctx context.Context, session *sessions.Session, shopId int, itemId int, variantId int, data *models.ItemVariantUpdate) error {
	return h.WithAuthorizeModifyShop(ctx, session, shopId, func(pq *db.PgxQueries) error {
		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}

		err = pq.UpdateItemVariant(ctx, shopId, itemId, variantId, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func (h *Handler) DeleteItemVariant(ctx context.Context, session *sessions.Session, shopId int, itemId int, variantId int) error {
	return h.WithAuthorizeModifyShop(ctx, session, shopId, func(pq *db.PgxQueries) error {
		err := pq.DeleteItemVariant(ctx, shopId, itemId, variantId)
		if err != nil {
			return err
		}

		return nil
	})
}
