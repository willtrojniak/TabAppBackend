package notifications

import (
	"log/slog"

	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services/events"
)

type NotificationData struct {
	Field string
	Value string
}

type Notification interface {
	IsDisabledFor(*models.User, *models.Shop) bool
	SlackChannel(*models.Shop) string
	Heading() string
	SubHeading() string
	ResourceURL() string
	Data() []NotificationData
}

type Notifier interface {
	Name() string
	NotifyUsers(to []*models.User, n Notification) error
	NotifyShop(shop *models.Shop, n Notification) error
}

type NotificationService struct {
	logger  *slog.Logger
	drivers map[Notifier]bool
}

func NewNotificationService(l *slog.Logger, e *events.EventDispatcher) *NotificationService {
	n := &NotificationService{
		logger:  l,
		drivers: make(map[Notifier]bool),
	}

	events.Register(e, n.onTabCreate)
	events.Register(e, n.onTabBillPaid)
	events.Register(e, n.onDailyTabReport)

	return n
}

func (n *NotificationService) RegisterDriver(d Notifier, enabled bool) {
	n.drivers[d] = enabled
}

func (n *NotificationService) NotifyShop(shop *models.Shop, notification Notification) {
	for driver, enabled := range n.drivers {
		if !enabled {
			continue
		}

		go func() {
			err := driver.NotifyShop(shop, notification)
			if err != nil {
				n.logger.Warn("Error while sending shop notification.", "driver", driver.Name(), "err", err)
			}
		}()
	}
}

func (n *NotificationService) NotifyUsers(to []*models.User, notification Notification) {
	for driver, enabled := range n.drivers {
		if !enabled {
			continue
		}

		go func() {
			err := driver.NotifyUsers(to, notification)
			if err != nil {
				n.logger.Warn("Error while sending notification.", "driver", driver.Name(), "err", err)
			}
		}()
	}
}
