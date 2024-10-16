package shop

import (
	"context"
	"math"

	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateItem(ctx context.Context, session *sessions.Session, data *models.ItemCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, data.ShopId)
	if err != nil {
		return err
	}

	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}
	rounded_price := float32(math.Round(float64(*data.BasePrice)*100) / 100)
	data.BasePrice = &rounded_price

	err = h.store.CreateItem(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetItems(ctx context.Context, shopId int) ([]models.ItemOverview, error) {
	items, err := h.store.GetItems(ctx, shopId)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (h *Handler) UpdateItem(ctx context.Context, session *sessions.Session, shopId int, itemId int, data *models.ItemUpdate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}
	rounded_price := float32(math.Round(float64(*data.BasePrice)*100) / 100)
	data.BasePrice = &rounded_price

	err = h.store.UpdateItem(ctx, shopId, itemId, data)
	if err != nil {
		return err
	}

	return nil

}

func (h *Handler) GetItem(ctx context.Context, shopId int, itemId int) (models.Item, error) {
	item, err := h.store.GetItem(ctx, shopId, itemId)
	if err != nil {
		return models.Item{}, err
	}

	return item, nil
}

func (h *Handler) DeleteItem(ctx context.Context, session *sessions.Session, shopId int, itemId int) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}
	h.logger.Debug("Deleting item", "id", itemId)
	err = h.store.DeleteItem(ctx, shopId, itemId)
	if err != nil {
		return err
	}
	h.logger.Debug("Deleted item", "id", itemId)

	return nil
}

func (h *Handler) CreateItemVariant(ctx context.Context, session *sessions.Session, data *models.ItemVariantCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, data.ShopId)
	if err != nil {
		return err
	}

	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.CreateItemVariant(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) UpdateItemVariant(ctx context.Context, session *sessions.Session, shopId int, itemId int, variantId int, data *models.ItemVariantUpdate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.UpdateItemVariant(ctx, shopId, itemId, variantId, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) DeleteItemVariant(ctx context.Context, session *sessions.Session, shopId int, itemId int, variantId int) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = h.store.DeleteItemVariant(ctx, shopId, itemId, variantId)
	if err != nil {
		return err
	}

	return nil
}
