package db

import (
	"context"
	"errors"
	"time"

	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/jackc/pgx/v5"
)

func (q *PgxQueries) CreateTab(ctx context.Context, data *models.TabCreate, status models.TabStatus) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		row := q.tx.QueryRow(ctx, `
    INSERT INTO tabs 
      (shop_id, owner_id, payment_method, organization, display_name,
      start_date, end_date, daily_start_time, daily_end_time, active_days_of_wk,
      dollar_limit_per_order, verification_method, payment_details, billing_interval_days, status) 
    VALUES (@shopId, @ownerId, @paymentMethod, @organization, @displayName,
            @startDate, @endDate, @dailyStartTime, @dailyEndTime, @activeDaysOfWk,
            @dollarLimitPerOrder, @verificationMethod, @paymentDetails, @billingIntervalDays, @status)
    RETURNING id`,
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
				"status":              status,
			})

		var tabId int
		err := row.Scan(&tabId)
		if err != nil {
			return handlePgxError(err)
		}

		err = q.setTabUsers(ctx, data.ShopId, tabId, data.VerificationList)
		if err != nil {
			return err
		}

		err = q.setTabLocations(ctx, data.ShopId, tabId, data.LocationIds)
		if err != nil {
			return err
		}
		return nil
	})
}

func (q *PgxQueries) UpdateTab(ctx context.Context, shopId int, tabId int, data *models.TabUpdate) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		_, err := q.tx.Exec(ctx, `
    UPDATE tabs SET
      (payment_method, organization, display_name,
      start_date, end_date, daily_start_time, daily_end_time, active_days_of_wk,
      dollar_limit_per_order, verification_method, payment_details, billing_interval_days) 
    = (@paymentMethod, @organization, @displayName,
            @startDate, @endDate, @dailyStartTime, @dailyEndTime, @activeDaysOfWk,
            @dollarLimitPerOrder, @verificationMethod, @paymentDetails, @billingIntervalDays)
    WHERE id = @tabId AND shop_id = @shopId`,
			pgx.NamedArgs{
				"shopId":              shopId,
				"tabId":               tabId,
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

		err = q.setTabUsers(ctx, shopId, tabId, data.VerificationList)
		if err != nil {
			return err
		}

		err = q.setTabLocations(ctx, shopId, tabId, data.LocationIds)
		if err != nil {
			return err
		}
		return nil
	})
}

func (q *PgxQueries) ApproveTab(ctx context.Context, shopId int, tabId int) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		_, err := q.tx.Exec(ctx, `
    UPDATE tabs SET
      payment_method = u.payment_method,
      organization = u.organization,
      display_name = u.display_name,
      start_date = u.start_date,
      end_date = u.end_date,
      daily_start_time = u.daily_start_time,
      daily_end_time = u.daily_end_time,
      active_days_of_wk = u.active_days_of_wk,
      dollar_limit_per_order = u.dollar_limit_per_order,
      verification_method = u.verification_method,
      payment_details = u.payment_details,
      billing_interval_days = u.billing_interval_days
    FROM tab_updates AS u
    WHERE tabs.id = @tabId AND tabs.shop_id = @shopId 
      AND u.shop_id = tabs.shop_id AND u.tab_id = tabs.id`,
			pgx.NamedArgs{
				"shopId": shopId,
				"tabId":  tabId,
			})
		if err != nil {
			return handlePgxError(err)
		}
		result, err := q.tx.Exec(ctx, `
    UPDATE tabs SET
      status = @status
    WHERE tabs.id = @tabId AND tabs.shop_id = @shopId 
    `, pgx.NamedArgs{
			"shopId": shopId,
			"tabId":  tabId,
			"status": models.TAB_STATUS_CONFIRMED,
		})
		if err != nil {
			return handlePgxError(err)
		}

		if result.RowsAffected() == 0 {
			return services.NewNotFoundServiceError(nil)
		}

		_, err = q.tx.Exec(ctx, `
    DELETE FROM tab_updates
    WHERE shop_id = @shopId AND tab_id = @tabId`,
			pgx.NamedArgs{
				"shopId": shopId,
				"tabId":  tabId,
			})
		if err != nil {
			return handlePgxError(err)
		}

		// Copy over new locations
		// 1. Delete old locations
		_, err = q.tx.Exec(ctx, `
    DELETE FROM tab_locations
    WHERE shop_id = @shopId AND tab_id = @tabId
  `,
			pgx.NamedArgs{
				"shopId": shopId,
				"tabId":  tabId,
			})
		if err != nil {
			return handlePgxError(err)
		}

		// 2. Add new locations
		_, err = q.tx.Exec(ctx, `
    INSERT INTO tab_locations (SELECT * FROM tab_update_locations AS u
    WHERE u.shop_id = @shopId AND u.tab_id = @tabId)
    ON CONFLICT DO NOTHING
  `,
			pgx.NamedArgs{
				"shopId": shopId,
				"tabId":  tabId,
			})
		if err != nil {
			return handlePgxError(err)
		}

		// 3. Delete new locations
		_, err = q.tx.Exec(ctx, `
    DELETE FROM tab_update_locations
    WHERE shop_id = @shopId AND tab_id = @tabId
  `,
			pgx.NamedArgs{
				"shopId": shopId,
				"tabId":  tabId,
			})
		if err != nil {
			return handlePgxError(err)
		}

		return nil
	})

}

