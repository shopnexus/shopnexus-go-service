package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
)

type ListTagsParams struct {
	model.PaginationParams
	Tag         *string
	Description *string
}

func (s *ProductService) GetTag(ctx context.Context, tag string) (model.Tag, error) {
	return s.repo.GetTag(ctx, tag)
}

func (s *ProductService) ListTags(ctx context.Context, params ListTagsParams) (model.PaginateResult[model.Tag], error) {
	count, err := s.repo.CountTags(ctx, repository.ListTagsParams{
		PaginationParams: params.PaginationParams,
		Tag:              params.Tag,
		Description:      params.Description,
	})
	if err != nil {
		return model.PaginateResult[model.Tag]{}, err
	}

	data, err := s.repo.ListTags(ctx, repository.ListTagsParams{
		PaginationParams: params.PaginationParams,
		Tag:              params.Tag,
		Description:      params.Description,
	})
	if err != nil {
		return model.PaginateResult[model.Tag]{}, err
	}

	return model.PaginateResult[model.Tag]{
		Data:  data,
		Total: count,
		Page:  params.Page,
		Limit: params.Limit,
	}, nil
}

func (s *ProductService) CreateTag(ctx context.Context, tag model.Tag) error {
	return s.repo.CreateTag(ctx, tag)
}

type UpdateTagParams struct {
	Tag         string
	NewTag      *string
	Description *string
}

func (s *ProductService) UpdateTag(ctx context.Context, params UpdateTagParams) error {
	return s.repo.UpdateTag(ctx, repository.UpdateTagParams{
		Tag:         params.Tag,
		NewTag:      params.NewTag,
		Description: params.Description,
	})
}

func (s *ProductService) DeleteTag(ctx context.Context, tag string) error {
	return s.repo.DeleteTag(ctx, tag)
}
