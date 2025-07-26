package notifications

import (
	"fmt"

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

type ShopDailyTabReportNotification struct {
	events.DailyTabReportEvent
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

func (n *NotificationService) onDailyTabReport(e events.DailyTabReportEvent) {
	n.NotifyShop(e.Shop, &ShopDailyTabReportNotification{e})
}

func (n *TabRequestNotification) IsDisabledFor(u *models.User, s *models.Shop) bool {
	return !authorization.HasRole(u, s, authorization.ROLE_SHOP_MANAGE_TABS)
}
func (n *TabRequestNotification) SlackChannel(s *models.Shop) string { return s.TabRequestSlackChannel }
func (n *TabRequestNotification) Heading() string {
	return fmt.Sprintf("New Tab Request - %s", n.Tab.DisplayName)
}
func (n *TabRequestNotification) SubHeading() string {
	return fmt.Sprintf("A new tab request has been submitted for %s", n.Shop.Name)
}
func (n *TabRequestNotification) ResourceURL() string {
	return fmt.Sprintf("%s/shops/%v/tabs/%v", env.Envs.UI_URI, n.Shop.Id, n.Tab.Id)
}
func (n *TabRequestNotification) Data() []NotificationData {
	return []NotificationData{
		{Field: "Display Name", Value: n.Tab.DisplayName},
		{Field: "Organization", Value: n.Tab.Organization},
		{Field: "Contact", Value: n.TabOwner.Name},
		{Field: "Contact Email", Value: n.TabOwner.Email},
		{Field: "Start Date", Value: fmt.Sprintf("%s %v, %v", n.Tab.StartDate.Month.String(), n.Tab.StartDate.Day, n.Tab.StartDate.Year)},
		{Field: "End Date", Value: fmt.Sprintf("%s %v, %v", n.Tab.EndDate.Month.String(), n.Tab.EndDate.Day, n.Tab.EndDate.Year)},
		{Field: "Start Time", Value: n.Tab.DailyStartTime.String()},
		{Field: "End Time", Value: n.Tab.DailyEndTime.String()},
	}
}

func (n *TabBillPaidNotification) IsDisabledFor(u *models.User, s *models.Shop) bool {
	return !authorization.HasRole(u, s, authorization.ROLE_SHOP_MANAGE_TABS)
}
func (n *TabBillPaidNotification) SlackChannel(s *models.Shop) string {
	return s.TabBillReceiptSlackChannel
}
func (n *TabBillPaidNotification) Heading() string {
	return fmt.Sprintf("New Tab Bill Receipt")
}
func (n *TabBillPaidNotification) SubHeading() string {
	return fmt.Sprintf("Payment received for %s at %s",
		n.Tab.DisplayName,
		n.Shop.Name)
}
func (n *TabBillPaidNotification) ResourceURL() string {
	return fmt.Sprintf("%s/shops/%v/tabs/%v", env.Envs.UI_URI, n.Shop.Id, n.Tab.Id)
}
func (n *TabBillPaidNotification) Data() []NotificationData {
	return []NotificationData{
		{Field: "Tab", Value: n.Tab.DisplayName},
		{Field: "Total", Value: fmt.Sprintf("$%.2f", n.Bill.Total())},
		{Field: "Bill Start Date", Value: fmt.Sprintf("%s %v, %v", n.Bill.StartDate.Month.String(), n.Bill.StartDate.Day, n.Bill.StartDate.Year)},
		{Field: "Bill End Date", Value: fmt.Sprintf("%s %v, %v", n.Bill.EndDate.Month.String(), n.Bill.EndDate.Day, n.Bill.EndDate.Year)},
	}
}

func (n *ShopDailyTabReportNotification) IsDisabledFor(u *models.User, s *models.Shop) bool {
	return !authorization.HasRole(u, s, authorization.ROLE_SHOP_READ_TABS)
}
func (n *ShopDailyTabReportNotification) SlackChannel(s *models.Shop) string {
	return s.DailyUpdateSlackChannel
}
func (n *ShopDailyTabReportNotification) Heading() string {
	return fmt.Sprintf("Daily Tab Report - %s", n.Shop.Name)
}
func (n *ShopDailyTabReportNotification) SubHeading() string {
	return fmt.Sprintf("%v Active Tabs Today", len(n.Tabs))
}
func (n *ShopDailyTabReportNotification) ResourceURL() string {
	return fmt.Sprintf("%s/shops/%v/tabs", env.Envs.UI_URI, n.Shop.Id)
}
func (n *ShopDailyTabReportNotification) Data() []NotificationData {
	data := make([]NotificationData, len(n.Tabs))
	for i, t := range n.Tabs {
		data[i] = NotificationData{Field: t.DisplayName,
			Value: fmt.Sprintf("%s - %s\nLimit: $%v\nVerification: %s",
				t.DailyStartTime.String(),
				t.DailyEndTime.String(),
				t.DollarLimitPerOrder,
				t.VerificationMethod),
		}
	}
	return data
}
