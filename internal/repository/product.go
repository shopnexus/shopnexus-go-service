package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/util"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) GetBrand(ctx context.Context, id int64) (model.Brand, error) {
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
		Offset:      params.Offset,
		Limit:       params.Limit,
		Name:        *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Description: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Description),
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
			Images:      brand.Images,
		}
	}
	return result, nil
}

func (r *Repository) CreateBrand(ctx context.Context, brand model.Brand) (model.Brand, error) {
	row, err := r.sqlc.CreateBrand(ctx, sqlc.CreateBrandParams{
		Name:        brand.Name,
		Description: brand.Description,
		Images:      brand.Images,
	})
	if err != nil {
		return model.Brand{}, err
	}

	return model.Brand{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Images:      row.Images,
	}, nil
}

func (r *Repository) UpdateBrand(ctx context.Context, brand model.Brand) error {
	return r.sqlc.UpdateBrand(ctx, sqlc.UpdateBrandParams{
		ID:          brand.ID,
		Name:        *pgxutil.PtrToPgtype(&pgtype.Text{}, &brand.Name),
		Description: *pgxutil.PtrToPgtype(&pgtype.Text{}, &brand.Description),
	})
}

func (r *Repository) DeleteBrand(ctx context.Context, id int64) error {
	return r.sqlc.DeleteBrand(ctx, id)
}

func (r *Repository) GetProductModel(ctx context.Context, id int64) (model.ProductModel, error) {
	productModel, err := r.sqlc.GetProductModel(ctx, id)
	if err != nil {
		return model.ProductModel{}, err
	}

	return model.ProductModel{
		ID:          productModel.ID,
		BrandID:     productModel.BrandID,
		Name:        productModel.Name,
		Description: productModel.Description,
		ListPrice:   productModel.ListPrice,
		Images:      productModel.Images,
		Tags:        productModel.Tags,
	}, nil
}

type ListProductModelsParams struct {
	model.PaginationParams
	BrandID              int64
	Name                 *string
	Description          *string
	ListPrice            *int64
	DateManufacturedFrom *int64
	DateManufacturedTo   *int64
}

func (r *Repository) ListProductModels(ctx context.Context, params ListProductModelsParams) ([]model.ProductModel, error) {
	productModels, err := r.sqlc.ListProductModels(ctx, sqlc.ListProductModelsParams{
		Offset:               params.Offset,
		Limit:                params.Limit,
		BrandID:              *pgxutil.PtrToPgtype(&pgtype.Int8{}, &params.BrandID),
		Name:                 *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Description:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Description),
		ListPrice:            *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ListPrice),
		DateManufacturedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateManufacturedFrom)),
		DateManufacturedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateManufacturedTo)),
	})

	if err != nil {
		return nil, err
	}

	result := make([]model.ProductModel, len(productModels))
	for i, productModel := range productModels {
		result[i] = model.ProductModel{
			ID:          productModel.ID,
			BrandID:     productModel.BrandID,
			Name:        productModel.Name,
			Description: productModel.Description,
			ListPrice:   productModel.ListPrice,
			Images:      productModel.Images,
			Tags:        productModel.Tags,
		}
	}

	return result, nil
}

func (r *Repository) CreateProductModel(ctx context.Context, productModel model.ProductModel) (model.ProductModel, error) {
	row, err := r.sqlc.CreateProductModel(ctx, sqlc.CreateProductModelParams{
		BrandID:     productModel.BrandID,
		Name:        productModel.Name,
		Description: productModel.Description,
		ListPrice:   productModel.ListPrice,
		Images:      productModel.Images,
		Tags:        productModel.Tags,
	})
	if err != nil {
		return model.ProductModel{}, err
	}

	return model.ProductModel{
		ID:          row.ID,
		BrandID:     row.BrandID,
		Name:        row.Name,
		Description: row.Description,
		ListPrice:   row.ListPrice,
		Images:      row.Images,
		Tags:        row.Tags,
	}, nil
}

func (r *Repository) UpdateProductModel(ctx context.Context, productModel model.ProductModel) error {
	return r.sqlc.UpdateProductModel(ctx, sqlc.UpdateProductModelParams{
		ID:               productModel.ID,
		BrandID:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, &productModel.BrandID),
		Name:             *pgxutil.PtrToPgtype(&pgtype.Text{}, &productModel.Name),
		Description:      *pgxutil.PtrToPgtype(&pgtype.Text{}, &productModel.Description),
		ListPrice:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, &productModel.ListPrice),
		DateManufactured: *pgxutil.ValueToPgtype(&pgtype.Timestamptz{}, time.UnixMilli(productModel.DateManufactured)),
	})
}

