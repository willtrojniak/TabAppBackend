package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/jackc/pgx/v5"
)

func (s *PgxStore) CreateLocation(ctx context.Context, data *models.LocationCreate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	defer tx.Rollback(ctx)
	row := tx.QueryRow(ctx,
		`INSERT INTO locations (shop_id, name) VALUES  (@shopId, @name) RETURNING id`,
		pgx.NamedArgs{
			"shopId": data.ShopId,
			"name":   data.Name,
		})
	var categoryId int
	err = row.Scan(&categoryId)
	if err != nil {
		return handlePgxError(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) UpdateLocation(ctx context.Context, shopId int, locationId int, data *models.LocationUpdate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, `
    UPDATE locations SET name = @name
    WHERE shop_id = @shopId AND id = @locationId`,
		pgx.NamedArgs{
			"name":       data.Name,
			"shopId":     shopId,
			"locationId": locationId,
		})
	if err != nil {
		return handlePgxError(err)
	}

	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) DeleteLocation(ctx context.Context, shopId int, locationId int) error {

	result, err := s.pool.Exec(ctx, `
    DELETE FROM locations 
    WHERE shop_id = @shopId AND id = @locationId`,
		pgx.NamedArgs{
			"shopId":     shopId,
			"locationId": locationId,
		})
	if err != nil {
		return handlePgxError(err)
	}

	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}

	return nil
}
