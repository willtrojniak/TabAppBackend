package user

import (
	"context"
	"log/slog"

	"github.com/willtrojniak/TabAppBackend/db"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services"
	"github.com/willtrojniak/TabAppBackend/services/authorization"
	"github.com/willtrojniak/TabAppBackend/services/sessions"
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

func (h *Handler) GetUser(ctx context.Context, session *sessions.AuthedSession) (*models.User, error) {
	user, err := db.WithTxRet(ctx, h.store, func(q *db.PgxQueries) (*models.User, error) {
		return q.GetUser(ctx, session.UserId)
	})
	if err != nil {
		h.logger.Error("Failed to get user from database", "userId", session.UserId, "err", err)
		return nil, err
	}
	return user, nil

}

func (h *Handler) UpdateUser(ctx context.Context, session *sessions.AuthedSession, userId string, data *models.UserUpdate) error {

	h.logger.Debug("Updating user", "id", userId)

	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	err = db.WithTx(ctx, h.store, func(q *db.PgxQueries) error {
		subject, err := q.GetUser(ctx, session.UserId)
		if err != nil {
			return err
		}
		target, err := q.GetUser(ctx, userId)
		if err != nil {
			return err
		}

		ok, err := authorization.AuthorizeUserAction(subject, target, authorization.USER_ACTION_UPDATE)
		if err != nil {
			return err
		}
		if !ok {
			return services.NewUnauthorizedServiceError(nil)
		}

		return q.UpdateUser(ctx, userId, data)
	})

	if err != nil {
		h.logger.Error("Error updating user", "id", userId, "error", err)
		return err
	}

	h.logger.Debug("Updated user", "id", userId)

	return nil
}
