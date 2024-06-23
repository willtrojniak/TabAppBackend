package user

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services/auth"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
)


type Handler struct{
  store types.UserStore;
  auth *auth.Handler;
}

func NewHandler(store types.UserStore, auth *auth.Handler) *Handler {
  return &Handler{
    store: store,
    auth: auth,
  };
}

func (h *Handler) CreateUser(context context.Context, data *types.UserCreate) (*uuid.UUID, error) {

  id, err := h.store.CreateUser(context, data);
  if err != nil {
    return nil, err;
  }
  return id, nil;
}
