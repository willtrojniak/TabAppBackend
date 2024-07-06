package db

import (
	"context"
	"errors"
	"net/http"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *PgxStore) CreateShop(ctx context.Context, data *types.ShopCreate) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return services.NewInternalServiceError(err)
	}

	_, err = s.pool.Exec(ctx,
		`INSERT INTO shops (id, owner_id, name) VALUES ($1, $2, $3)`, id, data.OwnerId, data.Name)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) && pgerr.Code == pgerrcode.UniqueViolation {
			return services.NewServiceError(err, http.StatusConflict, nil)
		}
		return services.NewInternalServiceError(err)
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
		return nil, services.NewInternalServiceError(err)
	}

	shops, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Shop])
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}
	return shops, nil

}
