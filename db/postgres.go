package db

import (
	"context"
	"errors"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, config *pgxpool.Config) (*PgxStore, error) {
	pool, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		return nil, err
	}

	pg := &PgxStore{
		pool: pool,
	}

	return pg, nil
}

func handlePgxError(err error) error {
	var pgerr *pgconn.PgError

	if errors.Is(err, pgx.ErrNoRows) {
		return services.NewNotFoundServiceError(err)
	}

	if errors.As(err, &pgerr) && pgerr.Code == pgerrcode.UniqueViolation {
		return services.NewDataConflictServiceError(err)
	}
	return services.NewInternalServiceError(err)
}
