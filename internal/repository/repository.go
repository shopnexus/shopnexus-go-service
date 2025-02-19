package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
)

var _ RepositoryInterface = (*Repository)(nil)

type RepositoryInterface interface {
	Begin(ctx context.Context) (*Repository, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Repository struct {
	db   pgxutil.DBTX
	sqlc *sqlc.Queries
}

func NewRepository(db pgxutil.DBTX) *Repository {
	return &Repository{
		db:   db,
		sqlc: sqlc.New(db),
	}
}

func (r *Repository) Begin(ctx context.Context) (*Repository, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return NewRepository(tx), nil
}

func (r *Repository) Commit(ctx context.Context) error {
	return wrapError(r.db.Commit(ctx))
}

func (r *Repository) Rollback(ctx context.Context) error {
	return wrapError(r.db.Rollback(ctx))
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	return err
}
