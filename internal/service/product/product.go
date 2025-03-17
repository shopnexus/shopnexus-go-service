package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
)

type ProductService struct {
	repo *repository.Repository
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
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
		NextCursor: nil,
	}, nil
}

type CreateBrandParams struct {
	UserID int64
	model.Brand
}

func (s *ProductService) CreateBrand(ctx context.Context, params CreateBrandParams) (model.Brand, error) {
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

func (s *ProductService) GetProduct(ctx context.Context, params model.ProductIdentifier) (model.Product, error) {
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
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
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

func (s *ProductService) DeleteProduct(ctx context.Context, params model.ProductIdentifier) error {
	return s.repo.DeleteProduct(ctx, params)
}

type ProductServiceInterface interface {
	GetBrand(ctx context.Context, id int64) (model.Brand, error)
	ListBrands(ctx context.Context, params ListBrandsParams) (model.PaginateResult[model.Brand], error)
	CreateBrand(ctx context.Context, params CreateBrandParams) (model.Brand, error)
	UpdateBrand(ctx context.Context, params UpdateBrandParams) error
	DeleteBrand(ctx context.Context, id int64) error

	GetProductModel(ctx context.Context, id int64) (model.ProductModel, error)
	ListProductModels(ctx context.Context, params ListProductModelsParams) (model.PaginateResult[model.ProductModel], error)
	CreateProductModel(ctx context.Context, params CreateProductModelParams) (model.ProductModel, error)
	UpdateProductModel(ctx context.Context, params UpdateProductModelParams) error
	DeleteProductModel(ctx context.Context, id int64) error

	GetProduct(ctx context.Context, params model.ProductIdentifier) (model.Product, error)
	ListProducts(ctx context.Context, params ListProductsParams) (model.PaginateResult[model.Product], error)
	CreateProduct(ctx context.Context, product model.Product) (model.Product, error)
	UpdateProduct(ctx context.Context, params UpdateProductParams) error
	DeleteProduct(ctx context.Context, params model.ProductIdentifier) error

	CreateSale(ctx context.Context, sale model.Sale) (model.Sale, error)
	DeleteSale(ctx context.Context, id int64) error

	CreateTag(ctx context.Context, tag model.Tag) (model.Tag, error)
	DeleteTag(ctx context.Context, name string) error
}
