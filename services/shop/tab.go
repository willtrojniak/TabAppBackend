package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
)

func (h *Handler) CreateTab(ctx context.Context, session *sessions.Session, data *types.TabCreate) error {

	err := types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.CreateTab(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetTabs(ctx context.Context, session *sessions.Session, shopId *uuid.UUID) ([]types.Tab, error) {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return nil, err
	}

	tabs, err := h.store.GetTabs(ctx, shopId)
	if err != nil {
		return nil, err
	}

	return tabs, nil
}
