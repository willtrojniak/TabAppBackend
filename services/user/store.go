package user

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
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

func (s *store) CreateUser(context context.Context, user *types.UserCreate) (*types.User, error) {
  _, err := s.pool.Exec(context, `
    INSERT INTO users (id, email, name) VALUES ($1, $2, $3) ON CONFLICT (email) DO NOTHING`, uuid.New(), user.Email, user.Name);

  if err != nil {
    return &types.User{}, err;
  }

  return &types.User{}, nil;
}
