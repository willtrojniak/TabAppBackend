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

	err = s.setItemCategories(ctx, tx, &data.ShopId, &id, data.CategoryIds)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) GetItems(ctx context.Context, shopId *uuid.UUID) ([]types.ItemOverview, error) {
	rows, err := s.pool.Query(ctx, `
    SELECT items.base_price, items.name, items.id
    FROM items
    WHERE items.shop_id = @shopId`,
		pgx.NamedArgs{
			"shopId": shopId,
		})

	if err != nil {
		return nil, handlePgxError(err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[types.ItemOverview])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return items, nil

}

func (s *PgxStore) GetItem(ctx context.Context, shopId *uuid.UUID, itemId *uuid.UUID) (types.Item, error) {
	rows, err := s.pool.Query(ctx, `
    SELECT items.id, items.name, items.base_price,
    COALESCE(json_agg(item_categories ORDER BY item_categories.name) FILTER (WHERE item_categories.id IS NOT NULL), '[]') AS categories, 
    COALESCE(json_agg(item_variants ORDER BY item_variants.index) FILTER (WHERE item_variants.id IS NOT NULL), '[]') AS variants
    FROM items
    LEFT JOIN items_to_categories ON items.shop_id = items_to_categories.shop_id AND items.id = items_to_categories.item_id
    LEFT JOIN item_categories ON items_to_categories.shop_id = item_categories.shop_id AND items_to_categories.item_category_id = item_categories.id
    LEFT JOIN item_variants ON items.shop_id = item_variants.shop_id AND items.id = item_variants.item_id
    WHERE items.shop_id = @shopId AND items.id = @itemId
    GROUP BY items.shop_id, items.id`,
		pgx.NamedArgs{
			"shopId": shopId,
			"itemId": itemId,
		})

	if err != nil {
		return types.Item{}, handlePgxError(err)
	}

	item, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[types.Item])
	if err != nil {
		return types.Item{}, handlePgxError(err)
	}

	return item, nil

}

func (s *PgxStore) UpdateItem(ctx context.Context, shopId *uuid.UUID, itemId *uuid.UUID, data *types.ItemUpdate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `
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

	err = s.setItemCategories(ctx, tx, shopId, itemId, data.CategoryIds)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
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

func (s *PgxStore) CreateItemVariant(ctx context.Context, data *types.ItemVariantCreate) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return services.NewInternalServiceError(err)
	}

	_, err = s.pool.Exec(ctx, `
    INSERT INTO item_variants (shop_id, item_id, id, name, price, index) VALUES (@shopId, @itemId, @id, @name, @price, @index)`,
		pgx.NamedArgs{
			"shopId": data.ShopId,
			"itemId": data.ItemId,
			"id":     id,
			"name":   data.Name,
			"price":  data.Price,
			"index":  data.Index,
		})

	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) UpdateItemVariant(ctx context.Context, shopId *uuid.UUID, itemId *uuid.UUID, variantId *uuid.UUID, data *types.ItemVariantUpdate) error {
	result, err := s.pool.Exec(ctx, `
    UPDATE item_variants SET (name, price, index) = (@name, @price, @index)
    WHERE id = @id AND item_id = @itemId AND shop_id = @shopId`,
		pgx.NamedArgs{
			"shopId": shopId,
			"itemId": itemId,
			"id":     variantId,
			"name":   data.Name,
			"price":  data.Price,
			"index":  data.Index,
		})

	if err != nil {
		return handlePgxError(err)
	}

	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}

	return nil
}

func (s *PgxStore) DeleteItemVariant(ctx context.Context, shopId *uuid.UUID, itemId *uuid.UUID, variantId *uuid.UUID) error {
	result, err := s.pool.Exec(ctx, `
    DELETE FROM item_variants 
    WHERE id = @id AND item_id = @itemId AND shop_id = @shopId`,
		pgx.NamedArgs{
			"shopId": shopId,
			"itemId": itemId,
			"id":     variantId,
		})

	if err != nil {
		return handlePgxError(err)
	}

	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}

	return nil
}

func (s *PgxStore) setItemCategories(ctx context.Context, tx pgx.Tx, shopId *uuid.UUID, itemId *uuid.UUID, categoryIds []uuid.UUID) error {
	_, err := tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_items_to_categories (LIKE items_to_categories INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_items_to_categories"}, []string{"shop_id", "item_id", "item_category_id", "index"}, pgx.CopyFromSlice(len(categoryIds), func(i int) ([]any, error) {
		return []any{shopId, itemId, categoryIds[i], 0}, nil
	}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx, `
    INSERT INTO items_to_categories SELECT * FROM _temp_upsert_items_to_categories ON CONFLICT DO NOTHING`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM items_to_categories WHERE shop_id = @shopId AND item_id = @itemId AND NOT (item_category_id = ANY (@categories))`,
		pgx.NamedArgs{
			"shopId":     shopId,
			"itemId":     itemId,
			"categories": categoryIds,
		})
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) setCategoryItems(ctx context.Context, tx pgx.Tx, shopId *uuid.UUID, categoryId *uuid.UUID, itemIds []uuid.UUID) error {
	_, err := tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_items_to_categories (LIKE items_to_categories INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_items_to_categories"}, []string{"shop_id", "item_id", "item_category_id", "index"}, pgx.CopyFromSlice(len(itemIds), func(i int) ([]any, error) {
		return []any{shopId, itemIds[i], categoryId, i}, nil
	}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx, `
    INSERT INTO items_to_categories SELECT * FROM _temp_upsert_items_to_categories ON CONFLICT DO NOTHING`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM items_to_categories WHERE shop_id = @shopId AND item_category_id = @categoryId AND NOT (item_id = ANY (@itemIds))`,
		pgx.NamedArgs{
			"shopId":     shopId,
			"itemIds":    itemIds,
			"categoryId": categoryId,
		})
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}
