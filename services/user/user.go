package user

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

func (h *Handler) CreateUser(ctx context.Context, data *models.UserCreate) (*models.User, error) {
	h.logger.Debug("Creating user", "id", data.Id)
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return nil, err
	}

	user, err := db.WithTxRet(ctx, h.store, func(q *db.PgxQueries) (*models.User, error) {
		return q.CreateUser(ctx, data)
	})
	if err != nil {
		return nil, err
	}

	h.logger.Debug("Created user", "id", data.Id)
	return user, nil
}

func (h *Handler) GetUser(ctx context.Context, session *sessions.Session) (*models.User, error) {
	userId, err := session.GetUserId()
	if err != nil {
		return nil, err
	}

	user, err := db.WithTxRet(ctx, h.store, func(q *db.PgxQueries) (*models.User, error) {
		return q.GetUser(ctx, userId)
	})
	if err != nil {
		h.logger.Error("Failed to get user from database", "userId", userId, "err", err)
		return nil, err
	}
	return user, nil

}

func (h *Handler) UpdateUser(ctx context.Context, session *sessions.Session, userId string, data *models.UserUpdate) error {
	err := h.authorizeModifyUser(session, userId)
	if err != nil {
		return err
	}

	h.logger.Debug("Updating user", "id", userId)

	err = models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = db.WithTx(ctx, h.store, func(q *db.PgxQueries) error {
		return q.UpdateUser(ctx, userId, data)
	})
	if err != nil {
		h.logger.Error("Error updating user", "id", userId, "error", err)
		return err
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
