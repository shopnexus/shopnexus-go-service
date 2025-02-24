package service

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
)

type ProductService struct {
	repo    *repository.Repository
	account *AccountService
}

func NewProductService(repo *repository.Repository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetBrand(ctx context.Context, id int64) (model.Brand, error) {
	brand, err := s.repo.GetBrand(ctx, id)
	if err != nil {
		return model.Brand{}, err
	}

	return brand, nil
}

type ListBrandsParams = repository.ListBrandsParams

func (s *ProductService) ListBrands(ctx context.Context, params ListBrandsParams) (result model.PaginateResult[model.Brand], err error) {
	total, err := s.repo.CountBrands(ctx, params)
	if err != nil {
		return result, err
	}

	brands, err := s.repo.ListBrands(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.Brand]{
		Data:       brands,
		Total:      total,
		NextPage:   model.NextPage(params.Offset, params.Limit, total),
		NextCursor: nil,
	}, nil
}

type CreateBrandParams struct {
	UserID int64
	model.Brand
}

func (s *ProductService) CreateBrand(ctx context.Context, params CreateBrandParams) (model.Brand, error) {
	if isAdmin, err := s.account.IsAdmin(ctx, params.UserID); err != nil {
		return model.Brand{}, err
	} else if !isAdmin {
		return model.Brand{}, model.ErrForbidden
	}

	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return model.Brand{}, err
	}
	defer txRepo.Rollback(ctx)

	newBrand, err := txRepo.CreateBrand(ctx, params.Brand)
	if err != nil {
		return model.Brand{}, err
	}

	if err := txRepo.Commit(ctx); err != nil {
		return model.Brand{}, err
	}

	return newBrand, nil
}

type UpdateBrandParams struct {
	repository.UpdateBrandParams
}

func (s *ProductService) UpdateBrand(ctx context.Context, params UpdateBrandParams) error {
	// TODO: chuyển isAdmin này ra chỗ khác
	// if isAdmin, err := s.Account.IsAdmin(ctx, params.UserID); err != nil {
	// 	return err
	// } else if !isAdmin {
	// 	return model.ErrForbidden
	// }

	return s.repo.UpdateBrand(ctx, params.UpdateBrandParams)
}

func (s *ProductService) DeleteBrand(ctx context.Context, id int64) error {
	return s.repo.DeleteBrand(ctx, id)
}

func (s *ProductService) GetProductModel(ctx context.Context, id int64) (model.ProductModel, error) {
	productModel, err := s.repo.GetProductModel(ctx, id)
	if err != nil {
		return model.ProductModel{}, err
	}

	return productModel, nil
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
		Total:      total,
		NextPage:   model.NextPage(params.Offset, params.Limit, total),
		NextCursor: nil,
	}, nil
}

type CreateProductModelParams struct {
	UserID int64
	model.ProductModel
}

func (s *ProductService) CreateProductModel(ctx context.Context, params CreateProductModelParams) (model.ProductModel, error) {
	if isAdmin, err := s.account.IsAdmin(ctx, params.UserID); err != nil {
		return model.ProductModel{}, err
	} else if !isAdmin {
		return model.ProductModel{}, model.ErrForbidden
	}

	productModel, err := s.repo.CreateProductModel(ctx, params.ProductModel)
	if err != nil {
		return model.ProductModel{}, err
	}

	return productModel, nil
}

type UpdateProductModelParams = repository.UpdateProductModelParams

func (s *ProductService) UpdateProductModel(ctx context.Context, params UpdateProductModelParams) error {
	// if isAdmin, err := s.Account.IsAdmin(ctx, params.UserID); err != nil {
	// 	return err
	// } else if !isAdmin {
	// 	return model.ErrForbidden
	// }

	return s.repo.UpdateProductModel(ctx, params)
}

func (s *ProductService) DeleteProductModel(ctx context.Context, id int64) error {
	return s.repo.DeleteProductModel(ctx, id)
}

type ProductIdentifier = repository.ProductIdentifier

func (s *ProductService) GetProduct(ctx context.Context, params ProductIdentifier) (model.Product, error) {
	return s.repo.GetProduct(ctx, params)
}

type ListProductsParams = repository.ListProductsParams

func (s *ProductService) ListProducts(ctx context.Context, params ListProductsParams) (result model.PaginateResult[model.Product], err error) {
	total, err := s.repo.CountProducts(ctx, params)
	if err != nil {
		return result, err
	}

	products, err := s.repo.ListProducts(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.Product]{
		Data:       products,
		Total:      total,
		NextPage:   model.NextPage(params.Offset, params.Limit, total),
		NextCursor: nil,
	}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	newProduct, err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return model.Product{}, err
	}

	return newProduct, nil
}

type UpdateProductParams = repository.UpdateProductParams

func (s *ProductService) UpdateProduct(ctx context.Context, params UpdateProductParams) error {
	return s.repo.UpdateProduct(ctx, params)
}

func (s *ProductService) DeleteProduct(ctx context.Context, params ProductIdentifier) error {
	return s.repo.DeleteProduct(ctx, params)
}

func (s *ProductService) CreateSale(ctx context.Context, sale model.Sale) (model.Sale, error) {
	return s.repo.CreateSale(ctx, sale)
}

func (s *ProductService) DeleteSale(ctx context.Context, id int64) error {
	return s.repo.DeleteSale(ctx, id)
}

func (s *ProductService) CreateTag(ctx context.Context, tag model.Tag) (model.Tag, error) {
	err := s.repo.CreateTag(ctx, tag)
	if err != nil {
		return model.Tag{}, err
	}

	return tag, nil
}

func (s *ProductService) DeleteTag(ctx context.Context, name string) error {
	return s.repo.DeleteTag(ctx, name)
}
