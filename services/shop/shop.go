package shop

import (
	"context"
	"errors"
	"log/slog"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
)

type Handler struct {
	logger      *slog.Logger
	store       *db.PgxStore
	sessions    *sessions.Handler
	handleError services.HTTPErrorHandler
}

func NewHandler(store *db.PgxStore, sessions *sessions.Handler, handleError services.HTTPErrorHandler, logger *slog.Logger) *Handler {
	return &Handler{
		logger:      logger,
		sessions:    sessions,
		store:       store,
		handleError: handleError,
	}
}

func (h *Handler) CreateShop(ctx context.Context, session *sessions.Session, data *types.ShopCreate) error {
	h.logger.Debug("Creating shop")
	userId, err := session.GetUserId()
	if err != nil {
		return err
	}

	if data.OwnerId == "" {
		data.OwnerId = userId
	}

	// Data validation
	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	// Authorization
	if data.OwnerId != userId {
		return services.NewUnauthorizedServiceError(errors.New("Attempted to create shop for another user"))
	}

	err = h.store.CreateShop(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetShops(ctx context.Context, limit int, offset int) ([]types.Shop, error) {
	shops, err := h.store.GetShops(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return shops, nil

}

func (h *Handler) GetShopById(ctx context.Context, shopId uuid.UUID) (types.Shop, error) {
	shop, err := h.store.GetShopById(ctx, shopId)
	if err != nil {
		return types.Shop{}, err
	}
	return shop, err
}

func (h *Handler) UpdateShop(ctx context.Context, session *sessions.Session, data *types.ShopUpdate) error {
	err := h.AuthorizeModifyShop(ctx, session, data.Id)
	if err != nil {
		return err
	}

	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}
	h.logger.Debug("Updating Shop", "id", data.Id)

	err = h.store.UpdateShop(ctx, data)
	if err != nil {
		h.logger.Debug("Error Updating Shop", "id", data.Id, "error", err)
		return err
	}
	h.logger.Debug("Updated Shop", "id", data.Id)

	return nil
}

func (h *Handler) AuthorizeModifyShop(ctx context.Context, session *sessions.Session, targetShopId uuid.UUID) error {
	userId, err := session.GetUserId()
	if err != nil {
		return err
	}

	shop, err := h.store.GetShopById(ctx, targetShopId)
	if err != nil {
		return err
	}

	if shop.OwnerId != userId {
		return services.NewUnauthorizedServiceError(errors.New("Unauthorized"))
	}

	return nil
}
