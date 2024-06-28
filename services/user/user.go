package user

import (
	"context"
	"log/slog"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/sessions"
	"github.com/WilliamTrojniak/TabAppBackend/types"
)

type Handler struct {
	logger      *slog.Logger
	store       types.UserStore
	sessions    *sessions.Handler
	handleError services.HTTPErrorHandler
}

func NewHandler(store types.UserStore, sessions *sessions.Handler, handleError services.HTTPErrorHandler, logger *slog.Logger) *Handler {
	return &Handler{
		logger:      logger,
		sessions:    sessions,
		store:       store,
		handleError: handleError,
	}
}

func (h *Handler) CreateUser(context context.Context, data *types.UserCreate) error {
	h.logger.Debug("Creating user.")
	err := types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.CreateUser(context, data)
	if err != nil {
		return err
	}

	h.logger.Debug("User created.")
	return nil
}
