package notifications

import (
	"github.com/slack-go/slack"
	"github.com/willtrojniak/TabAppBackend/models"
)

type SlackDriver struct{}

func NewSlackDriver() *SlackDriver {
	return &SlackDriver{}
}

func (d *SlackDriver) Name() string                         { return "Slack" }
func (d *SlackDriver) IsDisabledFor(user *models.User) bool { return false }
func (d *SlackDriver) NotifyShop(shop *models.Shop, n Notification) error {
	client := slack.New(shop.SlackAccessToken.String())
	_, _, _, err := client.SendMessage("#test", slack.MsgOptionText(n.Subject(), false))
	return err
}
func (d *SlackDriver) NotifyUsers([]*models.User, Notification) error { return nil }
