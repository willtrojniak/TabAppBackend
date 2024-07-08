package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *PgxStore) CreateCategory(ctx context.Context, data *types.CategoryCreate) error {
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
		`INSERT INTO item_categories (id, shop_id, name, index) VALUES (@id, @shopId, @name, @index)`,
		pgx.NamedArgs{
			"id":     id,
			"shopId": data.ShopId,
			"name":   data.Name,
			"index":  data.Index,
		})
	if err != nil {
		return handlePgxError(err)
	}

	// TODO: Implement setting category items

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) GetCategories(ctx context.Context, shopId *uuid.UUID) ([]types.Category, error) {

	// TODO: Add associated item ids

	rows, err := s.pool.Query(ctx,
		`SELECT item_categories.* FROM item_categories
    WHERE item_categories.shop_id = @shopId
    ORDER BY item_categories.index, item_categories.name`,
		pgx.NamedArgs{
			"shopId": shopId,
		})
	if err != nil {
		return nil, handlePgxError(err)
	}

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Category])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return categories, nil
}

func (s *PgxStore) UpdateCategory(ctx context.Context, shopId *uuid.UUID, categoryId *uuid.UUID, data *types.CategoryUpdate) error {

	result, err := s.pool.Exec(ctx, `
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

	return nil
}

func (s *PgxStore) DeleteCategory(ctx context.Context, shopId *uuid.UUID, categoryId *uuid.UUID) error {

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
