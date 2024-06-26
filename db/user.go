package db

import (
	"context"
	"log/slog"

	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/google/uuid"
)

func (s *PgxStore) CreateUser(context context.Context, data *types.UserCreate) (*uuid.UUID, error) {
  row := s.pool.QueryRow(context, `
    INSERT INTO users (id, email, name) VALUES ($1, $2, $3) ON CONFLICT (email) DO UPDATE
      SET name = excluded.name RETURNING (id)`, uuid.New(), data.Email, data.Name);

  id := uuid.UUID{};
  err := row.Scan(&id);
  if err != nil {
    slog.Warn("Error creating user.", "err", err.Error(), "email", data.Email, "name", data.Name);
    return &uuid.UUID{}, err;
  }

  return &id, nil;
}
