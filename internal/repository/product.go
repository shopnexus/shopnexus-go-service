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

func ToPgtype(p model.ProductIdentifier) ProductIdentifierPg {
	return ProductIdentifierPg{
		ID:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, p.ID),
		SerialID: *pgxutil.PtrToPgtype(&pgtype.Text{}, p.SerialID),
	}
}

func (r *Repository) GetProduct(ctx context.Context, id model.ProductIdentifier) (model.Product, error) {
	row, err := r.sqlc.GetProduct(ctx, sqlc.GetProductParams(ToPgtype(id)))
	if err != nil {
		return model.Product{}, err
	}

	return model.Product{
		ID:             row.ID,
		SerialID:       row.SerialID,
		ProductModelID: row.ProductModelID,
		Quantity:       row.Quantity,
		Sold:           row.Sold,
		AddPrice:       row.AddPrice,
		IsActive:       row.IsActive,
		Metadata:       row.Metadata,
		DateCreated:    row.DateCreated.Time.UnixMilli(),
		DateUpdated:    row.DateUpdated.Time.UnixMilli(),
		Resources:      row.Resources,
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
			Quantity:       row.Quantity,
			Sold:           row.Sold,
			AddPrice:       row.AddPrice,
			IsActive:       row.IsActive,
			Metadata:       row.Metadata,
			DateCreated:    row.DateCreated.Time.UnixMilli(),
			DateUpdated:    row.DateUpdated.Time.UnixMilli(),
		}
	}

	return result, nil
}

type ListProductsParams struct {
	model.PaginationParams
	ProductModelID  *int64
	QuantityFrom    *int64
	QuantityTo      *int64
	SoldFrom        *int64
	SoldTo          *int64
	AddPriceFrom    *int64
	AddPriceTo      *int64
	IsActive        *bool
	Metadata        []byte
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

// TODO: fix list pagination product (missing fields in sqlc)
func (r *Repository) CountProducts(ctx context.Context, params ListProductsParams) (int64, error) {
	return r.sqlc.CountProducts(ctx, sqlc.CountProductsParams{
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		QuantityFrom:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.QuantityFrom),
		QuantityTo:      *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.QuantityTo),
		SoldFrom:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.SoldFrom),
		SoldTo:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.SoldTo),
		AddPriceFrom:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AddPriceFrom),
		AddPriceTo:      *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AddPriceTo),
		IsActive:        *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		Metadata:        params.Metadata,
		DateCreatedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (r *Repository) ListProducts(ctx context.Context, params ListProductsParams) ([]model.Product, error) {
	products, err := r.sqlc.ListProducts(ctx, sqlc.ListProductsParams{
		Offset:          params.Offset(),
		Limit:           params.Limit,
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		QuantityFrom:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.QuantityFrom),
		QuantityTo:      *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.QuantityTo),
		SoldFrom:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.SoldFrom),
		SoldTo:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.SoldTo),
		AddPriceFrom:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AddPriceFrom),
		AddPriceTo:      *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AddPriceTo),
		IsActive:        *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.IsActive),
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
			SerialID:       product.SerialID,
			ProductModelID: product.ProductModelID,
			Quantity:       product.Quantity,
			Sold:           product.Sold,
			AddPrice:       product.AddPrice,
			IsActive:       product.IsActive,
			Metadata:       product.Metadata,
			DateCreated:    product.DateCreated.Time.UnixMilli(),
			DateUpdated:    product.DateUpdated.Time.UnixMilli(),
			Resources:      product.Resources,
		}
	}

	return result, nil
}

func (r *Repository) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	row, err := r.sqlc.CreateProduct(ctx, sqlc.CreateProductParams{
		SerialID:       product.SerialID,
		ProductModelID: product.ProductModelID,
		Quantity:       product.Quantity,
		Sold:           product.Sold,
		AddPrice:       product.AddPrice,
		IsActive:       product.IsActive,
		Metadata:       product.Metadata,
		Resources:      product.Resources,
	})
	if err != nil {
		return model.Product{}, err
	}

	return model.Product{
		ID:             row.ID,
		SerialID:       product.SerialID,
		ProductModelID: product.ProductModelID,
		Quantity:       product.Quantity,
		Sold:           product.Sold,
		AddPrice:       product.AddPrice,
		IsActive:       product.IsActive,
		Metadata:       product.Metadata,
		DateCreated:    time.Now().UnixMilli(),
		DateUpdated:    time.Now().UnixMilli(),
		Resources:      row.Resources,
	}, nil
}

type UpdateProductParams struct {
	ID             int64
	SerialID       *string
	ProductModelID *int64
	Quantity       *int64
	Sold           *int64
	AddPrice       *int64
	IsActive       *bool
	Metadata       []byte
}

func (r *Repository) UpdateProduct(ctx context.Context, params UpdateProductParams) error {
	return r.sqlc.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:             params.ID,
		SerialID:       *pgxutil.PtrToPgtype(&pgtype.Text{}, params.SerialID),
		ProductModelID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		Quantity:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Quantity),
		Sold:           *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Sold),
		AddPrice:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AddPrice),
		IsActive:       *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		Metadata:       params.Metadata,
	})
}

func (r *Repository) DeleteProduct(ctx context.Context, id model.ProductIdentifier) error {
	return r.sqlc.DeleteProduct(ctx, sqlc.DeleteProductParams(ToPgtype(id)))
}

func (r *Repository) GetResources(ctx context.Context, ownerID int64) ([]string, error) {
	return r.sqlc.GetResources(ctx, ownerID)
}

func (r *Repository) AddResources(ctx context.Context, ownerID int64, resouces []string) error {
	return r.sqlc.AddResources(ctx, sqlc.AddResourcesParams{
		OwnerID:   ownerID,
		Resources: resouces,
	})
}

func (r *Repository) RemoveResources(ctx context.Context, ownerID int64, resources []string) error {
	return r.sqlc.RemoveResources(ctx, sqlc.RemoveResourcesParams{
		OwnerID:   ownerID,
		Resources: resources,
	})
}
