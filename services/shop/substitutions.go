package shop

import (
	"context"

	"github.com/willtrojniak/TabAppBackend/db"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services/authorization"
	"github.com/willtrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateSubstitutionGroup(ctx context.Context, session *sessions.AuthedSession, data *models.SubstitutionGroupCreate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return WithAuthorizeShopAction(ctx, h.store, session, data.ShopId, authorization.SHOP_ACTION_CREATE_SUBSTITUTION, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.CreateSubstitutionGroup(ctx, data)
	})
}

func (h *Handler) UpdateSubstitutionGroup(ctx context.Context, session *sessions.AuthedSession, shopId int, substitutionGroupId int, data *models.SubstitutionGroupUpdate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}
	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_UPDATE_SUBSTITUTION, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.UpdateSubstitutionGroup(ctx, shopId, substitutionGroupId, data)
	})
}

func (h *Handler) GetSubstitutionGroups(ctx context.Context, session *sessions.AuthedSession, shopId int) (substitutions []models.SubstitutionGroup, err error) {
	err = WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_READ_SUBSTITUTIONS, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		substitutions, err = pq.GetSubstitutionGroups(ctx, shopId)
		return err
	})
	return substitutions, err
}

func (h *Handler) DeleteSubstitutionGroup(ctx context.Context, session *sessions.AuthedSession, shopId int, substitutionGroupId int) error {
	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_DELETE_SUBSTITUTION, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.DeleteSubstitutionGroup(ctx, shopId, substitutionGroupId)
	})
}
