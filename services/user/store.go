package user

import (
	"context"
	"fmt"

	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxstore struct {
  pool *pgxpool.Pool;
} 

func NewPGXStore(pool *pgxpool.Pool) *pgxstore {
  return &pgxstore{
    pool: pool,
  };
}

func (s *pgxstore) CreateUser(context context.Context, data *types.UserCreate) (*uuid.UUID, error) {
  row := s.pool.QueryRow(context, `
    INSERT INTO users (id, email, name) VALUES ($1, $2, $3) ON CONFLICT (email) DO UPDATE
      SET name = excluded.name RETURNING (id)`, uuid.New(), data.Email, data.Name);

  id := uuid.UUID{};
  err := row.Scan(&id);
  if err != nil {
    return &uuid.UUID{}, err;
  }
  fmt.Printf("Scanned: %v\n", id.String());

  return &id, nil;
}
