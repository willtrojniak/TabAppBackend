package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/authorization"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) InviteUserToShop(ctx context.Context, session *sessions.AuthedSession, shopId int, userData *models.ShopUserCreate) error {
	err := models.ValidateData(userData, h.logger)
	if err != nil {
		return err
	}
	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_INVITE_USER, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		owner, err := pq.GetUser(ctx, shop.OwnerId)
		if err != nil {
			return err
		}

		// Dissallow inviting the shop owner as a member
		if owner.Email == userData.Email {
			return services.NewDataConflictServiceError(nil)
		}

		return pq.AddUserToShop(ctx, shopId, userData)
	})
}

func (h *Handler) RemoveUserFromShop(ctx context.Context, session *sessions.AuthedSession, shopId int, data *models.ShopUserCreate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_REMOVE_USER, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.RemoveUserFromShop(ctx, shopId, data)
	})
}

func (h *Handler) AcceptInviteToShop(ctx context.Context, session *sessions.AuthedSession, shopId int) error {
	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		return pq.ConfirmShopInvite(ctx, shopId, session.UserId)
	})
}