func (q *PgxQueries) CloseTab(ctx context.Context, shopId int, tabId int) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		result, err := q.tx.Exec(ctx, `
    UPDATE tabs SET
      status = @status
    WHERE tabs.id = @tabId AND tabs.shop_id = @shopId 
    `, pgx.NamedArgs{
			"shopId": shopId,
			"tabId":  tabId,
			"status": models.TAB_STATUS_CLOSED,
		})
		if err != nil {
			return handlePgxError(err)
		}

		if result.RowsAffected() == 0 {
			return services.NewNotFoundServiceError(nil)
		}

		_, err = q.tx.Exec(ctx, `
    DELETE FROM tab_updates
    WHERE shop_id = @shopId AND tab_id = @tabId`,
			pgx.NamedArgs{
				"shopId": shopId,
				"tabId":  tabId,
			})
		if err != nil {
			return handlePgxError(err)
		}

		return nil
	})
}

func (q *PgxQueries) MarkTabBillPaid(ctx context.Context, shopId int, tabId int, billId int) error {
	endDate := models.DateOf(time.Now())
	println(endDate.String())
	result, err := q.tx.Exec(ctx, `
    UPDATE tab_bills 
    SET (is_paid, end_date) = (true, @endDate)
    WHERE shop_id = @shopId AND tab_id = @tabId AND id = @billId 
    `, pgx.NamedArgs{
		"shopId":  shopId,
		"tabId":   tabId,
		"billId":  billId,
		"endDate": endDate,
	})
	if err != nil {
		return handlePgxError(err)
	}
	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}

	return nil
}

func (q *PgxQueries) SetTabUpdates(ctx context.Context, shopId int, tabId int, data *models.TabUpdate) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		_, err := q.tx.Exec(ctx, `
    INSERT INTO tab_updates 
      (shop_id, tab_id, payment_method, organization, display_name,
      start_date, end_date, daily_start_time, daily_end_time, active_days_of_wk,
      dollar_limit_per_order, verification_method, payment_details, billing_interval_days) 
    VALUES (@shopId, @tabId, @paymentMethod, @organization, @displayName,
            @startDate, @endDate, @dailyStartTime, @dailyEndTime, @activeDaysOfWk,
            @dollarLimitPerOrder, @verificationMethod, @paymentDetails, @billingIntervalDays)
    ON CONFLICT (shop_id, tab_id) DO UPDATE SET
      (payment_method, organization, display_name,
      start_date, end_date, daily_start_time, daily_end_time, active_days_of_wk,
      dollar_limit_per_order, verification_method, payment_details, billing_interval_days) 
    = (excluded.payment_method, excluded.organization, excluded.display_name,
      excluded.start_date, excluded.end_date, excluded.daily_start_time, excluded.daily_end_time, excluded.active_days_of_wk,
      excluded.dollar_limit_per_order, excluded.verification_method, excluded.payment_details, excluded.billing_interval_days)`,
			pgx.NamedArgs{
				"shopId":              shopId,
				"tabId":               tabId,
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

		err = q.setTabUsers(ctx, shopId, tabId, data.VerificationList)
		if err != nil {
			return err
		}

		err = q.SetTabUpdateLocations(ctx, shopId, tabId, data.LocationIds)
		if err != nil {
			return err
		}

		err = q.Commit(ctx)
		if err != nil {
			return handlePgxError(err)
		}

		return nil
	})
}

