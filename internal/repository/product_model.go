package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/util"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) GetProductModel(ctx context.Context, id int64) (model.ProductModel, error) {
	productModel, err := r.sqlc.GetProductModel(ctx, id)
	if err != nil {
		return model.ProductModel{}, err
	}

	return model.ProductModel{
		ID:               productModel.ID,
		BrandID:          productModel.BrandID,
		Name:             productModel.Name,
		Description:      productModel.Description,
		ListPrice:        productModel.ListPrice,
		DateManufactured: productModel.DateManufactured.Time.UnixMilli(),
		Resources:        productModel.Resources,
		Tags:             productModel.Tags,
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
		Offset:               params.Offset(),
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
			ID:               productModel.ID,
			BrandID:          productModel.BrandID,
			Name:             productModel.Name,
			Description:      productModel.Description,
			ListPrice:        productModel.ListPrice,
			DateManufactured: productModel.DateManufactured.Time.UnixMilli(),
			Resources:        productModel.Resources,
			Tags:             productModel.Tags,
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
