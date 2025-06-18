package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services"
)

func (q *PgxQueries) CreateCategory(ctx context.Context, data *models.CategoryCreate) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		row := q.tx.QueryRow(ctx,
			`INSERT INTO item_categories (shop_id, name, index) VALUES  (@shopId, @name, @index) RETURNING id`,
			pgx.NamedArgs{
				"shopId": data.ShopId,
				"name":   data.Name,
				"index":  data.Index,
			})
		var categoryId int
		err := row.Scan(&categoryId)
		if err != nil {
			return handlePgxError(err)
		}

		err = q.setCategoryItems(ctx, data.ShopId, categoryId, data.ItemIds)
		if err != nil {
			return err
		}

		return nil
	})
}

func (q *PgxQueries) GetCategories(ctx context.Context, shopId int) ([]models.Category, error) {
	rows, err := q.tx.Query(ctx,
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

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Category])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return categories, nil
}

func (q *PgxQueries) UpdateCategory(ctx context.Context, shopId int, categoryId int, data *models.CategoryUpdate) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		result, err := q.tx.Exec(ctx, `
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

		err = q.setCategoryItems(ctx, shopId, categoryId, data.ItemIds)
		if err != nil {
			return err
		}

		return nil
	})
}

func (q *PgxQueries) DeleteCategory(ctx context.Context, shopId int, categoryId int) error {
	result, err := q.tx.Exec(ctx, `
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
