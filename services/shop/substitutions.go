package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateSubstitutionGroup(ctx context.Context, session *sessions.Session, data *models.SubstitutionGroupCreate) error {
	err := h.AuthorizeModifyShop(ctx, session, data.ShopId)
	if err != nil {
		return err
	}

	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.CreateSubstitutionGroup(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) UpdateSubstitutionGroup(ctx context.Context, session *sessions.Session, shopId int, substitutionGroupId int, data *models.SubstitutionGroupUpdate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.UpdateSubstitutionGroup(ctx, shopId, substitutionGroupId, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetSubstitutionGroups(ctx context.Context, shopId int) ([]models.SubstitutionGroup, error) {
	groups, err := h.store.GetSubstitutionGroups(ctx, shopId)
	if err != nil {
		return nil, err
	}

	return groups, err

}

func (h *Handler) DeleteSubstitutionGroup(ctx context.Context, session *sessions.Session, shopId int, substitutionGroupId int) error {
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
