package shop

import (
	"context"
	"errors"
	"log/slog"

	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
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

	err = h.store.CreateShop(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetShops(ctx context.Context, limit int, offset int) ([]models.ShopOverview, error) {
	shops, err := h.store.GetShops(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return shops, nil
}

func (h *Handler) GetShopsByUserId(ctx context.Context, userId string) ([]models.ShopOverview, error) {
	shops, err := h.store.GetShopsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return shops, err
}

func (h *Handler) GetShopById(ctx context.Context, shopId int) (models.Shop, error) {
	shop, err := h.store.GetShopById(ctx, shopId)
	if err != nil {
		return models.Shop{}, err
	}
	return shop, err
}

func (h *Handler) UpdateShop(ctx context.Context, session *sessions.Session, shopId int, data *models.ShopUpdate) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	h.logger.Debug("Updating Shop", "id", shopId)
	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.UpdateShop(ctx, shopId, data)
	if err != nil {
		h.logger.Debug("Error Updating Shop", "id", shopId, "error", err)
		return err
	}
	h.logger.Debug("Updated Shop", "id", shopId)

	return nil
}

func (h *Handler) DeleteShop(ctx context.Context, session *sessions.Session, shopId int) error {
	err := h.AuthorizeModifyShop(ctx, session, shopId)
	if err != nil {
		return err
	}

	err = h.store.DeleteShop(ctx, shopId)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) AuthorizeModifyShop(ctx context.Context, session *sessions.Session, targetShopId int) error {
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
