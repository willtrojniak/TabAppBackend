package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/types"
)

func (h *Handler) CreateLocation(ctx context.Context, session *sessions.Session, data *types.LocationCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, data.ShopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.CreateLocation(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) UpdateLocation(ctx context.Context, session *sessions.Session, shopId int, locationId int, data *types.LocationUpdate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}
	h.logger.Debug("Updating location", "shopId", shopId, "locationId", locationId)

	err = h.store.UpdateLocation(ctx, shopId, locationId, data)
	if err != nil {
		return err
	}
	h.logger.Debug("Updated location", "shopId", shopId, "locationId", locationId)

	return nil
}

func (h *Handler) DeleteLocation(ctx context.Context, session *sessions.Session, shopId int, locationId int) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	h.logger.Debug("Deleting location", "shopId", shopId, "locationId", locationId)

	err = h.store.DeleteLocation(ctx, shopId, locationId)
	if err != nil {
		return err
	}
	h.logger.Debug("Deleted location", "shopId", shopId, "locationId", locationId)

	return nil

}
