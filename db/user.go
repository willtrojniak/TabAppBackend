package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/types"
)

func (s *PgxStore) CreateUser(context context.Context, data *types.UserCreate) error {
  _, err := s.pool.Exec(context, `
    INSERT INTO users (id, email, name) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE
      SET email = excluded.email, name = excluded.name`, data.Id, data.Email, data.Name);

  if err != nil {
    return err;
  }


  return nil;
}
