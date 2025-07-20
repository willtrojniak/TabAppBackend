package shop

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/willtrojniak/TabAppBackend/db"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services"
	"github.com/willtrojniak/TabAppBackend/services/authorization"
	"github.com/willtrojniak/TabAppBackend/services/sessions"
	"golang.org/x/oauth2"
)

func (h *Handler) InstallSlack(ctx context.Context, session *sessions.AuthedSession, shopId int, token *oauth2.Token) error {
	return WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_INSTALL_SLACK, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.AddShopSlackToken(ctx, shopId, models.Token(token.AccessToken))
	})
}

func (h *Handler) GetShopSlackChannels(ctx context.Context, session *sessions.AuthedSession, shopId int) (channels []models.SlackChannel, err error) {
	err = WithAuthorizeShopAction(ctx, h.store, session, shopId, authorization.SHOP_ACTION_READ_SLACK_CHANNELS, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		if !shop.SlackIntegrated {
			return services.NewUnauthorizedServiceError(nil)
		}

		client := slack.New(shop.SlackAccessToken.String())

		var cursor string
		for {
			channelData, next, err := client.GetConversationsContext(ctx, &slack.GetConversationsParameters{Types: []string{"public_channel", "private_channel"}, ExcludeArchived: true, Cursor: cursor})
			cursor = next
			if err != nil {
				h.logger.Warn("Shop.GetShopSlackChannels", "err", err)
				return services.NewInternalServiceError(err)
			}

			for _, c := range channelData {
				if c.IsChannel {
					channels = append(channels, models.SlackChannel{Name: c.Name, IsPrivate: c.IsPrivate})
				}
			}
			if cursor == "" {
				break
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (h *Handler) UpdateShopSlackChannels(ctx context.Context, s *sessions.AuthedSession, shopId int, data *models.ShopSlackDataUpdate) error {
	err := models.ValidateData(data, h.logger)
	if err != nil {
		return err
	}

	return WithAuthorizeShopAction(ctx, h.store, s, shopId, authorization.SHOP_ACTION_UPDATE_SLACK_CHANNELS, func(pq *db.PgxQueries, user *models.User, shop *models.Shop) error {
		return pq.UpdateShopSlackChannels(ctx, shopId, data)
	})
}
