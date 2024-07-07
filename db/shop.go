package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *PgxStore) CreateShop(ctx context.Context, data *types.ShopCreate) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return services.NewInternalServiceError(err)
	}

	_, err = s.pool.Exec(ctx,
		`INSERT INTO shops (id, owner_id, name) VALUES ($1, $2, $3)`, id, data.OwnerId, data.Name)
	if err != nil {
		return handlePgxError(err)
	}
	return nil
}

func (s *PgxStore) GetShops(ctx context.Context, limit int, offset int) ([]types.Shop, error) {
	// TODO: Dynamically change limit and offset
	rows, err := s.pool.Query(ctx,
		`SELECT shops.*, array_remove(array_agg(payment_methods.method), NULL) as payment_methods FROM shops
    LEFT JOIN payment_methods on shops.id = payment_methods.shop_id
    GROUP BY shops.id
    LIMIT @limit OFFSET @offset`,
		pgx.NamedArgs{
			"limit":  limit,
			"offset": offset,
		})
	if err != nil {
		return nil, handlePgxError(err)
	}

	shops, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Shop])
	if err != nil {
		return nil, handlePgxError(err)
	}
	return shops, nil

}

func (s *PgxStore) GetShopById(ctx context.Context, shopId *uuid.UUID) (types.Shop, error) {
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

	shop, err := pgx.CollectOneRow(row, pgx.RowToStructByName[types.Shop])
	if err != nil {
		return types.Shop{}, handlePgxError(err)
	}

	return shop, nil

}

func (s *PgxStore) UpdateShop(ctx context.Context, shopId *uuid.UUID, data *types.ShopUpdate) error {
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

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"payment_methods"}, []string{"shop_id", "method"}, pgx.CopyFromSlice(len(data.PaymentMethods), func(i int) ([]any, error) {
		return []any{shopId, data.PaymentMethods[i]}, nil
	}))

	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM payment_methods AS p WHERE p.shop_id = @shopId AND NOT (p.method = ANY (@methods))`,
		pgx.NamedArgs{
			"shopId":  shopId,
			"methods": data.PaymentMethods,
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

func (s *PgxStore) DeleteShop(ctx context.Context, shopId *uuid.UUID) error {
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
