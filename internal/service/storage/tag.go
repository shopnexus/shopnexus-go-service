package storage

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
	util "shopnexus-go-service/internal/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *ServiceImpl) GetTag(ctx context.Context, tag string) (model.Tag, error) {
	row, err := r.sqlc.GetTag(ctx, tag)
	if err != nil {
		return model.Tag{}, err
	}

	return model.Tag{
		Tag:         row.Tag,
		Description: row.Description,
	}, nil
}

type ListTagsParams struct {
	model.PaginationParams
	Tag         *string
	Description *string
}

func (r *ServiceImpl) CountTags(ctx context.Context, params ListTagsParams) (int64, error) {
	return r.sqlc.CountTags(ctx, sqlc.CountTagsParams{
		Tag:         *PtrToPgtype(&pgtype.Text{}, params.Tag),
		Description: *PtrToPgtype(&pgtype.Text{}, params.Description),
	})
}

func (r *ServiceImpl) ListTags(ctx context.Context, params ListTagsParams) ([]model.Tag, error) {
	rows, err := r.sqlc.ListTags(ctx, sqlc.ListTagsParams{
		Limit:       params.Limit,
		Offset:      params.Offset(),
		Tag:         *PtrToPgtype(&pgtype.Text{}, params.Tag),
		Description: *PtrToPgtype(&pgtype.Text{}, params.Description),
	})
	if err != nil {
		return nil, err
	}

	tags := make([]model.Tag, 0, len(rows))
	for _, row := range rows {
		tags = append(tags, model.Tag{
			Tag:         row.Tag,
			Description: row.Description,
		})
	}

	return tags, nil
}

func (r *ServiceImpl) CreateTag(ctx context.Context, tag model.Tag) error {
	return r.sqlc.CreateTag(ctx, sqlc.CreateTagParams{
		Tag:         tag.Tag,
		Description: tag.Description,
	})
}

type UpdateTagParams struct {
	Tag         string
	NewTag      *string
	Description *string
}

func (r *ServiceImpl) UpdateTag(ctx context.Context, params UpdateTagParams) error {
	return r.sqlc.UpdateTag(ctx, sqlc.UpdateTagParams{
		Tag:         params.Tag,
		NewTag:      *PtrToPgtype(&pgtype.Text{}, params.NewTag),
		Description: *PtrToPgtype(&pgtype.Text{}, params.Description),
	})
}

func (r *ServiceImpl) DeleteTag(ctx context.Context, tag string) error {
	return r.sqlc.DeleteTag(ctx, tag)
}

func (r *ServiceImpl) CountProductModelsOnTag(ctx context.Context, tag string) (int64, error) {
	return r.sqlc.CountProductModelsOnTag(ctx, tag)
}

func (r *ServiceImpl) GetTags(ctx context.Context, productModelID int64) ([]string, error) {
	return r.sqlc.GetTags(ctx, productModelID)
}

func (r *ServiceImpl) AddTags(ctx context.Context, productModelID int64, tags []string) error {
	return r.sqlc.AddTags(ctx, sqlc.AddTagsParams{
		ProductModelID: productModelID,
		Tags:           tags,
	})
}

func (r *ServiceImpl) RemoveTags(ctx context.Context, productModelID int64, tags []string) error {
	return r.sqlc.RemoveTags(ctx, sqlc.RemoveTagsParams{
		ProductModelID: productModelID,
		Tags:           tags,
	})
}

func (r *ServiceImpl) UpdateTags(ctx context.Context, productModelID int64, tags []string) error {
	current, err := r.GetTags(ctx, productModelID)
	if err != nil {
		return err
	}

	added, removed := util.Diff(current, tags)
	if err := r.AddTags(ctx, productModelID, added); err != nil {
		return err
	}

	if err := r.RemoveTags(ctx, productModelID, removed); err != nil {
		return err
	}

	return nil
}
