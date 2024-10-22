package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/jackc/pgx/v5"
)

func (q *PgxQueries) CreateShop(ctx context.Context, data *models.ShopCreate) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		row := q.tx.QueryRow(ctx,
			`INSERT INTO shops (owner_id, name) VALUES (@ownerId, @name) RETURNING id`,
			pgx.NamedArgs{
				"ownerId": data.OwnerId,
				"name":    data.Name,
			})
		var shopId int
		err := row.Scan(&shopId)
		if err != nil {
			return handlePgxError(err)
		}

		err = q.setShopPaymentMethods(ctx, shopId, data.PaymentMethods)
		if err != nil {
			return handlePgxError(err)
		}
		return nil
	})
}

func (q *PgxQueries) GetShops(ctx context.Context, params *models.GetShopsQueryParams) ([]models.ShopOverview, error) {
	if params == nil {
		return nil, services.NewInternalServiceError(nil)
	}

	rows, err := q.tx.Query(ctx,
		`SELECT shops.*, array_remove(array_agg(payment_methods.method), NULL) as payment_methods FROM shops
    LEFT JOIN payment_methods on shops.id = payment_methods.shop_id
    LEFT JOIN shop_users ON shops.id = shop_users.shop_id
    WHERE ((@userId::text is NULL) OR ((@userId::text = shop_users.user_id OR @userId::text = shops.owner_id) = @isMember))
    AND (
      (@pending::boolean is NULL) 
      OR (@userId::text = shops.owner_id AND @pending::boolean = FALSE) 
      OR (@userId::text = shop_users.user_id AND @pending::boolean != shop_users.confirmed))
    GROUP BY shops.id
    ORDER BY shops.name
    LIMIT @limit OFFSET @offset`,
		pgx.NamedArgs{
			"limit":    params.Limit,
			"offset":   params.Offset,
			"pending":  params.IsPending,
			"isMember": params.IsMember,
			"userId":   params.UserId,
		})
	if err != nil {
		return nil, handlePgxError(err)
	}

	shops, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.ShopOverview])
	if err != nil {
		return nil, handlePgxError(err)
	}
	return shops, nil
}

func (q *PgxQueries) GetShopById(ctx context.Context, shopId int) (models.Shop, error) {
	row, err := q.tx.Query(ctx,
		`SELECT shops.*, 
      array_remove(array_agg(payment_methods.method), NULL) as payment_methods,
      (SELECT COALESCE(json_agg(locations.*) FILTER (WHERE locations.id IS NOT NULL), '[]') AS locations
       FROM locations
       WHERE locations.shop_id = shops.id
      ) AS locations
    FROM shops
    LEFT JOIN payment_methods on shops.id = payment_methods.shop_id
    WHERE shops.id = @shopId
    GROUP BY shops.id`,
		pgx.NamedArgs{
			"shopId": shopId,
		})
	if err != nil {
		return models.Shop{}, handlePgxError(err)
	}

	shop, err := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[models.Shop])
	if err != nil {
		return models.Shop{}, handlePgxError(err)
	}

	return shop, nil
}

func (q *PgxQueries) UpdateShop(ctx context.Context, shopId int, data *models.ShopUpdate) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		_, err := q.tx.Exec(ctx,
			`UPDATE shops SET name = @name WHERE shops.id = @shopId`,
			pgx.NamedArgs{
				"name":   data.Name,
				"shopId": shopId,
			})
		if err != nil {
			return handlePgxError(err)
		}

		err = q.setShopPaymentMethods(ctx, shopId, data.PaymentMethods)
		if err != nil {
			return err
		}
		return nil
	})
}

func (q *PgxQueries) DeleteShop(ctx context.Context, shopId int) error {
	_, err := q.tx.Exec(ctx,
		`DELETE FROM shops WHERE shops.id = @shopId`,
		pgx.NamedArgs{
			"shopId": shopId,
		})
	if err != nil {
		return handlePgxError(err)
	}
	return nil
}

