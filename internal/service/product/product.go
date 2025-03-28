package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
)

type ProductService struct {
	repo *repository.Repository
}

var _ ProductServiceInterface = (*ProductService)(nil)

func NewProductService(repo *repository.Repository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, params model.ProductIdentifier) (model.Product[any], error) {
	return s.repo.GetProduct(ctx, params)
}

type ListProductsParams = repository.ListProductsParams

func (s *ProductService) ListProducts(ctx context.Context, params ListProductsParams) (result model.PaginateResult[model.Product[any]], err error) {
	total, err := s.repo.CountProducts(ctx, params)
	if err != nil {
		return result, err
	}

	products, err := s.repo.ListProducts(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.Product[any]]{
		Data:       products,
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
		NextCursor: nil,
	}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product model.Product[any]) (model.Product[any], error) {
	newProduct, err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return model.Product[any]{}, err
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
	ListProductTypes(ctx context.Context, params ListProductTypesParams) ([]model.ProductType, error)

	GetProduct(ctx context.Context, params model.ProductIdentifier) (model.Product[any], error)
	ListProducts(ctx context.Context, params ListProductsParams) (model.PaginateResult[model.Product[any]], error)
	CreateProduct(ctx context.Context, product model.Product[any]) (model.Product[any], error)
	UpdateProduct(ctx context.Context, params UpdateProductParams) error
	DeleteProduct(ctx context.Context, params model.ProductIdentifier) error

	GetSale(ctx context.Context, id int64) (model.Sale, error)
	ListSales(ctx context.Context, params ListSalesParams) (model.PaginateResult[model.Sale], error)
	CreateSale(ctx context.Context, params CreateSaleParams) (model.Sale, error)
	UpdateSale(ctx context.Context, params UpdateSaleParams) error
	DeleteSale(ctx context.Context, id int64) error

	GetTag(ctx context.Context, tag string) (TagResponse, error)
	ListTags(ctx context.Context, params ListTagsParams) (model.PaginateResult[TagResponse], error)
	CreateTag(ctx context.Context, tag model.Tag) error
	UpdateTag(ctx context.Context, params UpdateTagParams) error
	DeleteTag(ctx context.Context, tag string) error
}
