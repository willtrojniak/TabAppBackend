package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/jackc/pgx/v5"
)

func (s *PgxStore) CreateCategory(ctx context.Context, data *types.CategoryCreate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	defer tx.Rollback(ctx)
	row := tx.QueryRow(ctx,
		`INSERT INTO item_categories (shop_id, name, index) VALUES  (@shopId, @name, @index) RETURNING id`,
		pgx.NamedArgs{
			"shopId": data.ShopId,
			"name":   data.Name,
			"index":  data.Index,
		})
	var categoryId int
	err = row.Scan(&categoryId)
	if err != nil {
		return handlePgxError(err)
	}

	err = s.setCategoryItems(ctx, tx, data.ShopId, categoryId, data.ItemIds)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) GetCategories(ctx context.Context, shopId int) ([]types.Category, error) {

	rows, err := s.pool.Query(ctx,
		`SELECT item_categories.*, array_remove(array_agg(items.id ORDER BY items_to_categories.index), null) AS item_ids FROM item_categories
    LEFT JOIN items_to_categories ON item_categories.shop_id = items_to_categories.shop_id AND item_categories.id = items_to_categories.item_category_id
    LEFT JOIN items ON items_to_categories.shop_id = items.shop_id AND items_to_categories.item_id = items.id
    WHERE item_categories.shop_id = @shopId
    GROUP BY item_categories.shop_id, item_categories.id
    ORDER BY item_categories.index, item_categories.name`,
		pgx.NamedArgs{
			"shopId": shopId,
		})
	if err != nil {
		return nil, handlePgxError(err)
	}

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[types.Category])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return categories, nil
}

func (s *PgxStore) UpdateCategory(ctx context.Context, shopId int, categoryId int, data *types.CategoryUpdate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, `
    UPDATE item_categories SET name = @name, index = @index WHERE shop_id = @shopId AND id = @categoryId`,
		pgx.NamedArgs{
			"name":       data.Name,
			"index":      data.Index,
			"shopId":     shopId,
			"categoryId": categoryId,
		})
	if err != nil {
		return handlePgxError(err)
	}

	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}

	err = s.setCategoryItems(ctx, tx, shopId, categoryId, data.ItemIds)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) DeleteCategory(ctx context.Context, shopId int, categoryId int) error {

	result, err := s.pool.Exec(ctx, `
    DELETE FROM item_categories WHERE shop_id = @shopId AND id = @categoryId`,
		pgx.NamedArgs{
			"shopId":     shopId,
			"categoryId": categoryId,
		})
	if err != nil {
		return handlePgxError(err)
	}

	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}

	return nil
}