func (q *PgxQueries) GetTabs(ctx context.Context, shopId int) ([]models.TabOverview, error) {
	rows, err := q.tx.Query(ctx, `
    SELECT 
      tabs.*, 
      (SELECT to_jsonb(tab_updates) as pending_updates
       FROM (SELECT tab_updates.*, 
             COALESCE(json_agg(locations.*) FILTER (WHERE locations.id IS NOT NULL), '[]') AS locations
             FROM tab_updates
             LEFT JOIN tab_update_locations ON tab_updates.shop_id = tab_update_locations.shop_id
               AND tab_updates.tab_id = tab_update_locations.shop_id
             LEFT JOIN locations ON locations.shop_id = tabs.shop_id AND locations.id = tab_update_locations.location_id
             WHERE tab_updates.shop_id = tabs.shop_id AND tab_updates.tab_id = tabs.id
             GROUP BY tab_updates.shop_id, tab_updates.tab_id
            ) AS tab_updates
      ) AS pending_updates,
      array_remove(array_agg(tab_users.email), null) as verification_list,
      EXISTS(
        SELECT tab_bills.id
        FROM tab_bills
        WHERE tab_bills.shop_id = tabs.shop_id AND tab_bills.tab_id = tabs.id AND tab_bills.is_paid = FALSE
        LIMIT 1
      ) as is_pending_balance,
      (SELECT COALESCE(json_agg(locations.*) FILTER (WHERE locations.id IS NOT NULL), '[]') AS locations
       FROM locations
       LEFT JOIN tab_locations ON tab_locations.shop_id = locations.shop_id AND tab_locations.location_id = locations.id
       WHERE tab_locations.tab_id = tabs.id
      ) AS locations
    FROM tabs
    LEFT JOIN tab_users ON tabs.shop_id = tab_users.shop_id AND tabs.id = tab_users.tab_id
    WHERE tabs.shop_id = @shopId
    GROUP BY tabs.shop_id, tabs.id
    ORDER BY tabs.display_name, tabs.start_date, tabs.end_date 
    `,
		pgx.NamedArgs{
			"shopId": shopId,
		})

	if err != nil {
		return nil, handlePgxError(err)
	}

	tabs, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.TabOverview])
	if err != nil {
		return nil, handlePgxError(err)
	}
	return tabs, nil
}

