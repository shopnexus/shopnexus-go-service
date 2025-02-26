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
		Resources:   brand.Resources,
	}, nil
}

type ListBrandsParams struct {
	model.PaginationParams
	Name        *string
	Description *string
}

func (r *Repository) CountBrands(ctx context.Context, params ListBrandsParams) (int64, error) {
	return r.sqlc.CountBrands(ctx, sqlc.CountBrandsParams{
		Name:        *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Description: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Description),
	})
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
			Resources:   brand.Resources,
		}
	}
	return result, nil
}

func (r *Repository) CreateBrand(ctx context.Context, brand model.Brand) (model.Brand, error) {
	row, err := r.sqlc.CreateBrand(ctx, sqlc.CreateBrandParams{
		Name:        brand.Name,
		Description: brand.Description,
		Resources:   brand.Resources,
	})
	if err != nil {
		return model.Brand{}, err
	}

	return model.Brand{
		ID:          row.ID,
		Name:        brand.Name,
		Description: brand.Description,
		Resources:   row.Resources,
	}, nil
}

type UpdateBrandParams struct {
	ID          int64
	Name        *string
	Description *string
}

func (r *Repository) UpdateBrand(ctx context.Context, params UpdateBrandParams) error {
	return r.sqlc.UpdateBrand(ctx, sqlc.UpdateBrandParams{
		ID:          params.ID,
		Name:        *pgxutil.PtrToPgtype(&pgtype.Text{}, &params.Name),
		Description: *pgxutil.PtrToPgtype(&pgtype.Text{}, &params.Description),
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
		Resources:   productModel.Resources,
		Tags:        productModel.Tags,
	}, nil
}

type ListProductModelsParams struct {
	model.PaginationParams
	BrandID              *int64
	Name                 *string
	Description          *string
	ListPriceFrom        *int64
	ListPriceTo          *int64
	DateManufacturedFrom *int64
	DateManufacturedTo   *int64
}

func (r *Repository) CountProductModels(ctx context.Context, params ListProductModelsParams) (int64, error) {
	return r.sqlc.CountProductModels(ctx, sqlc.CountProductModelsParams{
		BrandID:              *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.BrandID),
		Name:                 *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Description:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Description),
		ListPriceFrom:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ListPriceFrom),
		ListPriceTo:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ListPriceTo),
		DateManufacturedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateManufacturedFrom)),
		DateManufacturedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateManufacturedTo)),
	})
}

func (r *Repository) ListProductModels(ctx context.Context, params ListProductModelsParams) ([]model.ProductModel, error) {
	productModels, err := r.sqlc.ListProductModels(ctx, sqlc.ListProductModelsParams{
		Offset:               params.Offset,
		Limit:                params.Limit,
		BrandID:              *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.BrandID),
		Name:                 *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Description:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Description),
		ListPriceFrom:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ListPriceFrom),
		ListPriceTo:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ListPriceTo),
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
			Resources:   productModel.Resources,
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
		Resources:   productModel.Resources,
		Tags:        productModel.Tags,
	})
	if err != nil {
		return model.ProductModel{}, err
	}

	return model.ProductModel{
		ID:               row.ID,
		BrandID:          productModel.BrandID,
		Name:             productModel.Name,
		Description:      productModel.Description,
		ListPrice:        productModel.ListPrice,
		DateManufactured: productModel.DateManufactured,
		Resources:        row.Resources,
		Tags:             row.Tags,
	}, nil
}

type UpdateProductModelParams struct {
	ID               int64
	BrandID          *int64
	Name             *string
	Description      *string
	ListPrice        *int64
	DateManufactured *int64
}

func (r *Repository) UpdateProductModel(ctx context.Context, params UpdateProductModelParams) error {
	return r.sqlc.UpdateProductModel(ctx, sqlc.UpdateProductModelParams{
		ID:               params.ID,
		BrandID:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, &params.BrandID),
		Name:             *pgxutil.PtrToPgtype(&pgtype.Text{}, &params.Name),
		Description:      *pgxutil.PtrToPgtype(&pgtype.Text{}, &params.Description),
		ListPrice:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, &params.ListPrice),
		DateManufactured: *pgxutil.ValueToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateManufactured)),
	})
}

func (r *Repository) DeleteProductModel(ctx context.Context, id int64) error {
	return r.sqlc.DeleteProductModel(ctx, id)
}

// ----------- PRODUCT ------------

// ProductIdentifier is a struct to identify a product, either by ID or SerialID
type ProductIdentifier struct {
	ID       *int64
	SerialID *string
}

type ProductIdentifierPg struct {
	ID       pgtype.Int8
	SerialID pgtype.Text
}

func (p ProductIdentifier) ToPgtype() ProductIdentifierPg {
	return ProductIdentifierPg{
		ID:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, p.ID),
		SerialID: *pgxutil.PtrToPgtype(&pgtype.Text{}, p.SerialID),
	}
}

func (r *Repository) GetProduct(ctx context.Context, id ProductIdentifier) (model.Product, error) {
	row, err := r.sqlc.GetProduct(ctx, sqlc.GetProductParams(id.ToPgtype()))
	if err != nil {
		return model.Product{}, err
	}

	return model.Product{
		ID:             row.ID,
		SerialID:       row.SerialID,
		ProductModelID: row.ProductModelID,
		DateCreated:    row.DateCreated.Time.UnixMilli(),
		DateUpdated:    row.DateUpdated.Time.UnixMilli(),
	}, nil
}

func (r *Repository) GetAvailableProducts(ctx context.Context, productModelID, amount int64) ([]model.Product, error) {
	rows, err := r.sqlc.GetAvailableProducts(ctx, sqlc.GetAvailableProductsParams{
		ProductModelID: productModelID,
		Amount:         int32(amount),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Product, len(rows))
	for i, row := range rows {
		result[i] = model.Product{
			ID:             row.ID,
			SerialID:       row.SerialID,
			ProductModelID: row.ProductModelID,
			DateCreated:    row.DateCreated.Time.UnixMilli(),
			DateUpdated:    row.DateUpdated.Time.UnixMilli(),
		}
	}

	return result, nil
}

type ListProductsParams struct {
	model.PaginationParams
	ProductModelID  *int64
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (r *Repository) CountProducts(ctx context.Context, params ListProductsParams) (int64, error) {
	return r.sqlc.CountProducts(ctx, sqlc.CountProductsParams{
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		DateCreatedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (r *Repository) ListProducts(ctx context.Context, params ListProductsParams) ([]model.Product, error) {
	products, err := r.sqlc.ListProducts(ctx, sqlc.ListProductsParams{
		Offset:          params.Offset,
		Limit:           params.Limit,
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
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

type UpdateProductParams struct {
	ID             int64
	SerialID       *string
	ProductModelID *int64
}

func (r *Repository) UpdateProduct(ctx context.Context, params UpdateProductParams) error {
	return r.sqlc.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:             params.ID,
		SerialID:       *pgxutil.PtrToPgtype(&pgtype.Text{}, params.SerialID),
		ProductModelID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
	})
}

func (r *Repository) DeleteProduct(ctx context.Context, id ProductIdentifier) error {
	return r.sqlc.DeleteProduct(ctx, sqlc.DeleteProductParams(id.ToPgtype()))
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

func (r *Repository) CreateTag(ctx context.Context, tag model.Tag) error {
	return r.sqlc.CreateTag(ctx, sqlc.CreateTagParams{
		TagName:     tag.Name,
		Description: tag.Description,
	})
}

func (r *Repository) DeleteTag(ctx context.Context, name string) error {
	return r.sqlc.DeleteTag(ctx, name)
}
