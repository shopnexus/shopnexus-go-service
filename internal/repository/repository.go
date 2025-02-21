package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db   pgxutil.DBTX
	sqlc *sqlc.Queries
}

type RepositoryTx struct {
	*Repository
	tx pgx.Tx
}

func NewRepository(db pgxutil.DBTX) *Repository {
	return &Repository{
		db:   db,
		sqlc: sqlc.New(db),
	}
}

func (r *Repository) Begin(ctx context.Context) (*RepositoryTx, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &RepositoryTx{
		Repository: NewRepository(tx),
		tx:         tx,
	}, nil
}

func (r *RepositoryTx) Commit(ctx context.Context) error {
	return wrapError(r.tx.Commit(ctx))
}

func (r *RepositoryTx) Rollback(ctx context.Context) error {
	return wrapError(r.tx.Rollback(ctx))
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	return err
}
