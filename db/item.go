package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *PgxStore) CreateItem(ctx context.Context, data *types.ItemCreate) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return services.NewInternalServiceError(err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`INSERT INTO items (id, shop_id, name, base_price) VALUES (@id, @shopId, @name, @basePrice)`,
		pgx.NamedArgs{
			"id":        id,
			"shopId":    data.ShopId,
			"name":      data.Name,
			"basePrice": data.BasePrice,
		})
	if err != nil {
		return handlePgxError(err)
	}

	// TODO: Implement setting item categories, variants, substitions and addons

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) GetItems(ctx context.Context, shopId *uuid.UUID) ([]types.Item, error) {
	rows, err := s.pool.Query(ctx, `
    SELECT * FROM items WHERE shop_id = @shopId`,
		pgx.NamedArgs{
			"shopId": shopId,
		})

	if err != nil {
		return nil, handlePgxError(err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Item])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return items, nil

}

func (s *PgxStore) UpdateItem(ctx context.Context, shopId *uuid.UUID, itemId *uuid.UUID, data *types.ItemUpdate) error {
	result, err := s.pool.Exec(ctx, `
    UPDATE items SET name = @name, base_price = @base_price
    WHERE shop_id = @shopId AND id = @itemId`,
		pgx.NamedArgs{
			"name":       data.Name,
			"base_price": data.BasePrice,
			"shopId":     shopId,
			"itemId":     itemId,
		})

	if err != nil {
		return handlePgxError(err)
	}

	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}

	return nil
}

func (s *PgxStore) DeleteItem(ctx context.Context, shopId *uuid.UUID, itemId *uuid.UUID) error {
	result, err := s.pool.Exec(ctx, `
    DELETE FROM items 
    WHERE shop_id = @shopId AND id = @itemId`,
		pgx.NamedArgs{
			"shopId": shopId,
			"itemId": itemId,
		})
	if err != nil {
		return handlePgxError(err)
	}

	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}

	return nil
}
