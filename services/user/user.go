package user

import (
	"context"
	"errors"
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
	h.logger.Debug("Creating user", "id", data.Id)
	err := types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.CreateUser(context, data)
	if err != nil {
		return err
	}

	h.logger.Debug("Created user", "id", data.Id)
	return nil
}

func (h *Handler) GetUser(context context.Context, session *sessions.Session) (*types.User, error) {
	userId, err := session.GetUserId()
	if err != nil {
		return nil, err
	}

	user, err := h.store.GetUser(context, userId)
	if err != nil {
		h.logger.Error("Failed to get user from database", "userId", userId, "err", err)
		return nil, services.NewInternalServiceError(err)
	}
	return user, nil

}

func (h *Handler) UpdateUser(context context.Context, session *sessions.Session, userId string, data *types.UserUpdate) error {
	err := h.authorizeModifyUser(session, userId)
	if err != nil {
		return err
	}

	h.logger.Debug("Updating user", "id", userId)

	err = types.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = h.store.UpdateUser(context, userId, data)
	if err != nil {
		h.logger.Error("Error updating user", "id", userId, "error", err)
		return services.NewInternalServiceError(err)
	}
	h.logger.Debug("Updated user", "id", userId)

	return nil
}

func (h *Handler) authorizeModifyUser(session *sessions.Session, targetUserId string) error {
	userId, err := session.GetUserId()
	if err != nil {
		return err
	}

	if userId != targetUserId {
		return services.NewUnauthorizedServiceError(errors.New("Attempt to modify other user's data"))
	}

	return nil

}
