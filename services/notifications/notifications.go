package notifications

import (
	"log/slog"
	"slices"

	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services/events"
)

type Notification interface {
	IsDisabledFor(*models.User) bool
	Subject() string
	HTML() (string, error)
}

type Notifier interface {
	Name() string
	IsDisabledFor(to *models.User) bool
	Send(to []*models.User, n Notification) error
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

	return n
}

func (n *NotificationService) RegisterDriver(d Notifier, enabled bool) {
	n.drivers[d] = enabled
}

func (n *NotificationService) Send(to []*models.User, notification Notification) {
	to = slices.DeleteFunc(to, notification.IsDisabledFor)

	for driver, enabled := range n.drivers {
		if !enabled {
			continue
		}
		driverTo := slices.DeleteFunc(to, driver.IsDisabledFor)

		go func() {
			err := driver.Send(driverTo, notification)
			if err != nil {
				n.logger.Warn("Error while sending notification.", "driver", driver.Name(), "err", err)
			}
		}()
	}
}