func (q *PgxQueries) GetTabById(ctx context.Context, shopId int, tabId int) (models.Tab, error) {
	rows, err := q.tx.Query(ctx, `
    SELECT tabs.*, 
      EXISTS(
        SELECT tab_bills.id
        FROM tab_bills
        WHERE tab_bills.shop_id = tabs.shop_id AND tab_bills.tab_id = tabs.id AND tab_bills.is_paid = FALSE
        LIMIT 1
      ) as is_pending_balance,
      (SELECT to_jsonb(tab_updates) as pending_updates
       FROM (SELECT tab_updates.*, 
             COALESCE(json_agg(locations.*) FILTER (WHERE locations.id IS NOT NULL), '[]') AS locations
             FROM tab_updates
             LEFT JOIN tab_update_locations ON tab_updates.shop_id = tab_update_locations.shop_id
               AND tab_updates.tab_id = tab_update_locations.shop_id
             LEFT JOIN locations ON locations.shop_id = tabs.shop_id AND locations.id = tab_update_locations.location_id
             WHERE tab_updates.shop_id = tabs.shop_id AND tab_updates.tab_id = tabs.id
             GROUP BY tab_updates.shop_id, tab_updates.tab_id
            ) AS tab_updates
      ) AS pending_updates,
      (SELECT COALESCE(json_agg(locations.*) FILTER (WHERE locations.id IS NOT NULL), '[]') AS locations
       FROM locations
       LEFT JOIN tab_locations ON tab_locations.shop_id = locations.shop_id AND tab_locations.location_id = locations.id
       WHERE tab_locations.tab_id = tabs.id
      ) AS locations,
      (SELECT COALESCE(json_agg(tab_bills) FILTER (WHERE tab_bills.id IS NOT NULL), '[]') AS bills
        FROM 
        (SELECT tab_bills.*, 
          (SELECT COALESCE(json_agg(items) FILTER (WHERE items.id IS NOT NULL), '[]') AS items
            FROM
            (SELECT items.*, oi.quantity,
              (SELECT COALESCE(json_agg(variants) FILTER (WHERE variants.id IS NOT NULL), '[]') AS variants
                FROM
                (SELECT iv.*, ov.quantity
                  FROM order_variants AS ov
                  LEFT JOIN item_variants AS iv ON ov.shop_id = iv.shop_id AND iv.item_id = ov.item_id AND iv.id = ov.variant_id
                  WHERE ov.shop_id = oi.shop_id AND ov.tab_id = oi.tab_id AND ov.bill_id = oi.bill_id AND ov.item_id = oi.item_id) AS variants
            ) AS variants
              FROM order_items AS oi
              LEFT JOIN items ON items.shop_id = oi.shop_id AND items.id = oi.item_id
              WHERE oi.shop_id = tab_bills.shop_id AND oi.tab_id = tab_bills.tab_id AND oi.bill_id = tab_bills.id) AS items
          ) AS items
          FROM tab_bills
          WHERE tab_bills.shop_id = tabs.shop_id AND tab_bills.tab_id = tabs.id
          GROUP BY tab_bills.shop_id, tab_bills.tab_id, tab_bills.id
          ) as tab_bills
      ) AS bills,
      (SELECT array_remove(array_agg(tab_users.email), null)
       FROM tab_users
       WHERE tab_users.shop_id = tabs.shop_id AND tab_users.tab_id = tabs.id
      ) AS verification_list
    FROM tabs
    WHERE tabs.shop_id = @shopId AND tabs.id = @tabId
    GROUP BY tabs.shop_id, tabs.id`,
		pgx.NamedArgs{
			"shopId": shopId,
			"tabId":  tabId,
		})

	if err != nil {
		return models.Tab{}, handlePgxError(err)
	}

	tab, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Tab])
	if err != nil {
		return models.Tab{}, handlePgxError(err)
	}
	return tab, nil
}

func (q *PgxQueries) SetTabUsers(ctx context.Context, shopId int, tabId int, emails []string) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		return q.setTabUsers(ctx, shopId, tabId, emails)
	})
}

