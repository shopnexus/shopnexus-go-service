package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/util"

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

func (r *Repository) GetProduct(ctx context.Context, id model.ProductIdentifier) (model.Product[any], error) {
	row, err := r.sqlc.GetProduct(ctx, sqlc.GetProductParams(ToPgtype(id)))
	if err != nil {
		return model.Product[any]{}, err
	}

	return model.Product[any]{
		ID:             row.ID,
		SerialID:       row.SerialID,
		ProductModelID: row.ProductModelID,
		Metadata:       row.Metadata,
		DateCreated:    row.DateCreated.Time.UnixMilli(),
		DateUpdated:    row.DateUpdated.Time.UnixMilli(),
	}, nil
}

func (r *Repository) GetAvailableProducts(ctx context.Context, productModelID, amount int64) ([]model.Product[any], error) {
	rows, err := r.sqlc.GetAvailableProducts(ctx, sqlc.GetAvailableProductsParams{
		ProductModelID: productModelID,
		Amount:         int32(amount),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Product[any], len(rows))
	for i, row := range rows {
		result[i] = model.Product[any]{
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

func (r *Repository) ListProducts(ctx context.Context, params ListProductsParams) ([]model.Product[any], error) {
	products, err := r.sqlc.ListProducts(ctx, sqlc.ListProductsParams{
		Offset:          params.Offset(),
		Limit:           params.Limit,
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		DateCreatedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedTo)),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Product[any], len(products))
	for i, product := range products {
		result[i] = model.Product[any]{
			ID:             product.ID,
			SerialID:       product.SerialID,
			ProductModelID: product.ProductModelID,
			DateCreated:    product.DateCreated.Time.UnixMilli(),
			DateUpdated:    product.DateUpdated.Time.UnixMilli(),
		}
	}

	return result, nil
}

func (r *Repository) CreateProduct(ctx context.Context, product model.Product[any]) (model.Product[any], error) {
	row, err := r.sqlc.CreateProduct(ctx, sqlc.CreateProductParams{
		ProductModelID: product.ProductModelID,
	})
	if err != nil {
		return model.Product[any]{}, err
	}

	return model.Product[any]{
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

func (r *Repository) DeleteProduct(ctx context.Context, id model.ProductIdentifier) error {
	return r.sqlc.DeleteProduct(ctx, sqlc.DeleteProductParams(ToPgtype(id)))
}
