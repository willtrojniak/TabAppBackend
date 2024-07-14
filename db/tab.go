package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *PgxStore) CreateTab(ctx context.Context, data *types.TabCreate) error {

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `
    INSERT INTO tabs 
      (shop_id, owner_id, payment_method, organization, display_name,
      start_date, end_date, daily_start_time, daily_end_time, active_days_of_wk,
      dollar_limit_per_order, verification_method, payment_details, billing_interval_days) 
    VALUES (@shopId, @ownerId, @paymentMethod, @organization, @displayName,
            @startDate, @endDate, @dailyStartTime, @dailyEndTime, @activeDaysOfWk,
            @dollarLimitPerOrder, @verificationMethod, @paymentDetails, @billingIntervalDays)`,
		pgx.NamedArgs{
			"shopId":              data.ShopId,
			"ownerId":             data.OwnerId,
			"paymentMethod":       data.PaymentMethod,
			"organization":        data.Organization,
			"displayName":         data.DisplayName,
			"startDate":           data.StartDate,
			"endDate":             data.EndDate,
			"dailyStartTime":      data.DailyStartTime,
			"dailyEndTime":        data.DailyEndTime,
			"activeDaysOfWk":      data.ActiveDaysOfWk,
			"dollarLimitPerOrder": data.DollarLimitPerOrder,
			"verificationMethod":  data.VerificationMethod,
			"paymentDetails":      data.PaymentDetails,
			"billingIntervalDays": data.BillingIntervalDays,
		})
	if err != nil {
		return handlePgxError(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) GetTabs(ctx context.Context, shopId *uuid.UUID) ([]types.Tab, error) {
	rows, err := s.pool.Query(ctx, `
    SELECT * FROM tabs
    WHERE shop_id = @shopId`,
		pgx.NamedArgs{
			"shopId": shopId,
		})

	if err != nil {
		return nil, handlePgxError(err)
	}

	tabs, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[types.Tab])
	if err != nil {
		return nil, handlePgxError(err)
	}
	return tabs, nil
}
