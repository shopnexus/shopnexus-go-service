package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
)

type TagResponse struct {
	model.Tag
	ProductCount int32
}

type ListTagsParams struct {
	model.PaginationParams
	Tag         *string
	Description *string
}

func (s *ProductService) GetTag(ctx context.Context, tag string) (TagResponse, error) {
	tagModel, err := s.repo.GetTag(ctx, tag)
	if err != nil {
		return TagResponse{}, err
	}

	count, err := s.repo.CountProductModelsOnTag(ctx, tag)
	if err != nil {
		return TagResponse{}, err
	}

	return TagResponse{
		Tag:          tagModel,
		ProductCount: int32(count),
	}, nil
}

func (s *ProductService) ListTags(ctx context.Context, params ListTagsParams) (result model.PaginateResult[TagResponse], err error) {
	count, err := s.repo.CountTags(ctx, repository.ListTagsParams{
		PaginationParams: params.PaginationParams,
		Tag:              params.Tag,
		Description:      params.Description,
	})
	if err != nil {
		return result, err
	}

	data, err := s.repo.ListTags(ctx, repository.ListTagsParams{
		PaginationParams: params.PaginationParams,
		Tag:              params.Tag,
		Description:      params.Description,
	})
	if err != nil {
		return result, err
	}

	result.Data = make([]TagResponse, 0, len(data))
	for _, d := range data {
		count, err := s.repo.CountProductModelsOnTag(ctx, d.Tag)
		if err != nil {
			return result, err
		}

		result.Data = append(result.Data, TagResponse{
			Tag:          d,
			ProductCount: int32(count),
		})
	}

	return model.PaginateResult[TagResponse]{
		Data:  result.Data,
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
