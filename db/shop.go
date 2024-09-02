package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/jackc/pgx/v5"
)

func (s *PgxStore) CreateShop(ctx context.Context, data *types.ShopCreate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	defer tx.Rollback(ctx)
	row := tx.QueryRow(ctx,
		`INSERT INTO shops (owner_id, name) VALUES (@ownerId, @name) RETURNING id`,
		pgx.NamedArgs{
			"ownerId": data.OwnerId,
			"name":    data.Name,
		})
	var shopId int
	err = row.Scan(&shopId)
	if err != nil {
		return handlePgxError(err)
	}

	err = s.setShopPaymentMethods(ctx, tx, shopId, data.PaymentMethods)
	if err != nil {
		return handlePgxError(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) GetShops(ctx context.Context, limit int, offset int) ([]types.Shop, error) {
	rows, err := s.pool.Query(ctx,
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

	shops, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[types.Shop])
	if err != nil {
		return nil, handlePgxError(err)
	}
	return shops, nil

}

func (s *PgxStore) GetShopsByUserId(ctx context.Context, userId string) ([]types.Shop, error) {
	rows, err := s.pool.Query(ctx,
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

	shops, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[types.Shop])
	if err != nil {
		return nil, handlePgxError(err)
	}
	return shops, nil

}

func (s *PgxStore) GetShopById(ctx context.Context, shopId int) (types.Shop, error) {
	row, err := s.pool.Query(ctx,
		`SELECT shops.*, array_remove(array_agg(payment_methods.method), NULL) as payment_methods FROM shops
    LEFT JOIN payment_methods on shops.id = payment_methods.shop_id
    WHERE shops.id = @shopId
    GROUP BY shops.id`,
		pgx.NamedArgs{
			"shopId": shopId,
		})
	if err != nil {
		return types.Shop{}, handlePgxError(err)
	}

	shop, err := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[types.Shop])
	if err != nil {
		return types.Shop{}, handlePgxError(err)
	}

	return shop, nil

}

func (s *PgxStore) UpdateShop(ctx context.Context, shopId int, data *types.ShopUpdate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE shops SET name = @name WHERE shops.id = @shopId`,
		pgx.NamedArgs{
			"name":   data.Name,
			"shopId": shopId,
		})
	if err != nil {
		return handlePgxError(err)
	}

	err = s.setShopPaymentMethods(ctx, tx, shopId, data.PaymentMethods)
	if err != nil {
		return handlePgxError(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) DeleteShop(ctx context.Context, shopId int) error {
	_, err := s.pool.Exec(ctx,
		`DELETE FROM shops WHERE shops.id = @shopId`,
		pgx.NamedArgs{
			"shopId": shopId,
		})
	if err != nil {
		return handlePgxError(err)
	}
	return nil
}

func (s *PgxStore) setShopPaymentMethods(ctx context.Context, tx pgx.Tx, shopId int, methods []string) error {
	_, err := tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_payment_methods (LIKE payment_methods INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_payment_methods"}, []string{"shop_id", "method"}, pgx.CopyFromSlice(len(methods), func(i int) ([]any, error) {
		return []any{shopId, methods[i]}, nil
	}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx, `
    INSERT INTO payment_methods SELECT * FROM _temp_upsert_payment_methods ON CONFLICT DO NOTHING`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx,
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
