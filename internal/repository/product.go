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

type ProductIdentifierPg struct {
	ID       pgtype.Int8
	SerialID pgtype.Text
}

func (r *RepositoryImpl) GetProduct(ctx context.Context, id int64) (model.Product, error) {
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

func (r *RepositoryImpl) CountProducts(ctx context.Context, params ListProductsParams) (int64, error) {
	return r.sqlc.CountProducts(ctx, sqlc.CountProductsParams{
		ID:              *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ID),
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		QuantityFrom:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.QuantityFrom),
		QuantityTo:      *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.QuantityTo),
		SoldFrom:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.SoldFrom),
		SoldTo:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.SoldTo),
		AddPriceFrom:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AddPriceFrom),
		AddPriceTo:      *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AddPriceTo),
		IsActive:        *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		CanCombine:      *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.CanCombine),
		Metadata:        params.Metadata,
		DateCreatedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (r *RepositoryImpl) ListProducts(ctx context.Context, params ListProductsParams) ([]model.Product, error) {
	products, err := r.sqlc.ListProducts(ctx, sqlc.ListProductsParams{
		Offset:          params.Offset(),
		Limit:           params.Limit,
		ID:              *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ID),
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		QuantityFrom:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.QuantityFrom),
		QuantityTo:      *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.QuantityTo),
		SoldFrom:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.SoldFrom),
		SoldTo:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.SoldTo),
		AddPriceFrom:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AddPriceFrom),
		AddPriceTo:      *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AddPriceTo),
		IsActive:        *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		CanCombine:      *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.CanCombine),
		Metadata:        params.Metadata,
		DateCreatedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedTo)),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Product, len(products))
	for i, product := range products {
		result[i] = model.Product{
			ID:             product.ID,
			ProductModelID: product.ProductModelID,
			Quantity:       product.Quantity,
			Sold:           product.Sold,
			AddPrice:       product.AddPrice,
			IsActive:       product.IsActive,
			CanCombine:     product.CanCombine,
			Metadata:       product.Metadata,
			DateCreated:    product.DateCreated.Time.UnixMilli(),
			DateUpdated:    product.DateUpdated.Time.UnixMilli(),
			Resources:      product.Resources,
		}
	}

	return result, nil
}

func (r *RepositoryImpl) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	row, err := r.sqlc.CreateProduct(ctx, sqlc.CreateProductParams{
		ProductModelID: product.ProductModelID,
		Quantity:       product.Quantity,
		Sold:           product.Sold,
		AddPrice:       product.AddPrice,
		CanCombine:     product.CanCombine,
		IsActive:       product.IsActive,
		Metadata:       product.Metadata,
		Resources:      product.Resources,
	})
	if err != nil {
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
		Resources:      row.Resources,
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
	Metadata       []byte
}

func (r *RepositoryImpl) UpdateProduct(ctx context.Context, params UpdateProductParams) error {
	return r.sqlc.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:             params.ID,
		ProductModelID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		Quantity:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Quantity),
		Sold:           *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Sold),
		AddPrice:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AddPrice),
		CanCombine:     *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.CanCombine),
		IsActive:       *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		Metadata:       params.Metadata,
	})
}

func (r *RepositoryImpl) UpdateProductSold(ctx context.Context, ids []int64, amount int64) error {
	return r.sqlc.UpdateProductSold(ctx, sqlc.UpdateProductSoldParams{
		Ids:    ids,
		Amount: amount,
	})
}

func (r *RepositoryImpl) DeleteProduct(ctx context.Context, id int64) error {
	return r.sqlc.DeleteProduct(ctx, id)
}

func (r *RepositoryImpl) GetResources(ctx context.Context, ownerID int64) ([]string, error) {
	return r.sqlc.GetResources(ctx, ownerID)
}

func (r *RepositoryImpl) AddResources(ctx context.Context, ownerID int64, resouces []string) error {
	return r.sqlc.AddResources(ctx, sqlc.AddResourcesParams{
		OwnerID:   ownerID,
		Resources: resouces,
	})
}

func (r *RepositoryImpl) RemoveResources(ctx context.Context, ownerID int64, resources []string) error {
	return r.sqlc.RemoveResources(ctx, sqlc.RemoveResourcesParams{
		OwnerID:   ownerID,
		Resources: resources,
	})
}
