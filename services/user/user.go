package user

import (
	"context"
	"log/slog"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/services/auth"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
)


type Handler struct{
  logger *slog.Logger;
  store types.UserStore;
  auth *auth.Handler;
  handleError services.HTTPErrorHandler
}

func NewHandler(store types.UserStore, auth *auth.Handler, handleError services.HTTPErrorHandler, logger *slog.Logger) *Handler {
  return &Handler{
    logger: logger,
    auth: auth,
    store: store,
    handleError: handleError,
  };
}

func (h *Handler) CreateUser(context context.Context, data *types.UserCreate) (*uuid.UUID, error) {
  h.logger.Debug("Creating user.");
  err := types.ValidateData(data, h.logger);
  if err != nil {
    return nil, err;
  }

  id, err := h.store.CreateUser(context, data);
  if err != nil {
    return nil, err;
  }

  h.logger.Debug("User created.", "id", id);
  return id, nil;
}

