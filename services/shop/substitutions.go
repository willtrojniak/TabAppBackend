package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateSubstitutionGroup(ctx context.Context, session *sessions.Session, data *models.SubstitutionGroupCreate) error {
	return h.WithAuthorize(ctx, session, data.ShopId, ROLE_USER_MANAGE_ITEMS, func(pq *db.PgxQueries) error {
		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}

		err = pq.CreateSubstitutionGroup(ctx, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func (h *Handler) UpdateSubstitutionGroup(ctx context.Context, session *sessions.Session, shopId int, substitutionGroupId int, data *models.SubstitutionGroupUpdate) error {
	return h.WithAuthorize(ctx, session, shopId, ROLE_USER_MANAGE_ITEMS, func(pq *db.PgxQueries) error {
		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}

		err = pq.UpdateSubstitutionGroup(ctx, shopId, substitutionGroupId, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func (h *Handler) GetSubstitutionGroups(ctx context.Context, shopId int) ([]models.SubstitutionGroup, error) {
	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.SubstitutionGroup, error) {
		return pq.GetSubstitutionGroups(ctx, shopId)
	})
}

func (h *Handler) DeleteSubstitutionGroup(ctx context.Context, session *sessions.Session, shopId int, substitutionGroupId int) error {
	return h.WithAuthorize(ctx, session, shopId, ROLE_USER_MANAGE_ITEMS, func(pq *db.PgxQueries) error {
		return pq.DeleteSubstitutionGroup(ctx, shopId, substitutionGroupId)
	})
}