func (q *PgxQueries) setTabUsers(ctx context.Context, shopId int, tabId int, emails []string) error {
	_, err := q.tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_tab_users (LIKE tab_users INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_tab_users"}, []string{"shop_id", "tab_id", "email"}, pgx.CopyFromSlice(len(emails), func(i int) ([]any, error) {
		return []any{shopId, tabId, emails[i]}, nil
	}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.Exec(ctx, `
    INSERT INTO tab_users SELECT * FROM _temp_upsert_tab_users ON CONFLICT DO NOTHING`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.Exec(ctx,
		`DELETE FROM tab_users WHERE shop_id = @shopId AND tab_id = @tabId AND NOT (email = ANY (@emails))`,
		pgx.NamedArgs{
			"shopId": shopId,
			"tabId":  tabId,
			"emails": emails,
		})
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (q *PgxQueries) setTabLocations(ctx context.Context, shopId int, tabId int, locationIds []int) error {
	_, err := q.tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_tab_locations (LIKE tab_locations INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_tab_locations"}, []string{"shop_id", "tab_id", "location_id"}, pgx.CopyFromSlice(len(locationIds), func(i int) ([]any, error) {
		return []any{shopId, tabId, locationIds[i]}, nil
	}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.Exec(ctx, `
    INSERT INTO tab_locations SELECT * FROM _temp_upsert_tab_locations ON CONFLICT DO NOTHING`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.Exec(ctx,
		`DELETE FROM tab_locations WHERE shop_id = @shopId AND tab_id = @tabId AND NOT (location_id = ANY (@locationIds))`,
		pgx.NamedArgs{
			"shopId":      shopId,
			"tabId":       tabId,
			"locationIds": locationIds,
		})
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (q *PgxQueries) SetTabUpdateLocations(ctx context.Context, shopId int, tabId int, locationIds []int) error {
	_, err := q.tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_tab_update_locations (LIKE tab_update_locations INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_tab_update_locations"}, []string{"shop_id", "tab_id", "location_id"}, pgx.CopyFromSlice(len(locationIds), func(i int) ([]any, error) {
		return []any{shopId, tabId, locationIds[i]}, nil
	}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.Exec(ctx, `
    INSERT INTO tab_update_locations SELECT * FROM _temp_upsert_tab_update_locations ON CONFLICT DO NOTHING`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.Exec(ctx,
		`DELETE FROM tab_update_locations WHERE shop_id = @shopId AND tab_id = @tabId AND NOT (location_id = ANY (@locationIds))`,
		pgx.NamedArgs{
			"shopId":      shopId,
			"tabId":       tabId,
			"locationIds": locationIds,
		})
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (q *PgxQueries) insertBill(ctx context.Context, shopId int, tabId int, startDate models.Date, endDate models.Date) (int, error) {
	var billId int
	println("start ", startDate.String(), "end", endDate.String())

	row := q.tx.QueryRow(ctx, `
    INSERT INTO tab_bills 
    (shop_id, tab_id, start_date, end_date) VALUES (@shopId, @tabId, @startDate, @endDate) RETURNING id
    `, pgx.NamedArgs{
		"shopId":    shopId,
		"tabId":     tabId,
		"startDate": startDate,
		"endDate":   endDate,
	})
	err := row.Scan(&billId)
	if err != nil {
		return 0, handlePgxError(err)
	}

	return billId, nil
}

func (q *PgxQueries) getTargetBill(ctx context.Context, shopId int, tabId int) (int, error) {
	// TODO: Implement and change to get tab overview by id
	tab, err := q.GetTabById(ctx, shopId, tabId)
	if err != nil {
		return 0, err
	}

	rows, _ := q.tx.Query(ctx, `
    SELECT b.id, b.start_date, b.end_date, b.is_paid
    FROM tab_bills b
    LEFT JOIN tabs ON b.shop_id = tabs.shop_id AND b.tab_id = tabs.id
    WHERE b.shop_id = @shopId AND b.tab_id = @tabId
    ORDER BY is_paid, end_date DESC
    `,
		pgx.NamedArgs{
			"shopId": shopId,
			"tabId":  tabId,
		})

	bill, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.BillOverview])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			endDate := tab.EndDate
			if tab.StartDate.AddDays(tab.BillingIntervalDays).Before(endDate.Date) {
				endDate = models.Date{Date: tab.StartDate.AddDays(tab.BillingIntervalDays - 1)}
			}

			id, err := q.insertBill(ctx, shopId, tabId, tab.StartDate, endDate)
			if err != nil {
				return 0, handlePgxError(err)
			}
			return id, nil
		}

		return 0, handlePgxError(err)
	}

	now := time.Now()
	today := models.DateOf(now)

	if !bill.StartDate.After(today.Date) && !bill.EndDate.Before(today.Date) && !bill.IsPaid {
		return bill.Id, nil
	}

	endDate := tab.EndDate
	if bill.EndDate.AddDays(tab.BillingIntervalDays - 1).Before(endDate.Date) {
		endDate = models.Date{Date: bill.EndDate.AddDays(tab.BillingIntervalDays - 1)}
	}

	id, err := q.insertBill(ctx, shopId, tabId, bill.EndDate, endDate)
	if err != nil {
		return 0, handlePgxError(err)
	}
	return id, nil

}

func (q *PgxQueries) AddOrderToTab(ctx context.Context, shopId int, tabId int, data *models.BillOrderCreate) error {
	err := q.updateTabOrders(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
    INSERT INTO order_items SELECT * FROM _temp_upsert_order_items ON CONFLICT (shop_id, tab_id, bill_id, item_id) DO UPDATE
    SET quantity = order_items.quantity + excluded.quantity`)
		if err != nil {
			return handlePgxError(err)
		}
		_, err = tx.Exec(ctx, `
	   INSERT INTO order_variants SELECT * FROM _temp_upsert_order_variants ON CONFLICT (shop_id, tab_id, bill_id, item_id, variant_id) DO UPDATE
	   SET quantity = order_variants.quantity + excluded.quantity`)
		if err != nil {
			return handlePgxError(err)
		}
		return nil
	}, shopId, tabId, data)
	if err != nil {
		return err
	}

	return nil
}

func (q *PgxQueries) RemoveOrderFromTab(ctx context.Context, shopId int, tabId int, data *models.BillOrderCreate) error {
	err := q.updateTabOrders(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
      UPDATE order_items SET
        quantity = order_items.quantity - u.quantity
      FROM _temp_upsert_order_items AS u
      WHERE order_items.shop_id = u.shop_id
        AND order_items.tab_id = u.tab_id 
        AND order_items.bill_id = u.bill_id 
        AND order_items.item_id = u.item_id`)

		if err != nil {
			return handlePgxError(err)
		}
		_, err = tx.Exec(ctx, `
      UPDATE order_variants SET
        quantity = order_variants.quantity - u.quantity
      FROM _temp_upsert_order_variants AS u
      WHERE order_variants.shop_id = u.shop_id
        AND order_variants.tab_id = u.tab_id 
        AND order_variants.bill_id = u.bill_id 
        AND order_variants.item_id = u.item_id
        AND order_variants.variant_id = u.variant_id`)
		if err != nil {
			return handlePgxError(err)
		}
		return nil
	}, shopId, tabId, data)
	if err != nil {
		return err
	}

	return nil
}

