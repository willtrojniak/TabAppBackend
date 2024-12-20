package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/authorization"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

const (
	ROLE_USER_OWNER         uint32 = 1 << 0
	ROLE_USER_MANAGE_ITEMS  uint32 = 1 << 1
	ROLE_USER_MANAGE_TABS   uint32 = (1 << 2) | ROLE_USER_READ_TABS
	ROLE_USER_MANAGE_ORDERS uint32 = (1 << 3) | ROLE_USER_READ_TABS
	ROLE_USER_READ_TABS     uint32 = 1 << 4
)

func (h *Handler) GetShopUserPermissions(ctx context.Context, session *sessions.AuthedSession, shopId int) (uint32, error) {
	perms, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (uint32, error) {
		isOwner, roles, err := pq.GetShopUserPermissions(ctx, shopId, session.UserId)
		if err != nil {
			return 0, err
		}
		roles = roles << 1

		if isOwner {
			roles = roles | ROLE_USER_OWNER
		}
		return roles, nil
	})

	if err != nil {
		return 0, err
	}

	return perms, nil
}

func (h *Handler) GetShopUsers(ctx context.Context, session *sessions.AuthedSession, shopId int) (users []models.ShopUser, err error) {
	err = WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_READ_USERS, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		users, err = pq.GetShopUsers(ctx, shopId)
		return err
	})

	if err != nil {
		return nil, err
	}
	for i, user := range users {
		user.Roles = user.Roles << 1
		users[i] = user
	}
	return users, nil
}

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