func (r *Repository) DeleteProductModel(ctx context.Context, id int64) error {
	return r.sqlc.DeleteProductModel(ctx, id)
}

func (r *Repository) GetProduct(ctx context.Context, serialID string) (model.Product, error) {
	row, err := r.sqlc.GetProduct(ctx, serialID)
	if err != nil {
		return model.Product{}, err
	}

	return model.Product{
		SerialID:       row.SerialID,
		ProductModelID: row.ProductModelID,
		DateCreated:    row.DateCreated.Time.UnixMilli(),
		DateUpdated:    row.DateUpdated.Time.UnixMilli(),
	}, nil
}

type ListProductsParams struct {
	model.PaginationParams
	ProductModelID  int64
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (r *Repository) ListProducts(ctx context.Context, params ListProductsParams) ([]model.Product, error) {
	products, err := r.sqlc.ListProducts(ctx, sqlc.ListProductsParams{
		Offset:          params.Offset,
		Limit:           params.Limit,
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, &params.ProductModelID),
		DateCreatedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedTo)),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Product, len(products))
	for i, product := range products {
		result[i] = model.Product{
			SerialID:       product.SerialID,
			ProductModelID: product.ProductModelID,
			DateCreated:    product.DateCreated.Time.UnixMilli(),
			DateUpdated:    product.DateUpdated.Time.UnixMilli(),
		}
	}

	return result, nil
}

func (r *Repository) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	row, err := r.sqlc.CreateProduct(ctx, sqlc.CreateProductParams{
		ProductModelID: product.ProductModelID,
	})
	if err != nil {
		return model.Product{}, err
	}

	return model.Product{
		SerialID:       row.SerialID,
		ProductModelID: row.ProductModelID,
		DateCreated:    row.DateCreated.Time.UnixMilli(),
		DateUpdated:    row.DateUpdated.Time.UnixMilli(),
	}, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, product model.Product) error {
	return r.sqlc.UpdateProduct(ctx, sqlc.UpdateProductParams{
		SerialID:       product.SerialID,
		ProductModelID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, &product.ProductModelID),
	})
}

func (r *Repository) DeleteProduct(ctx context.Context, serialID string) error {
	return r.sqlc.DeleteProduct(ctx, serialID)
}

func (r *Repository) CreateSale(ctx context.Context, sale model.Sale) (model.Sale, error) {
	row, err := r.sqlc.CreateSale(ctx, sqlc.CreateSaleParams{
		TagName:         *pgxutil.PtrToPgtype(&pgtype.Text{}, sale.Tag),
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, &sale.ProductModelID),
		DateStarted:     *pgxutil.ValueToPgtype(&pgtype.Timestamptz{}, time.UnixMilli(sale.DateStarted)),
		DateEnded:       *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(sale.DateEnded)),
		Quantity:        sale.Quantity,
		Used:            sale.Used,
		IsActive:        sale.IsActive,
		DiscountPercent: *pgxutil.PtrToPgtype(&pgtype.Int8{}, sale.DiscountPercent),
		DiscountPrice:   *pgxutil.PtrToPgtype(&pgtype.Int8{}, sale.DiscountPrice),
	})
	if err != nil {
		return model.Sale{}, err
	}

	return model.Sale{
		ID:              row.ID,
		Tag:             pgxutil.PgtypeToPtr[string](row.TagName),
		ProductModelID:  pgxutil.PgtypeToPtr[int64](row.ProductModelID),
		DateStarted:     row.DateStarted.Time.UnixMilli(),
		DateEnded:       util.PtrTimeToMilis(pgxutil.PgtypeToPtr[time.Time](row.DateEnded)),
		Quantity:        row.Quantity,
		Used:            row.Used,
		IsActive:        row.IsActive,
		DiscountPercent: pgxutil.PgtypeToPtr[int64](row.DiscountPercent),
		DiscountPrice:   pgxutil.PgtypeToPtr[int64](row.DiscountPrice),
	}, nil
}

func (r *Repository) DeleteSale(ctx context.Context, id int64) error {
	return r.sqlc.DeleteSale(ctx, id)
}

func (r *Repository) CreateTag(ctx context.Context, tag model.Tag) (model.Tag, error) {
	row, err := r.sqlc.CreateTag(ctx, sqlc.CreateTagParams{
		TagName:     tag.Name,
		Description: tag.Description,
	})
	if err != nil {
		return model.Tag{}, err
	}

	return model.Tag{
		Name:        row.TagName,
		Description: row.Description,
	}, nil
}

func (r *Repository) DeleteTag(ctx context.Context, name string) error {
	return r.sqlc.DeleteTag(ctx, name)
}
