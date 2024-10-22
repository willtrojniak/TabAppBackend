package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

const (
	ROLE_USER_OWNER         uint32 = 1 << 0
	ROLE_USER_MANAGE_ITEMS  uint32 = 1 << 1
	ROLE_USER_MANAGE_TABS   uint32 = (1 << 2) | ROLE_USER_READ_TABS
	ROLE_USER_MANAGE_ORDERS uint32 = (1 << 3) | ROLE_USER_READ_TABS
	ROLE_USER_READ_TABS     uint32 = 1 << 4
)

func (h *Handler) GetShopUserPermissions(ctx context.Context, session *sessions.Session, shopId int) (uint32, error) {
	userId, err := session.GetUserId()
	if err != nil {
		return 0, err
	}
	perms, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (uint32, error) {
		isOwner, roles, err := pq.GetShopUserPermissions(ctx, shopId, userId)
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

func (h *Handler) InviteUserToShop(ctx context.Context, session *sessions.Session, shopId int, userData *models.ShopUserCreate) error {
	err := models.ValidateData(userData, h.logger)
	if err != nil {
		return err
	}
	return h.WithAuthorize(ctx, session, shopId, ROLE_USER_OWNER, func(pq *db.PgxQueries) error {
		shop, err := pq.GetShopById(ctx, shopId)
		if err != nil {
			return err
		}
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

func (h *Handler) RemoveInviteToShop(ctx context.Context, session *sessions.Session, shopId int, data *models.ShopUserCreate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		authErr := h.Authorize(ctx, session, shopId, ROLE_USER_OWNER, pq)
		if authErr == nil {
			return pq.RemoveUserFromShop(ctx, shopId, data)
		}

		user, err := h.users.GetUser(ctx, session)
		if err != nil {
			return err
		}
		if user.Email != data.Email {
			return services.NewUnauthorizedServiceError(nil)
		}

		return pq.RemoveUserFromShop(ctx, shopId, data)
	})
}

func (h *Handler) AcceptInviteToShop(ctx context.Context, session *sessions.Session, shopId int) error {
	userId, err := session.GetUserId()
	if err != nil {
		return err
	}
	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		return pq.ConfirmShopInvite(ctx, shopId, userId)
	})
}
