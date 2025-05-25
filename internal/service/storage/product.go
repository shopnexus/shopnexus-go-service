package storage

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/utils/ptr"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type ProductIdentifierPg struct {
	ID       pgtype.Int8
	SerialID pgtype.Text
}

func (r *ServiceImpl) GetProduct(ctx context.Context, id int64) (model.Product, error) {
	row, err := r.sqlc.GetProduct(ctx, id)
	if err != nil {
		return model.Product{}, err
	}

	return model.Product{
		ID:             row.ID,
		ProductModelID: row.ProductModelID,
		Quantity:       row.Quantity,
		Sold:           row.Sold,
		AddPrice:       row.AddPrice,
		IsActive:       row.IsActive,
		CanCombine:     row.CanCombine,
		Metadata:       row.Metadata,
		DateCreated:    row.DateCreated.Time.UnixMilli(),
		DateUpdated:    row.DateUpdated.Time.UnixMilli(),
		Resources:      row.Resources,
	}, nil
}

type ListProductsParams struct {
	model.PaginationParams
	ID              *int64
	ProductModelID  *int64
	QuantityFrom    *int64
	QuantityTo      *int64
	SoldFrom        *int64
	SoldTo          *int64
	AddPriceFrom    *int64
	AddPriceTo      *int64
	IsActive        *bool
	CanCombine      *bool
	Metadata        []byte
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (r *ServiceImpl) CountProducts(ctx context.Context, params ListProductsParams) (int64, error) {
	return r.sqlc.CountProducts(ctx, sqlc.CountProductsParams{
		ID:              *PtrToPgtype(&pgtype.Int8{}, params.ID),
		ProductModelID:  *PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		QuantityFrom:    *PtrToPgtype(&pgtype.Int8{}, params.QuantityFrom),
		QuantityTo:      *PtrToPgtype(&pgtype.Int8{}, params.QuantityTo),
		SoldFrom:        *PtrToPgtype(&pgtype.Int8{}, params.SoldFrom),
		SoldTo:          *PtrToPgtype(&pgtype.Int8{}, params.SoldTo),
		AddPriceFrom:    *PtrToPgtype(&pgtype.Int8{}, params.AddPriceFrom),
		AddPriceTo:      *PtrToPgtype(&pgtype.Int8{}, params.AddPriceTo),
		IsActive:        *PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		CanCombine:      *PtrToPgtype(&pgtype.Bool{}, params.CanCombine),
		Metadata:        params.Metadata,
		DateCreatedFrom: *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (r *ServiceImpl) ListProducts(ctx context.Context, params ListProductsParams) ([]model.Product, error) {
	productRows, err := r.sqlc.ListProducts(ctx, sqlc.ListProductsParams{
		Offset:          params.Offset(),
		Limit:           params.Limit,
		ID:              *PtrToPgtype(&pgtype.Int8{}, params.ID),
		ProductModelID:  *PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		QuantityFrom:    *PtrToPgtype(&pgtype.Int8{}, params.QuantityFrom),
		QuantityTo:      *PtrToPgtype(&pgtype.Int8{}, params.QuantityTo),
		SoldFrom:        *PtrToPgtype(&pgtype.Int8{}, params.SoldFrom),
		SoldTo:          *PtrToPgtype(&pgtype.Int8{}, params.SoldTo),
		AddPriceFrom:    *PtrToPgtype(&pgtype.Int8{}, params.AddPriceFrom),
		AddPriceTo:      *PtrToPgtype(&pgtype.Int8{}, params.AddPriceTo),
		IsActive:        *PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		CanCombine:      *PtrToPgtype(&pgtype.Bool{}, params.CanCombine),
		Metadata:        params.Metadata,
		DateCreatedFrom: *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedTo)),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Product, len(productRows))
	for i, productRow := range productRows {
		result[i] = model.Product{
			ID:             productRow.ID,
			ProductModelID: productRow.ProductModelID,
			Quantity:       productRow.Quantity,
			Sold:           productRow.Sold,
			AddPrice:       productRow.AddPrice,
			IsActive:       productRow.IsActive,
			CanCombine:     productRow.CanCombine,
			Metadata:       productRow.Metadata,
			DateCreated:    productRow.DateCreated.Time.UnixMilli(),
			DateUpdated:    productRow.DateUpdated.Time.UnixMilli(),
			Resources:      productRow.Resources,
		}
	}

	return result, nil
}

func (r *ServiceImpl) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	row, err := r.sqlc.CreateProduct(ctx, sqlc.CreateProductParams{
		ProductModelID: product.ProductModelID,
		Quantity:       product.Quantity,
		Sold:           product.Sold,
		AddPrice:       product.AddPrice,
		CanCombine:     product.CanCombine,
		IsActive:       product.IsActive,
		Metadata:       product.Metadata,
	})
	if err != nil {
		return model.Product{}, err
	}

	if err := r.AddResources(ctx, row.ID, model.ResourceTypeProduct, product.Resources); err != nil {
		return model.Product{}, err
	}

	return model.Product{
		ID:             row.ID,
		ProductModelID: product.ProductModelID,
		Quantity:       product.Quantity,
		Sold:           product.Sold,
		AddPrice:       product.AddPrice,
		IsActive:       product.IsActive,
		CanCombine:     product.CanCombine,
		Metadata:       product.Metadata,
		DateCreated:    time.Now().UnixMilli(),
		DateUpdated:    time.Now().UnixMilli(),
		Resources:      product.Resources,
	}, nil
}

type UpdateProductParams struct {
	ID             int64
	ProductModelID *int64
	Quantity       *int64
	Sold           *int64
	AddPrice       *int64
	CanCombine     *bool
	IsActive       *bool
	Metadata       *[]byte
	Resources      *[]string
}

func (r *ServiceImpl) UpdateProduct(ctx context.Context, params UpdateProductParams) error {
	storageParams := sqlc.UpdateProductParams{
		ID:             params.ID,
		ProductModelID: *PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		Quantity:       *PtrToPgtype(&pgtype.Int8{}, params.Quantity),
		Sold:           *PtrToPgtype(&pgtype.Int8{}, params.Sold),
		AddPrice:       *PtrToPgtype(&pgtype.Int8{}, params.AddPrice),
		CanCombine:     *PtrToPgtype(&pgtype.Bool{}, params.CanCombine),
		IsActive:       *PtrToPgtype(&pgtype.Bool{}, params.IsActive),
	}

	if params.Metadata != nil {
		storageParams.Metadata = *params.Metadata
	}

	if params.Resources != nil {
		if err := r.UpdateResources(ctx, params.ID, model.ResourceTypeProduct, *params.Resources); err != nil {
			return err
		}
	}

	return r.sqlc.UpdateProduct(ctx, storageParams)
}

func (r *ServiceImpl) UpdateProductSold(ctx context.Context, ids []int64, amount int64) error {
	return r.sqlc.UpdateProductSold(ctx, sqlc.UpdateProductSoldParams{
		Ids:    ids,
		Amount: amount,
	})
}

func (r *ServiceImpl) DeleteProduct(ctx context.Context, id int64) error {
	return r.sqlc.DeleteProduct(ctx, id)
}
