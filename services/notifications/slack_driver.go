package notifications

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/willtrojniak/TabAppBackend/models"
)

type SlackDriver struct{}

func NewSlackDriver() *SlackDriver {
	return &SlackDriver{}
}

func (d *SlackDriver) Name() string { return "Slack" }
func (d *SlackDriver) NotifyShop(shop *models.Shop, n Notification) error {
	channel := fmt.Sprintf("#%s", n.SlackChannel(shop))
	if channel == "#" {
		return nil
	}

	client := slack.New(shop.SlackAccessToken.String())
	_, _, _, err := client.SendMessage(channel, slack.MsgOptionBlocks(d.toMsg(n)...))
	return err
}
func (d *SlackDriver) NotifyUsers([]*models.User, Notification) error { return nil }

func (d *SlackDriver) toMsg(n Notification) []slack.Block {
	// Header
	headerText := slack.NewTextBlockObject("mrkdwn", n.Heading(), false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Fields
	var fields []*slack.TextBlockObject
	for _, kv := range n.Data() {
		fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s:* %s", kv.Field, kv.Value), false, false))
	}
	fieldsSection := slack.NewSectionBlock(nil, fields, nil)

	// Footer
	footerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<%s| View on CafeTrackr>", n.ResourceURL()), false, false)
	footerSection := slack.NewSectionBlock(footerText, nil, nil)

	return []slack.Block{
		headerSection,
		fieldsSection,
		footerSection,
	}
}