func (q *PgxQueries) setShopPaymentMethods(ctx context.Context, shopId int, methods []string) error {
	_, err := q.tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_payment_methods (LIKE payment_methods INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_payment_methods"}, []string{"shop_id", "method"}, pgx.CopyFromSlice(len(methods), func(i int) ([]any, error) {
		return []any{shopId, methods[i]}, nil
	}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.Exec(ctx, `
    INSERT INTO payment_methods SELECT * FROM _temp_upsert_payment_methods ON CONFLICT DO NOTHING`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.Exec(ctx,
		`DELETE FROM payment_methods AS p WHERE p.shop_id = @shopId AND NOT (p.method = ANY (@methods))`,
		pgx.NamedArgs{
			"shopId":  shopId,
			"methods": methods,
		})
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (q *PgxQueries) GetShopUserPermissions(ctx context.Context, shopId int, userId string) (bool, uint32, error) {
	row := q.tx.QueryRow(ctx, `
    SELECT shops.owner_id, 
      (CASE WHEN shop_users.roles IS NULL THEN 0 ELSE shop_users.roles END)
    FROM shops
    LEFT JOIN shop_users ON shops.id = shop_users.shop_id AND shop_users.confirmed = TRUE
    LEFT JOIN users ON shop_users.user_id = users.id
    WHERE shops.id = @shopId AND (shops.owner_id = @userId OR users.id = @userId OR users.id IS NULL)
    `,
		pgx.NamedArgs{
			"shopId": shopId,
			"userId": userId,
		})

	var ownerId string
	var roles uint32

	err := row.Scan(&ownerId, &roles)
	if err != nil {
		return false, 0, handlePgxError(err)
	}
	return ownerId == userId, roles, nil
}

func (q *PgxQueries) AddUserToShop(ctx context.Context, shopId int, user *models.ShopUserCreate) error {
	res, err := q.tx.Exec(ctx, `
		INSERT INTO shop_users (shop_id, user_id, roles, confirmed)
    SELECT @shopId, users.id, @roles, FALSE
    FROM users
    WHERE users.email = @email
    ON CONFLICT (shop_id, user_id) DO UPDATE
    SET roles = excluded.roles,
    updated_at = NOW()
    `,
		pgx.NamedArgs{
			"shopId": shopId,
			"email":  user.Email,
			"roles":  user.Roles,
		})

	if err != nil {
		return handlePgxError(err)
	}

	if res.RowsAffected() == 0 {
		return handlePgxError(pgx.ErrNoRows)
	}

	return nil
}

func (q *PgxQueries) RemoveUserFromShop(ctx context.Context, shopId int, user *models.ShopUserCreate) error {
	res, err := q.tx.Exec(ctx, `
    DELETE FROM shop_users
    USING users
    WHERE shop_users.user_id = users.id AND shop_users.shop_id = @shopId AND users.email = @email
    `,
		pgx.NamedArgs{
			"shopId": shopId,
			"email":  user.Email,
		})
	if err != nil {
		return handlePgxError(err)
	}

	if res.RowsAffected() == 0 {
		return handlePgxError(pgx.ErrNoRows)
	}

	return nil
}

func (q *PgxQueries) ConfirmShopInvite(ctx context.Context, shopId int, userId string) error {
	res, err := q.tx.Exec(ctx, `
    UPDATE shop_users
    SET confirmed = TRUE
    WHERE shop_id = @shopId AND user_id = @userId
    `,
		pgx.NamedArgs{
			"shopId": shopId,
			"userId": userId,
		})

	if err != nil {
		return handlePgxError(err)
	}

	if res.RowsAffected() == 0 {
		return handlePgxError(pgx.ErrNoRows)
	}

	return nil
}

func (pq *PgxQueries) GetShopUsers(ctx context.Context, shopId int) ([]models.ShopUser, error) {
	rows, err := pq.tx.Query(ctx, `
    SELECT shop_users.roles, shop_users.confirmed, shop_users.updated_at, users.*
    FROM shop_users
    LEFT JOIN users ON shop_users.user_id = users.id
    WHERE shop_users.shop_id = @shopId
    ORDER BY users.name
    `, pgx.NamedArgs{
		"shopId": shopId,
	})
	if err != nil {
		return nil, handlePgxError(err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.ShopUser])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return users, nil
}
