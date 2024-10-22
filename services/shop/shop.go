package shop

import (
	"context"
	"errors"
	"log/slog"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/services/user"
)

type Handler struct {
	logger      *slog.Logger
	store       *db.PgxStore
	sessions    *sessions.Handler
	users       *user.Handler
	handleError services.HTTPErrorHandler
}

func NewHandler(store *db.PgxStore, sessions *sessions.Handler, userHandler *user.Handler, handleError services.HTTPErrorHandler, logger *slog.Logger) *Handler {
	return &Handler{
		logger:      logger,
		sessions:    sessions,
		store:       store,
		users:       userHandler,
		handleError: handleError,
	}
}

func (h *Handler) CreateShop(ctx context.Context, session *sessions.Session, data *models.ShopCreate) error {
	h.logger.Debug("Creating shop")
	userId, err := session.GetUserId()
	if err != nil {
		return err
	}

	// Data validation
	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	// Authorization
	if data.OwnerId != userId {
		return services.NewUnauthorizedServiceError(errors.New("Attempted to create shop for another user"))
	}

	err = db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		return pq.CreateShop(ctx, data)
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetShops(ctx context.Context, params *models.GetShopsQueryParams) ([]models.ShopOverview, error) {
	if params == nil {
		params = &models.GetShopsQueryParams{
			Offset: 0,
			Limit:  10,
		}
	}

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.ShopOverview, error) {
		return pq.GetShops(ctx, params)
	})
}

func (h *Handler) GetShopById(ctx context.Context, shopId int) (models.Shop, error) {
	shop, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (models.Shop, error) {
		return pq.GetShopById(ctx, shopId)
	})
	if err != nil {
		return models.Shop{}, err
	}
	return shop, err
}

func (h *Handler) UpdateShop(ctx context.Context, session *sessions.Session, shopId int, data *models.ShopUpdate) error {
	return h.WithAuthorize(ctx, session, shopId, ROLE_USER_OWNER, func(pq *db.PgxQueries) error {
		h.logger.Debug("Updating Shop", "id", shopId)

		err := models.ValidateData(data, h.logger)
		if err != nil {
			return err
		}

		err = db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
			return pq.UpdateShop(ctx, shopId, data)
		})
		if err != nil {
			h.logger.Debug("Error Updating Shop", "id", shopId, "error", err)
			return err
		}

		h.logger.Debug("Updated Shop", "id", shopId)
		return nil
	})
}

func (h *Handler) DeleteShop(ctx context.Context, session *sessions.Session, shopId int) error {
	return h.WithAuthorize(ctx, session, shopId, ROLE_USER_OWNER, func(pq *db.PgxQueries) error {
		return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
			return pq.DeleteShop(ctx, shopId)
		})
	})
}

func (h *Handler) Authorize(ctx context.Context, session *sessions.Session, targetShopId int, roles uint32, pq *db.PgxQueries) error {
	userRoles, err := h.GetShopUserPermissions(ctx, session, targetShopId)
	if err != nil {
		return err
	}

	if (userRoles&ROLE_USER_OWNER) == ROLE_USER_OWNER || (userRoles&roles) == roles {
		return nil
	}

	return services.NewUnauthorizedServiceError(errors.New("Unauthorized"))
}

func (h *Handler) WithAuthorize(ctx context.Context, session *sessions.Session, targetShopId int, roles uint32, fn func(*db.PgxQueries) error) error {
	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		err := h.Authorize(ctx, session, targetShopId, roles, pq)
		if err != nil {
			return err
		}
		return fn(pq)
	})
}
