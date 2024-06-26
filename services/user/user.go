package user

import (
	"context"
	"log/slog"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)


type Handler struct{
  logger *slog.Logger;
  store types.UserStore;
  handleError services.HTTPErrorHandler
}

func NewHandler(store types.UserStore, handleError services.HTTPErrorHandler, logger *slog.Logger) *Handler {
  return &Handler{
    logger: logger,
    store: store,
    handleError: handleError,
  };
}

func (h *Handler) CreateUser(context context.Context, data *types.UserCreate) (*uuid.UUID, error) {
  h.logger.Debug("Creating user.");
  err := h.ValidateUser(data);
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

func (h *Handler) ValidateUser(data *types.UserCreate) error {
  err := types.Validate.Struct(data);

  if err != nil {
    if err, ok := err.(*validator.InvalidValidationError); ok {
      h.logger.Error("Error while attempting to validate user");
      return services.NewInternalServiceError(err);
    }

    errors := make(services.ValidationErrors);
    for _, err := range err.(validator.ValidationErrors) {
      errors[err.Field()] = services.ValidationError{Value: err.Value(), Error: err.Tag()};
    }
    h.logger.Debug("User validation failled...", "errors", errors);
    return services.NewValidationServiceError(err, errors);

  }

  return nil;
}
