package shop

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
)

func (h *Handler) CreateLocation(ctx context.Context, session *sessions.Session, data *models.LocationCreate) error {
	return h.WithAuthorize(ctx, session, data.ShopId, ROLE_USER_OWNER, func(pq *db.PgxQueries) error {
		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}

		err = pq.CreateLocation(ctx, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func (h *Handler) UpdateLocation(ctx context.Context, session *sessions.Session, shopId int, locationId int, data *models.LocationUpdate) error {
	return h.WithAuthorize(ctx, session, shopId, ROLE_USER_OWNER, func(pq *db.PgxQueries) error {
		h.logger.Debug("Updating location", "shopId", shopId, "locationId", locationId)
		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}

		err = pq.UpdateLocation(ctx, shopId, locationId, data)
		if err != nil {
			return err
		}
		h.logger.Debug("Updated location", "shopId", shopId, "locationId", locationId)

		return nil
	})
}

func (h *Handler) DeleteLocation(ctx context.Context, session *sessions.Session, shopId int, locationId int) error {
	return h.WithAuthorize(ctx, session, shopId, ROLE_USER_OWNER, func(pq *db.PgxQueries) error {
		h.logger.Debug("Deleting location", "shopId", shopId, "locationId", locationId)

		err := pq.DeleteLocation(ctx, shopId, locationId)
		if err != nil {
			return err
		}
		h.logger.Debug("Deleted location", "shopId", shopId, "locationId", locationId)

		return nil
	})
}
