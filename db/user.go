package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/jackc/pgx/v5"
)

func (s *PgxStore) CreateUser(context context.Context, data *types.UserCreate) error {
	_, err := s.pool.Exec(context, `
    INSERT INTO users (id, email, name) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE
      SET email = excluded.email, name = excluded.name`, data.Id, data.Email, data.Name)

	if err != nil {
		return err
	}

	return nil
}

func (s *PgxStore) GetUser(context context.Context, userId string) (*types.User, error) {
	row, _ := s.pool.Query(context,
		`SELECT * FROM users WHERE id = $1`, userId)

	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[types.User])
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *PgxStore) UpdateUser(context context.Context, userId string, data *types.UserUpdate) error {
	_, err := s.pool.Exec(context, `UPDATE users SET preferred_name = $1 WHERE id = $2`, data.PreferredName, userId)

	if err != nil {
		return err
	}

	return nil
}
