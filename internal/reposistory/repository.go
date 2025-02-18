package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ RepositoryInterface = (*Repository)(nil)

type RepositoryInterface interface {
	WithTx(tx pgx.Tx) *Repository
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Repository struct {
	pool *pgxpool.Pool
	sqlc *sqlc.Queries
}

func NewRepository(db pgxutil.DBTX) *Repository {
	return &Repository{
		sqlc: sqlc.New(db),
	}
}

func (r *Repository) WithTx(tx pgx.Tx) *Repository {
	return NewRepository(tx)
}

func (r *Repository) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	return err
}
