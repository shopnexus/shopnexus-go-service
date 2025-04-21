package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
)

func (s *ProductService) GetProductModel(ctx context.Context, id int64) (model.ProductModel, error) {
	productModel, err := s.repo.GetProductModel(ctx, id)
	if err != nil {
		return model.ProductModel{}, err
	}

	return productModel, nil
}

func (s *ProductService) GetProductSerialIDs(ctx context.Context, productID int64) ([]string, error) {
	return s.repo.GetProductSerialIDs(ctx, productID)
}

type ListProductModelsParams = repository.ListProductModelsParams

func (s *ProductService) ListProductModels(ctx context.Context, params ListProductModelsParams) (result model.PaginateResult[model.ProductModel], err error) {
	total, err := s.repo.CountProductModels(ctx, params)
	if err != nil {
		return result, err
	}

	productModels, err := s.repo.ListProductModels(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.ProductModel]{
		Data:       productModels,
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
		NextCursor: nil,
	}, nil
}

type CreateProductModelParams struct {
	UserID int64
	model.ProductModel
}

func (s *ProductService) CreateProductModel(ctx context.Context, params CreateProductModelParams) (model.ProductModel, error) {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return model.ProductModel{}, err
	}
	defer txRepo.Rollback(ctx)

	productModel, err := txRepo.CreateProductModel(ctx, params.ProductModel)
	if err != nil {
		return model.ProductModel{}, err
	}

	if err := txRepo.Commit(ctx); err != nil {
		return model.ProductModel{}, err
	}

	return productModel, nil
}

type UpdateProductModelParams = struct {
	ID               int64
	Type             *int64
	BrandID          *int64
	Name             *string
	Description      *string
	ListPrice        *int64
	DateManufactured *int64
	Resources        *[]string
	Tags             *[]string
}

func (s *ProductService) UpdateProductModel(ctx context.Context, params UpdateProductModelParams) error {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	if err = txRepo.UpdateProductModel(ctx, params); err != nil {
		return err
	}

	return txRepo.Commit(ctx)
}

func (s *ProductService) DeleteProductModel(ctx context.Context, id int64) error {
	return s.repo.DeleteProductModel(ctx, id)
}

type ListProductTypesParams = repository.ListProductTypesParams

func (s *ProductService) ListProductTypes(ctx context.Context, params ListProductTypesParams) ([]model.ProductType, error) {
	return s.repo.ListProductTypes(ctx, params)
}
