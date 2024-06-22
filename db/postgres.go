package db

import (
	"github.com/jackc/pgx"
)

type postgres struct {
  conn *pgx.ConnPool
  
}

func NewPostgresStorage(config pgx.ConnPoolConfig) (*postgres, error) {
  conn, err := pgx.NewConnPool(config);
  if err != nil {
    return &postgres{}, err;
  }

  pg := &postgres{
    conn: conn,
  }

  return pg, nil;
}


