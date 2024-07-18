package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/jackc/pgx/v5"
)

func (s *PgxStore) CreateItem(ctx context.Context, data *types.ItemCreate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	defer tx.Rollback(ctx)
	row := tx.QueryRow(ctx,
		`INSERT INTO items (shop_id, name, base_price) VALUES (@shopId, @name, @basePrice) RETURNING id`,
		pgx.NamedArgs{
			"shopId":    data.ShopId,
			"name":      data.Name,
			"basePrice": data.BasePrice,
		})
	var itemId int
	err = row.Scan(&itemId)
	if err != nil {
		return handlePgxError(err)
	}

	err = s.setItemCategories(ctx, tx, data.ShopId, itemId, data.CategoryIds)
	if err != nil {
		return err
	}

	err = s.setItemAddons(ctx, tx, data.ShopId, itemId, data.AddonIds)
	if err != nil {
		return err
	}

	err = s.setItemSubstitutionGroups(ctx, tx, data.ShopId, itemId, data.SubstitutionGroupIds)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) GetItems(ctx context.Context, shopId int) ([]types.ItemOverview, error) {
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

func (s *PgxStore) GetItem(ctx context.Context, shopId int, itemId int) (types.Item, error) {
	rows, err := s.pool.Query(ctx, `
    SELECT items.id, items.name, items.base_price,
    COALESCE(json_agg(item_categories ORDER BY item_categories.name) FILTER (WHERE item_categories.id IS NOT NULL), '[]') AS categories, 
    COALESCE(json_agg(item_variants ORDER BY item_variants.index) FILTER (WHERE item_variants.id IS NOT NULL), '[]') AS variants,
    COALESCE(json_agg(addons_table ORDER BY item_addons.index) FILTER (WHERE addons_table.id IS NOT NULL), '[]') AS addons,
    COALESCE(json_agg(substitution_groups ORDER BY substitution_groups.index) FILTER (WHERE substitution_groups.id IS NOT NULL), '[]') AS substitution_groups
    FROM items
    LEFT JOIN items_to_categories ON items.shop_id = items_to_categories.shop_id AND items.id = items_to_categories.item_id
    LEFT JOIN item_categories ON items_to_categories.shop_id = item_categories.shop_id AND items_to_categories.item_category_id = item_categories.id
    LEFT JOIN item_variants ON items.shop_id = item_variants.shop_id AND items.id = item_variants.item_id
    LEFT JOIN item_addons ON items.id = item_addons.item_id AND items.shop_id = item_addons.shop_id
    LEFT JOIN items AS addons_table ON item_addons.addon_id = addons_table.id AND item_addons.shop_id = addons_table.shop_id
    LEFT JOIN (
      SELECT items_to_item_substitution_groups.item_id, items_to_item_substitution_groups.shop_id, items_to_item_substitution_groups.index, item_substitution_groups.name, items_to_item_substitution_groups.substitution_group_id AS id,
      COALESCE(json_agg(items ORDER BY item_substitution_groups_to_items.index) FILTER (WHERE items.id IS NOT NULL), '[]') AS substitutions
      FROM items_to_item_substitution_groups
      LEFT JOIN item_substitution_groups ON 
        item_substitution_groups.id = items_to_item_substitution_groups.substitution_group_id
        AND item_substitution_groups.shop_id = items_to_item_substitution_groups.shop_id
      LEFT JOIN item_substitution_groups_to_items ON
        items_to_item_substitution_groups.substitution_group_id = item_substitution_groups_to_items.substitution_group_id
        AND items_to_item_substitution_groups.shop_id = item_substitution_groups_to_items.shop_id
      LEFT JOIN items ON
        item_substitution_groups_to_items.item_id = items.id
        AND item_substitution_groups_to_items.shop_id = items.shop_id
      GROUP BY items_to_item_substitution_groups.substitution_group_id, items_to_item_substitution_groups.item_id, items_to_item_substitution_groups.shop_id, items_to_item_substitution_groups.index, item_substitution_groups.name
    ) AS substitution_groups ON items.id = substitution_groups.item_id AND items.shop_id = substitution_groups.shop_id
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

func (s *PgxStore) UpdateItem(ctx context.Context, shopId int, itemId int, data *types.ItemUpdate) error {
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

	err = s.setItemAddons(ctx, tx, shopId, itemId, data.AddonIds)
	if err != nil {
		return err
	}

	err = s.setItemSubstitutionGroups(ctx, tx, shopId, itemId, data.SubstitutionGroupIds)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) DeleteItem(ctx context.Context, shopId int, itemId int) error {
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
	_, err := s.pool.Exec(ctx, `
    INSERT INTO item_variants (shop_id, item_id, id, name, price, index) VALUES (@shopId, @itemId, @id, @name, @price, @index)`,
		pgx.NamedArgs{
			"shopId": data.ShopId,
			"itemId": data.ItemId,
			"name":   data.Name,
			"price":  data.Price,
			"index":  data.Index,
		})

	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) UpdateItemVariant(ctx context.Context, shopId int, itemId int, variantId int, data *types.ItemVariantUpdate) error {
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

func (s *PgxStore) DeleteItemVariant(ctx context.Context, shopId int, itemId int, variantId int) error {
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

func (s *PgxStore) setItemCategories(ctx context.Context, tx pgx.Tx, shopId int, itemId int, categoryIds []int) error {
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

func (s *PgxStore) setCategoryItems(ctx context.Context, tx pgx.Tx, shopId int, categoryId int, itemIds []int) error {
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
    INSERT INTO items_to_categories SELECT * FROM _temp_upsert_items_to_categories ON CONFLICT (shop_id, item_id, item_category_id) DO UPDATE
    SET index = excluded.index`)
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

func (s *PgxStore) setItemAddons(ctx context.Context, tx pgx.Tx, shopId int, itemId int, addonItemIds []int) error {
	_, err := tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_item_addons (LIKE item_addons INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_item_addons"}, []string{"shop_id", "item_id", "addon_id", "index"}, pgx.CopyFromSlice(len(addonItemIds), func(i int) ([]any, error) {
		return []any{shopId, itemId, addonItemIds[i], i}, nil
	}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx, `
    INSERT INTO item_addons SELECT * FROM _temp_upsert_item_addons ON CONFLICT (shop_id, item_id, addon_id) DO UPDATE
    SET index = excluded.index`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM item_addons WHERE item_id = @itemId AND shop_id = @shopId AND NOT (addon_id = ANY (@addonItemIds))`,
		pgx.NamedArgs{
			"shopId":       shopId,
			"itemId":       itemId,
			"addonItemIds": addonItemIds,
		})
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) setItemSubstitutionGroups(ctx context.Context, tx pgx.Tx, shopId int, itemId int, substitutionGroupIds []int) error {
	_, err := tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_items_to_item_substitution_groups (LIKE items_to_item_substitution_groups INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_items_to_item_substitution_groups"},
		[]string{"shop_id", "substitution_group_id", "item_id", "index"}, pgx.CopyFromSlice(len(substitutionGroupIds), func(i int) ([]any, error) {
			return []any{shopId, substitutionGroupIds[i], itemId, i}, nil
		}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx, `
    INSERT INTO items_to_item_substitution_groups SELECT * FROM _temp_upsert_items_to_item_substitution_groups ON CONFLICT (shop_id, substitution_group_id, item_id) DO UPDATE
    SET index = excluded.index`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM items_to_item_substitution_groups WHERE item_id = @itemId AND shop_id = @shopId AND NOT (substitution_group_id = ANY (@substitutionGroupIds))`,
		pgx.NamedArgs{
			"shopId":               shopId,
			"itemId":               itemId,
			"substitutionGroupIds": substitutionGroupIds,
		})
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}
