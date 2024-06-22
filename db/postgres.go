package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
  conn *pgxpool.Pool
  
}

func NewPostgresStorage(context context.Context, config *pgxpool.Config) (*postgres, error) {
  conn, err := pgxpool.New(context, config.ConnString());
  if err != nil {
    return &postgres{}, err;
  }

  pg := &postgres{
    conn: conn,
  }

  return pg, nil;
}


