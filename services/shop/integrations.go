package shop

import (
	"context"

	"github.com/willtrojniak/TabAppBackend/db"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services/authorization"
	"github.com/willtrojniak/TabAppBackend/services/sessions"
	"golang.org/x/oauth2"
)

func (h *Handler) InstallSlack(ctx context.Context, session *sessions.AuthedSession, shopId int, token *oauth2.Token) error {
	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_INSTALL_SLACK, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.AddShopSlackToken(ctx, shopId, models.Token(token.AccessToken))
	})
}
