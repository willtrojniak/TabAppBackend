package db

import (
	"context"

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
