package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
)

func (h *Handler) CreateSubstitutionGroup(ctx context.Context, session *sessions.Session, data *types.SubstitutionGroupCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, &data.ShopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.CreateSubstitutionGroup(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) UpdateSubstitutionGroup(ctx context.Context, session *sessions.Session, shopId *uuid.UUID, substitutionGroupId *uuid.UUID, data *types.SubstitutionGroupUpdate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.UpdateSubstitutionGroup(ctx, shopId, substitutionGroupId, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetSubstitutionGroups(ctx context.Context, shopId *uuid.UUID) ([]types.SubstitutionGroup, error) {
	groups, err := h.store.GetSubstitutionGroups(ctx, shopId)
	if err != nil {
		return nil, err
	}

	return groups, err

}

func (h *Handler) DeleteSubstitutionGroup(ctx context.Context, session *sessions.Session, shopId *uuid.UUID, substitutionGroupId *uuid.UUID) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = h.store.DeleteSubstitutionGroup(ctx, shopId, substitutionGroupId)
	if err != nil {
		return err
	}

	return nil
}
