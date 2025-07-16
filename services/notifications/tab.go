package notifications

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/willtrojniak/TabAppBackend/env"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services/authorization"
	"github.com/willtrojniak/TabAppBackend/services/events"
)

type TabRequestNotification struct {
	events.TabCreateEvent
}

type TabBillPaidNotification struct {
	events.TabBillPaidEvent
}

func (n *NotificationService) onTabCreate(e events.TabCreateEvent) {
	if e.Tab.Status != models.TAB_STATUS_PENDING.String() {
		return
	}
	to := make([]*models.User, 0)
	for _, user := range e.Shop.Users {
		if user.IsConfirmed && authorization.HasRole(&user.User, e.Shop, authorization.ROLE_SHOP_MANAGE_TABS) {
			to = append(to, &user.User)
		}
	}
	n.NotifyShop(e.Shop, &TabRequestNotification{e})
}

func (n *NotificationService) onTabBillPaid(e events.TabBillPaidEvent) {
	to := make([]*models.User, 0, 2)
	for _, user := range e.Shop.Users {
		if user.IsConfirmed && authorization.HasRole(&user.User, e.Shop, authorization.ROLE_SHOP_READ_TABS) {
			to = append(to, &user.User)
		}
	}

	n.NotifyShop(e.Shop, &TabBillPaidNotification{e})
	n.NotifyUsers([]*models.User{e.TabOwner}, &TabBillPaidNotification{e})
}

func (n *TabRequestNotification) IsDisabledFor(u *models.User, s *models.Shop) bool {
	return !authorization.HasRole(u, s, authorization.ROLE_SHOP_MANAGE_TABS)
}

func (n *TabRequestNotification) Subject() string {
	return fmt.Sprintf("New Tab Request - %s", n.Tab.DisplayName)
}

func (n *TabRequestNotification) HTML() (string, error) {
	const templateName = "resources/templates/notifications/tab-request.html"

	type templateData struct {
		ShopName          string
		TabURL            string
		DisplayName       string
		RequestingOrg     string
		RequestingContact string
		RequestingEmail   string
		StartDate         string
		EndDate           string
		StartTime         string
		EndTime           string
	}

	data := templateData{
		ShopName:          n.Shop.Name,
		TabURL:            fmt.Sprintf("%s/shops/%v/tabs/%v", env.Envs.UI_URI, n.Shop.Id, n.Tab.Id),
		DisplayName:       n.Tab.DisplayName,
		RequestingOrg:     n.Tab.Organization,
		RequestingContact: n.TabOwner.Name,
		RequestingEmail:   n.TabOwner.Email,
		StartDate:         fmt.Sprintf("%s %v, %v", n.Tab.StartDate.Month.String(), n.Tab.StartDate.Date.Day, n.Tab.StartDate.Year),
		EndDate:           fmt.Sprintf("%s %v, %v", n.Tab.EndDate.Month.String(), n.Tab.EndDate.Date.Day, n.Tab.EndDate.Year),
		StartTime:         n.Tab.DailyStartTime.String(),
		EndTime:           n.Tab.DailyEndTime.String(),
	}

	tplate, err := template.ParseFiles(templateName)
	if err != nil {
		return "", err
	}

	var res bytes.Buffer
	err = tplate.Execute(&res, data)
	if err != nil {
		return "", err
	}

	return res.String(), nil

}

func (n *TabBillPaidNotification) IsDisabledFor(u *models.User, s *models.Shop) bool {
	return !authorization.HasRole(u, s, authorization.ROLE_SHOP_MANAGE_TABS)
}

func (n *TabBillPaidNotification) Subject() string {
	return fmt.Sprintf("Tab Receipt - %s - $%.2f",
		n.Tab.DisplayName,
		n.Bill.Total())
}

func (n *TabBillPaidNotification) HTML() (string, error) {
	const templateName = "resources/templates/notifications/tab-bill-paid.html"

	type templateData struct {
		Shop   *models.Shop
		TabURL string
		Tab    *models.Tab
		Bill   *models.Bill
		Total  string
	}

	data := templateData{
		Shop:   n.Shop,
		TabURL: fmt.Sprintf("%s/shops/%v/tabs/%v", env.Envs.UI_URI, n.Shop.Id, n.Tab.Id),
		Tab:    n.Tab,
		Bill:   n.Bill,
		Total:  fmt.Sprintf("%.2f", n.Bill.Total()),
	}

	tplate, err := template.ParseFiles(templateName)
	if err != nil {
		return "", err
	}

	var res bytes.Buffer
	err = tplate.Execute(&res, data)
	if err != nil {
		return "", err
	}

	return res.String(), nil

}
