package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
)

func (h *Handler) CreateCategory(ctx context.Context, session *sessions.Session, data *types.CategoryCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, &data.ShopId)
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

func (h *Handler) GetCategories(ctx context.Context, shopId *uuid.UUID) ([]types.Category, error) {

	categories, err := h.store.GetCategories(ctx, shopId)
	if err != nil {
		return nil, err
	}

	return categories, err

}
