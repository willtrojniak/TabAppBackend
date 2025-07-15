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
	tab      *models.Tab
	tabOwner *models.User
	shop     *models.Shop
}

type TabBillPaidNotification struct {
	events.TabBillPaidEvent
}

func NewTabRequestNotification(tab *models.Tab, tabOwner *models.User, shop *models.Shop) *TabRequestNotification {
	return &TabRequestNotification{
		tab:      tab,
		tabOwner: tabOwner,
		shop:     shop,
	}
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
	n.Send(to, NewTabRequestNotification(e.Tab, e.TabOwner, e.Shop))
}

func (n *NotificationService) onTabBillPaid(e events.TabBillPaidEvent) {
	to := make([]*models.User, 0, 2)
	for _, user := range e.Shop.Users {
		if user.IsConfirmed && authorization.HasRole(&user.User, e.Shop, authorization.ROLE_SHOP_READ_TABS) {
			to = append(to, &user.User)
		}
	}

	to = append(to, e.TabOwner)
	n.Send(to, &TabBillPaidNotification{e})
}

func (n *TabRequestNotification) IsDisabledFor(*models.User) bool {
	// FIXME: Add column for user to enable tab create notifications
	return false
}

func (n *TabRequestNotification) Subject() string {
	return fmt.Sprintf("New Tab Request - %s", n.tab.DisplayName)
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
		ShopName:          n.shop.Name,
		TabURL:            fmt.Sprintf("%s/shops/%v/tabs/%v", env.Envs.UI_URI, n.shop.Id, n.tab.Id),
		DisplayName:       n.tab.DisplayName,
		RequestingOrg:     n.tab.Organization,
		RequestingContact: n.tabOwner.Name,
		RequestingEmail:   n.tabOwner.Email,
		StartDate:         fmt.Sprintf("%s %v, %v", n.tab.StartDate.Month.String(), n.tab.StartDate.Date.Day, n.tab.StartDate.Year),
		EndDate:           fmt.Sprintf("%s %v, %v", n.tab.EndDate.Month.String(), n.tab.EndDate.Date.Day, n.tab.EndDate.Year),
		StartTime:         n.tab.DailyStartTime.String(),
		EndTime:           n.tab.DailyEndTime.String(),
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

func (n *TabBillPaidNotification) IsDisabledFor(*models.User) bool {
	// FIXME: Add column for user to enable tab create notifications
	return false
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
