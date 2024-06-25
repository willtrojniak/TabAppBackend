package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxStore struct {
  pool *pgxpool.Pool
  
}

func NewPostgresStorage(context context.Context, config *pgxpool.Config) (*PgxStore, error) {
  conn, err := pgxpool.New(context, config.ConnString());
  if err != nil {
    return &PgxStore{}, err;
  }

  pg := &PgxStore{
    pool: conn,
  }

  return pg, nil;
}


