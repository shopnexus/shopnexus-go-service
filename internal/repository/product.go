package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
)

func (r *Repository) GetBrand(ctx context.Context, id []byte) (model.Brand, error) {
	brand, err := r.sqlc.GetBrand(ctx, id)
	if err != nil {
		return model.Brand{}, err
	}

	return model.Brand{
		ID:          brand.ID,
		Name:        brand.Name,
		Description: brand.Description,
		Images:      brand.Images,
	}, nil
}

type ListBrandsParams struct {
	model.PaginationParams
	Name        *string
	Description *string
}

func (r *Repository) ListBrands(ctx context.Context, params ListBrandsParams) ([]model.Brand, error) {
	brands, err := r.sqlc.ListBrands(ctx, sqlc.ListBrandsParams{
		// Name:        params.Name,
		// Description: params.Description,
		// Limit:       params.Limit,
		// Offset:      params.Page,
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Brand, len(brands))
	for i, brand := range brands {
		result[i] = model.Brand{
			ID:          brand.ID,
			Name:        brand.Name,
			Description: brand.Description,
			// Images:      brand.Images,
		}
	}
	return result, nil
}

func (r *Repository) CreateBrand(ctx context.Context, brand model.Brand) (model.Brand, error) {
	return model.Brand{}, nil
}

func (r *Repository) UpdateBrand(ctx context.Context, brand model.Brand) (model.Brand, error) {
	return model.Brand{}, nil
}

func (r *Repository) DeleteBrand(ctx context.Context, id []byte) error {
	return nil
}

func (r *Repository) GetProductModel(ctx context.Context, id []byte) (model.ProductModel, error) {
	return model.ProductModel{}, nil
}

type ListProductModelsParams struct {
	model.PaginationParams
	BrandID              *[]byte
	Name                 *string
	Description          *string
	ListPrice            *float64
	DateManufacturedFrom *int64
	DateManufacturedTo   *int64
}

func (r *Repository) ListProductModels(ctx context.Context, params ListProductModelsParams) ([]model.ProductModel, error) {
	return nil, nil
}

func (r *Repository) CreateProductModel(ctx context.Context, productModel model.ProductModel) (model.ProductModel, error) {
	return model.ProductModel{}, nil
}

func (r *Repository) UpdateProductModel(ctx context.Context, productModel model.ProductModel) (model.ProductModel, error) {
	return model.ProductModel{}, nil
}

func (r *Repository) DeleteProductModel(ctx context.Context, id []byte) error {
	return nil
}

func (r *Repository) GetProduct(ctx context.Context, serialID []byte) (model.Product, error) {
	return model.Product{}, nil
}

type ListProductsParams struct {
	model.PaginationParams
	ProductModelID  *[]byte
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (r *Repository) ListProducts(ctx context.Context, params ListProductsParams) ([]model.Product, error) {
	return nil, nil
}

func (r *Repository) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	return model.Product{}, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	return model.Product{}, nil
}

func (r *Repository) DeleteProduct(ctx context.Context, serialID []byte) error {
	return nil
}

func (r *Repository) CreateSale(ctx context.Context, sale model.Sale) (model.Sale, error) {
	return model.Sale{}, nil
}

func (r *Repository) DeleteSale(ctx context.Context, id []byte) error {
	return nil
}

func (r *Repository) CreateTag(ctx context.Context, tag model.Tag) (model.Tag, error) {
	return model.Tag{}, nil
}

func (r *Repository) DeleteTag(ctx context.Context, name string) error {
	return nil
}
