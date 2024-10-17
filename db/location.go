package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/jackc/pgx/v5"
)

func (q *PgxQueries) CreateLocation(ctx context.Context, data *models.LocationCreate) error {
	row := q.tx.QueryRow(ctx,
		`INSERT INTO locations (shop_id, name) VALUES  (@shopId, @name) RETURNING id`,
		pgx.NamedArgs{
			"shopId": data.ShopId,
			"name":   data.Name,
		})
	var categoryId int
	err := row.Scan(&categoryId)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (q *PgxQueries) UpdateLocation(ctx context.Context, shopId int, locationId int, data *models.LocationUpdate) error {
	result, err := q.tx.Exec(ctx, `
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
	return nil
}

func (q *PgxQueries) DeleteLocation(ctx context.Context, shopId int, locationId int) error {
	result, err := q.tx.Exec(ctx, `
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