func (q *PgxQueries) updateTabOrders(ctx context.Context, updateFn func(pgx.Tx) error, shopId int, tabId int, data *models.BillOrderCreate) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		billId, err := q.getTargetBill(ctx, shopId, tabId)
		if err != nil {
			return err
		}

		_, err = q.tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_order_items (LIKE order_items INCLUDING ALL ) ON COMMIT DROP`)
		if err != nil {
			return handlePgxError(err)
		}
		_, err = q.tx.Exec(ctx, `
	   CREATE TEMPORARY TABLE _temp_upsert_order_variants (LIKE order_variants INCLUDING ALL ) ON COMMIT DROP`)
		if err != nil {
			return handlePgxError(err)
		}

		type itemOrder struct {
			id       int
			quantity int
		}

		type variantOrder struct {
			itemOrder
			variantId int
		}

		itemOrders := make([]itemOrder, 0)
		variantOrders := make([]variantOrder, 0)
		for _, i := range data.Items {
			itemOrders = append(itemOrders, itemOrder{id: i.Id, quantity: *i.Quantity})
			for _, v := range i.Variants {
				variantOrders = append(variantOrders, variantOrder{itemOrder: itemOrder{id: i.Id, quantity: *v.Quantity}, variantId: v.Id})
			}
		}

		_, err = q.tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_order_items"},
			[]string{"shop_id", "tab_id", "bill_id", "item_id", "quantity"}, pgx.CopyFromSlice(len(itemOrders), func(i int) ([]any, error) {
				return []any{shopId, tabId, billId, itemOrders[i].id, itemOrders[i].quantity}, nil
			}))
		if err != nil {
			return handlePgxError(err)
		}

		_, err = q.tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_order_variants"},
			[]string{"shop_id", "tab_id", "bill_id", "item_id", "variant_id", "quantity"}, pgx.CopyFromSlice(len(variantOrders), func(i int) ([]any, error) {
				return []any{shopId, tabId, billId, variantOrders[i].id, variantOrders[i].variantId, variantOrders[i].quantity}, nil
			}))
		if err != nil {
			return handlePgxError(err)
		}

		err = updateFn(q.tx)
		if err != nil {
			return err
		}

		return nil
	})
}
