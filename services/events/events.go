package events

import "github.com/willtrojniak/TabAppBackend/models"

type Event any

type TabCreateEvent struct {
	Tab      *models.Tab
	TabOwner *models.User
	Shop     *models.Shop
}

type TabUpdateEvent struct {
	Tab      *models.Tab
	TabOwner *models.User
	Shop     *models.Shop
}

type TabApproveEvent struct {
	Tab      *models.Tab
	TabOwner *models.User
	Shop     *models.Shop
}

type TabCloseEvent struct {
	Tab      *models.Tab
	TabOwner *models.User
	Shop     *models.Shop
}

type TabBillPaidEvent struct {
	Bill     *models.Bill
	Tab      *models.Tab
	TabOwner *models.User
	Shop     *models.Shop
}

type DailyTabReportEvent struct {
	Shop *models.Shop
	Tabs []models.TabOverview
}
