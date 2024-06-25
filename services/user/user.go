package user

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)


type Handler struct{
  store types.UserStore;
  handleError services.HTTPErrorHandler
}

func NewHandler(store types.UserStore, handleError services.HTTPErrorHandler) *Handler {
  return &Handler{
    store: store,
    handleError: handleError,
  };
}

func (h *Handler) CreateUser(context context.Context, data *types.UserCreate) (*uuid.UUID, error) {
  err := h.ValidateUser(data);
  if err != nil {
    return nil, err;
  }

  id, err := h.store.CreateUser(context, data);
  if err != nil {
    return nil, err;
  }
  return id, nil;
}

func (h *Handler) ValidateUser(data *types.UserCreate) error {
  err := types.Validate.Struct(data);

  if err != nil {
    if err, ok := err.(*validator.InvalidValidationError); ok {
      return services.NewInternalServiceError(err);
    }

    errors := make(services.ValidationErrors);
    for _, err := range err.(validator.ValidationErrors) {
      errors[err.Field()] = services.ValidationError{Value: err.Value(), Error: err.Tag()};
    }
    return services.NewValidationServiceError(err, errors);

  }

  return nil;
}
