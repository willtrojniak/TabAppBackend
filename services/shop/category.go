package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/types"
)

func (h *Handler) CreateCategory(ctx context.Context, session *sessions.Session, data *types.CategoryCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, data.ShopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.CreateCategory(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetCategories(ctx context.Context, shopId int) ([]types.Category, error) {

	categories, err := h.store.GetCategories(ctx, shopId)
	if err != nil {
		return nil, err
	}

	return categories, err

}

func (h *Handler) UpdateCategory(ctx context.Context, session *sessions.Session, shopId int, categoryId int, data *types.CategoryUpdate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}
	h.logger.Debug("Updating category", "shopId", shopId, "categoryId", categoryId)

	err = h.store.UpdateCategory(ctx, shopId, categoryId, data)
	if err != nil {
		return err
	}
	h.logger.Debug("Updated category", "shopId", shopId, "categoryId", categoryId)

	return nil
}

func (h *Handler) DeleteCategory(ctx context.Context, session *sessions.Session, shopId int, categoryId int) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	h.logger.Debug("Deleting category", "shopId", shopId, "categoryId", categoryId)

	err = h.store.DeleteCategory(ctx, shopId, categoryId)
	if err != nil {
		return err
	}
	h.logger.Debug("Deleted category", "shopId", shopId, "categoryId", categoryId)

	return nil

}
