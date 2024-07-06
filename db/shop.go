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
		`SELECT * FROM shops ORDER BY shops.name LIMIT @limit OFFSET @offset`,
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
		`SELECT * FROM shops WHERE shops.id = @shopId`,
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
	_, err := s.pool.Exec(ctx,
		`UPDATE shops SET name = @name WHERE shops.id = @shopId`,
		pgx.NamedArgs{
			"name":   data.Name,
			"shopId": shopId,
		})

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
