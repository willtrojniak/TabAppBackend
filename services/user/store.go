package user

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
  pool *pgxpool.Pool;
} 

func NewStore(pool *pgxpool.Pool) *store {
  return &store{
    pool: pool,
  };
}

func (s *store) CreateUser(context context.Context, user types.UserCreate) error {
  s.pool.Exec(context, "");
  return nil;
}
