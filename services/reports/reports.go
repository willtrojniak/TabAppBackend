package reports

import (
	"context"
	"slices"

	"github.com/willtrojniak/TabAppBackend/db"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services/events"
)

type ReportHandler struct {
	store      *db.PgxStore
	dispatcher *events.EventDispatcher
}

func NewReportHandler(store *db.PgxStore, dispatcher *events.EventDispatcher) *ReportHandler {
	return &ReportHandler{
		store:      store,
		dispatcher: dispatcher,
	}
}

func (rh *ReportHandler) GenerateShopTabOverview(ctx context.Context, shopId int) {
	var shop *models.Shop
	var tabs []models.TabOverview
	err := db.WithTx(ctx, rh.store, func(pq *db.PgxQueries) error {
		var err error
		shop, err = pq.GetShopById(ctx, shopId)
		if err != nil {
			return err
		}

		query := models.GetTabsQueryParams{ShopId: &shopId}
		tabs, err = pq.GetTabs(ctx, &query)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	tabs = slices.DeleteFunc(tabs, func(t models.TabOverview) bool {
		return !t.IsActiveToday()
	})

	slices.SortFunc(tabs, func(t1, t2 models.TabOverview) int {
		return int(t1.DailyStartTime.Minutes() - t2.DailyStartTime.Minutes())
	})

	events.Dispatch(rh.dispatcher, events.DailyTabReportEvent{Shop: shop, Tabs: tabs})
}
