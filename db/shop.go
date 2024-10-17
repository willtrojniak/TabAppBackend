package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/models"
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

func (q *PgxQueries) GetShops(ctx context.Context, limit int, offset int) ([]models.ShopOverview, error) {
	rows, err := q.tx.Query(ctx,
		`SELECT shops.*, array_remove(array_agg(payment_methods.method), NULL) as payment_methods FROM shops
    LEFT JOIN payment_methods on shops.id = payment_methods.shop_id
    GROUP BY shops.id
    ORDER BY shops.name
    LIMIT @limit OFFSET @offset`,
		pgx.NamedArgs{
			"limit":  limit,
			"offset": offset,
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

func (q *PgxQueries) GetShopsByUserId(ctx context.Context, userId string) ([]models.ShopOverview, error) {
	rows, err := q.tx.Query(ctx,
		`SELECT shops.*, array_remove(array_agg(payment_methods.method), NULL) as payment_methods FROM shops
    LEFT JOIN payment_methods on shops.id = payment_methods.shop_id
    WHERE shops.owner_id = @userId
    GROUP BY shops.id
    ORDER BY shops.name`,
		pgx.NamedArgs{
			"userId": userId,
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
