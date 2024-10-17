package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/models"
	"github.com/jackc/pgx/v5"
)

func (q *PgxQueries) CreateUser(ctx context.Context, data *models.UserCreate) (*models.User, error) {
	row, _ := q.tx.Query(ctx, `
    INSERT INTO users (id, email, name) VALUES (@id, @email, @name) ON CONFLICT (id) DO UPDATE
      SET email = excluded.email, name = excluded.name
      RETURNING *`,
		pgx.NamedArgs{
			"id":    data.Id,
			"email": data.Email,
			"name":  data.Name,
		})

	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.User])

	if err != nil {
		return nil, handlePgxError(err)
	}

	return user, nil
}

func (q *PgxQueries) GetUser(ctx context.Context, userId string) (*models.User, error) {
	row, _ := q.tx.Query(ctx,
		`SELECT * FROM users WHERE id = $1`, userId)

	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		return nil, handlePgxError(err)
	}
	return user, nil
}

func (q *PgxQueries) UpdateUser(ctx context.Context, userId string, data *models.UserUpdate) error {
	_, err := q.tx.Exec(ctx, `UPDATE users SET preferred_name = $1 WHERE id = $2`, data.PreferredName, userId)

	if err != nil {
		return handlePgxError(err)
	}

	return nil
}
