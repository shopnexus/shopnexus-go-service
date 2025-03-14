package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
)

func (r *Repository) CreateTag(ctx context.Context, tag model.Tag) error {
	return r.sqlc.CreateTag(ctx, sqlc.CreateTagParams{
		Tag:         tag.Name,
		Description: tag.Description,
	})
}

func (r *Repository) DeleteTag(ctx context.Context, name string) error {
	return r.sqlc.DeleteTag(ctx, name)
}
