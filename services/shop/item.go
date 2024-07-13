package shop

import (
	"context"
	"math"

	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
)

func (h *Handler) CreateItem(ctx context.Context, session *sessions.Session, data *types.ItemCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, &data.ShopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
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

func (h *Handler) GetItems(ctx context.Context, shopId *uuid.UUID) ([]types.ItemOverview, error) {
	items, err := h.store.GetItems(ctx, shopId)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (h *Handler) UpdateItem(ctx context.Context, session *sessions.Session, shopId *uuid.UUID, itemId *uuid.UUID, data *types.ItemUpdate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
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

func (h *Handler) GetItem(ctx context.Context, shopId *uuid.UUID, itemId *uuid.UUID) (types.Item, error) {
	item, err := h.store.GetItem(ctx, shopId, itemId)
	if err != nil {
		return types.Item{}, err
	}

	return item, nil
}

func (h *Handler) DeleteItem(ctx context.Context, session *sessions.Session, shopId *uuid.UUID, itemId *uuid.UUID) error {
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

func (h *Handler) CreateItemVariant(ctx context.Context, session *sessions.Session, data *types.ItemVariantCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, &data.ShopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.CreateItemVariant(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) UpdateItemVariant(ctx context.Context, session *sessions.Session, shopId *uuid.UUID, itemId *uuid.UUID, variantId *uuid.UUID, data *types.ItemVariantUpdate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.UpdateItemVariant(ctx, shopId, itemId, variantId, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) DeleteItemVariant(ctx context.Context, session *sessions.Session, shopId *uuid.UUID, itemId *uuid.UUID, variantId *uuid.UUID) error {
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
