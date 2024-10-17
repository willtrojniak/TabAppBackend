package db

import (
	"context"
	"errors"
	"log/slog"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxConn interface {
	Begin(ctx context.Context) (*PgxQueries, error)
}

type PgxStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, config *pgxpool.Config) (*PgxStore, error) {
	pool, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		return nil, err
	}

	pg := &PgxStore{
		pool: pool,
	}

	return pg, nil
}

func (s *PgxStore) Begin(ctx context.Context) (*PgxQueries, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, handlePgxError(err)
	}
	return &PgxQueries{tx: tx}, nil
}

type PgxQueries struct {
	tx pgx.Tx
}

func (q *PgxQueries) Begin(ctx context.Context) (*PgxQueries, error) {
	tx, err := q.tx.Begin(ctx)
	if err != nil {
		return nil, handlePgxError(err)
	}
	return &PgxQueries{tx: tx}, nil
}

func (q *PgxQueries) Rollback(ctx context.Context) error {
	return q.tx.Rollback(ctx)
}

func (q *PgxQueries) Commit(ctx context.Context) error {
	return q.tx.Commit(ctx)
}

func (q *PgxQueries) WithTx(ctx context.Context, fn func(*PgxQueries) error) error {
	return WithTx(ctx, q, fn)
}

func WithTx(ctx context.Context, conn PgxConn, fn func(*PgxQueries) error) error {
	q, err := conn.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}
	defer q.Rollback(ctx)
	err = fn(q)
	if err != nil {
		return err
	}

	err = q.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}
	return nil
}

func WithTxRet[T any](ctx context.Context, conn PgxConn, fn func(*PgxQueries) (T, error)) (T, error) {
	q, err := conn.Begin(ctx)
	if err != nil {
		return *new(T), handlePgxError(err)
	}
	defer q.Rollback(ctx)
	val, err := fn(q)
	if err != nil {
		return *new(T), err
	}

	err = q.Commit(ctx)
	if err != nil {
		return *new(T), handlePgxError(err)
	}
	return val, nil
}

func handlePgxError(err error) error {
	var pgerr *pgconn.PgError

	if errors.Is(err, pgx.ErrNoRows) {
		return services.NewNotFoundServiceError(err)
	}

	if errors.As(err, &pgerr) && pgerr.Code == pgerrcode.UniqueViolation {
		return services.NewDataConflictServiceError(err)
	}
	slog.Warn("Database operation failed", "err", err)
	return services.NewInternalServiceError(err)
}
