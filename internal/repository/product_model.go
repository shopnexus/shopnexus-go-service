package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/util"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *RepositoryImpl) GetProductModel(ctx context.Context, id int64) (model.ProductModel, error) {
	productModel, err := r.sqlc.GetProductModel(ctx, id)
	if err != nil {
		return model.ProductModel{}, err
	}

	resources, err := r.GetResources(ctx, productModel.ID, model.ResourceTypeProductModel)
	if err != nil {
		return model.ProductModel{}, err
	}

	tags, err := r.GetTags(ctx, productModel.ID)
	if err != nil {
		return model.ProductModel{}, err
	}

	return model.ProductModel{
		ID:               productModel.ID,
		Type:             productModel.Type,
		BrandID:          productModel.BrandID,
		Name:             productModel.Name,
		Description:      productModel.Description,
		ListPrice:        productModel.ListPrice,
		DateManufactured: productModel.DateManufactured.Time.UnixMilli(),
		Resources:        resources,
		Tags:             tags,
	}, nil
}

func (r *RepositoryImpl) GetProductSerialIDs(ctx context.Context, productID int64) ([]string, error) {
	return r.sqlc.GetProductSerialIDs(ctx, productID)
}

type ListProductModelsParams struct {
	model.PaginationParams
	Type                 *int64
	BrandID              *int64
	Name                 *string
	Description          *string
	ListPriceFrom        *int64
	ListPriceTo          *int64
	DateManufacturedFrom *int64
	DateManufacturedTo   *int64
}

func (r *RepositoryImpl) CountProductModels(ctx context.Context, params ListProductModelsParams) (int64, error) {
	return r.sqlc.CountProductModels(ctx, sqlc.CountProductModelsParams{
		Type:                 *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Type),
		BrandID:              *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.BrandID),
		Name:                 *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Description:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Description),
		ListPriceFrom:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ListPriceFrom),
		ListPriceTo:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ListPriceTo),
		DateManufacturedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateManufacturedFrom)),
		DateManufacturedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateManufacturedTo)),
	})
}

func (r *RepositoryImpl) ListProductModels(ctx context.Context, params ListProductModelsParams) ([]model.ProductModel, error) {
	productModels, err := r.sqlc.ListProductModels(ctx, sqlc.ListProductModelsParams{
		Offset:               params.Offset(),
		Limit:                params.Limit,
		Type:                 *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Type),
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
		resources, err := r.GetResources(ctx, productModel.ID, model.ResourceTypeProductModel)
		if err != nil {
			return nil, err
		}

		tags, err := r.GetTags(ctx, productModel.ID)
		if err != nil {
			return nil, err
		}

		result[i] = model.ProductModel{
			ID:               productModel.ID,
			Type:             productModel.Type,
			BrandID:          productModel.BrandID,
			Name:             productModel.Name,
			Description:      productModel.Description,
			ListPrice:        productModel.ListPrice,
			DateManufactured: productModel.DateManufactured.Time.UnixMilli(),
			Resources:        resources,
			Tags:             tags,
		}
	}

	return result, nil
}

func (r *RepositoryImpl) CreateProductModel(ctx context.Context, productModel model.ProductModel) (model.ProductModel, error) {
	row, err := r.sqlc.CreateProductModel(ctx, sqlc.CreateProductModelParams{
		Type:             productModel.Type,
		BrandID:          productModel.BrandID,
		Name:             productModel.Name,
		Description:      productModel.Description,
		ListPrice:        productModel.ListPrice,
		DateManufactured: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(&productModel.DateManufactured)),
	})
	if err != nil {
		return model.ProductModel{}, err
	}

	if err = r.AddResources(ctx, row.ID, model.ResourceTypeProductModel, productModel.Resources); err != nil {
		return model.ProductModel{}, err
	}

	if err = r.AddTags(ctx, row.ID, productModel.Tags); err != nil {
		return model.ProductModel{}, err
	}

	return model.ProductModel{
		ID:               row.ID,
		Type:             productModel.Type,
		BrandID:          productModel.BrandID,
		Name:             productModel.Name,
		Description:      productModel.Description,
		ListPrice:        productModel.ListPrice,
		DateManufactured: productModel.DateManufactured,
		Resources:        productModel.Resources,
		Tags:             productModel.Tags,
	}, nil
}

type UpdateProductModelParams struct {
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

func (r *RepositoryImpl) UpdateProductModel(ctx context.Context, params UpdateProductModelParams) error {
	row, err := r.sqlc.UpdateProductModel(ctx, sqlc.UpdateProductModelParams{
		ID:               params.ID,
		Type:             *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Type),
		BrandID:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.BrandID),
		Name:             *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Description:      *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Description),
		ListPrice:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ListPrice),
		DateManufactured: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateManufactured)),
	})
	if err != nil {
		return err
	}

	if params.Resources != nil {
		if err = r.UpdateResources(ctx, row.ID, model.ResourceTypeProductModel, *params.Resources); err != nil {
			return err
		}
	}

	if params.Tags != nil {
		if err = r.UpdateTags(ctx, row.ID, *params.Tags); err != nil {
			return err
		}
	}

	return nil
}

func (r *RepositoryImpl) DeleteProductModel(ctx context.Context, id int64) error {
	return r.sqlc.DeleteProductModel(ctx, id)
}

type ListProductTypesParams struct {
	model.PaginationParams
	Name *string
}

func (r *RepositoryImpl) CountProductTypes(ctx context.Context, params ListProductTypesParams) (int64, error) {
	return r.sqlc.CountProductTypes(ctx, *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name))
}

func (r *RepositoryImpl) ListProductTypes(ctx context.Context, params ListProductTypesParams) ([]model.ProductType, error) {
	productTypes, err := r.sqlc.ListProductTypes(ctx, sqlc.ListProductTypesParams{
		Offset: params.Offset(),
		Limit:  params.Limit,
		Name:   *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.ProductType, len(productTypes))
	for i, productType := range productTypes {
		result[i] = model.ProductType{
			ID:   productType.ID,
			Name: productType.Name,
		}
	}

	return result, nil
}
