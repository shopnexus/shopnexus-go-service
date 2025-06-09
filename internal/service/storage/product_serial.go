package storage

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/utils/ptr"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *ServiceImpl) GetProductSerial(ctx context.Context, serialID string) (model.ProductSerial, error) {
	row, err := r.sqlc.GetProductSerial(ctx, serialID)
	if err != nil {
		return model.ProductSerial{}, err
	}

	return model.ProductSerial{
		SerialID:    row.SerialID,
		ProductID:   row.ProductID,
		IsSold:      row.IsSold,
		IsActive:    row.IsActive,
		DateCreated: row.DateCreated.Time.UnixMilli(),
	}, nil
}

type ListProductSerialsParams struct {
	model.PaginationParams
	SerialID        *string
	ProductID       *int64
	IsSold          *bool
	IsActive        *bool
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (r *ServiceImpl) CountProductSerials(ctx context.Context, params ListProductSerialsParams) (int64, error) {
	return r.sqlc.CountProductSerials(ctx, sqlc.CountProductSerialsParams{
		SerialID:        *PtrToPgtype(&pgtype.Text{}, params.SerialID),
		ProductID:       *PtrToPgtype(&pgtype.Int8{}, params.ProductID),
		IsSold:          *PtrToPgtype(&pgtype.Bool{}, params.IsSold),
		IsActive:        *PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		DateCreatedFrom: *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (r *ServiceImpl) ListProductSerials(ctx context.Context, params ListProductSerialsParams) ([]model.ProductSerial, error) {
	serials, err := r.sqlc.ListProductSerials(ctx, sqlc.ListProductSerialsParams{
		Offset:          params.Offset(),
		Limit:           params.Limit,
		SerialID:        *PtrToPgtype(&pgtype.Text{}, params.SerialID),
		ProductID:       *PtrToPgtype(&pgtype.Int8{}, params.ProductID),
		IsSold:          *PtrToPgtype(&pgtype.Bool{}, params.IsSold),
		IsActive:        *PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		DateCreatedFrom: *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedTo)),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.ProductSerial, len(serials))
	for i, serial := range serials {
		result[i] = model.ProductSerial{
			SerialID:    serial.SerialID,
			ProductID:   serial.ProductID,
			IsSold:      serial.IsSold,
			IsActive:    serial.IsActive,
			DateCreated: serial.DateCreated.Time.UnixMilli(),
		}
	}

	return result, nil
}

func (r *ServiceImpl) CreateProductSerial(ctx context.Context, serial model.ProductSerial) (model.ProductSerial, error) {
	row, err := r.sqlc.CreateProductSerial(ctx, sqlc.CreateProductSerialParams{
		SerialID:  serial.SerialID,
		ProductID: serial.ProductID,
		IsSold:    serial.IsSold,
		IsActive:  serial.IsActive,
	})
	if err != nil {
		return model.ProductSerial{}, err
	}

	return model.ProductSerial{
		SerialID:    row.SerialID,
		ProductID:   row.ProductID,
		IsSold:      row.IsSold,
		IsActive:    row.IsActive,
		DateCreated: time.Now().UnixMilli(),
		DateUpdated: time.Now().UnixMilli(),
	}, nil
}

type UpdateProductSerialParams struct {
	SerialID string
	IsSold   *bool
	IsActive *bool
}

func (r *ServiceImpl) UpdateProductSerial(ctx context.Context, params UpdateProductSerialParams) error {
	return r.sqlc.UpdateProductSerial(ctx, sqlc.UpdateProductSerialParams{
		SerialID: params.SerialID,
		IsSold:   *PtrToPgtype(&pgtype.Bool{}, params.IsSold),
		IsActive: *PtrToPgtype(&pgtype.Bool{}, params.IsActive),
	})
}

func (r *ServiceImpl) DeleteProductSerial(ctx context.Context, serialID string) error {
	return r.sqlc.DeleteProductSerial(ctx, serialID)
}

func (r *ServiceImpl) MarkProductSerialsAsSold(ctx context.Context, serialIDs []string) error {
	return r.sqlc.MarkProductSerialsAsSold(ctx, serialIDs)
}

func (r *ServiceImpl) GetAvailableProducts(ctx context.Context, productID int64, amount int64) ([]model.ProductSerial, error) {
	rows, err := r.sqlc.GetAvailableProducts(ctx, sqlc.GetAvailableProductsParams{
		ProductID: productID,
		Amount:    int32(amount),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.ProductSerial, len(rows))
	for i, row := range rows {
		result[i] = model.ProductSerial{
			SerialID:    row.SerialID,
			ProductID:   row.ProductID,
			IsSold:      row.IsSold,
			IsActive:    row.IsActive,
			DateCreated: row.DateCreated.Time.UnixMilli(),
		}
	}

	return result, nil
}
