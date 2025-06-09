package storage

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/utils/ptr"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *ServiceImpl) GetSale(ctx context.Context, id int64) (model.Sale, error) {
	row, err := r.sqlc.GetSale(ctx, id)
	if err != nil {
		return model.Sale{}, err
	}

	var discountPercent *int32
	if row.DiscountPercent.Valid {
		discountPercent = &row.DiscountPercent.Int32
	}

	return model.Sale{
		ID: row.ID,

		DateCreated:      row.DateCreated.Time.UnixMilli(),
		DateStarted:      row.DateStarted.Time.UnixMilli(),
		DateEnded:        ptr.PtrTimeToMilis(PgtypeToPtr[time.Time](row.DateEnded)),
		CurrentStock:     row.CurrentStock.Int64,
		Used:             row.Used.Int64,
		IsActive:         row.IsActive,
		DiscountPercent:  discountPercent,
		DiscountPrice:    PgtypeToPtr[int64](row.DiscountPrice),
		MaxDiscountPrice: row.MaxDiscountPrice,
	}, nil
}

type GetLatestSaleParams struct {
	ProductModelID int64
	BrandID        int64
	Tags           []string
}

func (r *ServiceImpl) GetAvailableSales(ctx context.Context, params GetLatestSaleParams) ([]model.Sale, error) {
	rows, err := r.sqlc.GetAvailableSales(ctx, sqlc.GetAvailableSalesParams{
		ProductModelID: params.ProductModelID,
		BrandID:        params.BrandID,
		Tags:           params.Tags,
	})
	if err != nil {
		return nil, err
	}

	sales := make([]model.Sale, 0, len(rows))
	for _, row := range rows {
		var discountPercent *int32
		if row.DiscountPercent.Valid {
			discountPercent = &row.DiscountPercent.Int32
		}

		sales = append(sales, model.Sale{
			ID:               row.ID,
			Type:             model.SaleType(row.Type),
			ItemID:           row.ItemID,
			DateCreated:      row.DateCreated.Time.UnixMilli(),
			DateStarted:      row.DateStarted.Time.UnixMilli(),
			DateEnded:        ptr.PtrTimeToMilis(PgtypeToPtr[time.Time](row.DateEnded)),
			CurrentStock:     row.CurrentStock.Int64,
			Used:             row.Used.Int64,
			IsActive:         row.IsActive,
			DiscountPercent:  discountPercent,
			DiscountPrice:    PgtypeToPtr[int64](row.DiscountPrice),
			MaxDiscountPrice: row.MaxDiscountPrice,
		})
	}

	return sales, nil
}

type ListSalesParams struct {
	model.PaginationParams
	Type            *model.SaleType
	ItemID          *int64
	DateStartedFrom *int64
	DateStartedTo   *int64
	DateEndedFrom   *int64
	DateEndedTo     *int64
	IsActive        *bool
}

func (r *ServiceImpl) CountSales(ctx context.Context, params ListSalesParams) (int64, error) {
	return r.sqlc.CountSales(ctx, sqlc.CountSalesParams{
		Type:            *PtrBrandedToPgType(&pgtype.Text{}, params.Type),
		ItemID:          *PtrToPgtype(&pgtype.Int8{}, params.ItemID),
		DateStartedFrom: *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateStartedFrom)),
		DateStartedTo:   *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateStartedTo)),
		DateEndedFrom:   *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateEndedFrom)),
		DateEndedTo:     *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateEndedTo)),
		IsActive:        *PtrToPgtype(&pgtype.Bool{}, params.IsActive),
	})
}

func (r *ServiceImpl) ListSales(ctx context.Context, params ListSalesParams) ([]model.Sale, error) {
	rows, err := r.sqlc.ListSales(ctx, sqlc.ListSalesParams{
		Limit:           params.Limit,
		Offset:          params.Offset(),
		Type:            *PtrBrandedToPgType(&pgtype.Text{}, params.Type),
		ItemID:          *PtrToPgtype(&pgtype.Int8{}, params.ItemID),
		DateStartedFrom: *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateStartedFrom)),
		DateStartedTo:   *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateStartedTo)),
		DateEndedFrom:   *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateEndedFrom)),
		DateEndedTo:     *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateEndedTo)),
		IsActive:        *PtrToPgtype(&pgtype.Bool{}, params.IsActive),
	})
	if err != nil {
		return nil, err
	}

	sales := make([]model.Sale, 0, len(rows))
	for _, row := range rows {
		var discountPercent *int32
		if row.DiscountPercent.Valid {
			discountPercent = &row.DiscountPercent.Int32
		}

		sales = append(sales, model.Sale{
			ID:               row.ID,
			Type:             model.SaleType(row.Type),
			ItemID:           row.ItemID,
			DateCreated:      row.DateCreated.Time.UnixMilli(),
			DateStarted:      row.DateStarted.Time.UnixMilli(),
			DateEnded:        ptr.PtrTimeToMilis(PgtypeToPtr[time.Time](row.DateEnded)),
			CurrentStock:     row.CurrentStock.Int64,
			Used:             row.Used.Int64,
			IsActive:         row.IsActive,
			DiscountPercent:  discountPercent,
			DiscountPrice:    PgtypeToPtr[int64](row.DiscountPrice),
			MaxDiscountPrice: row.MaxDiscountPrice,
		})
	}

	return sales, nil
}

func (r *ServiceImpl) CreateSale(ctx context.Context, sale model.Sale) (model.Sale, error) {
	row, err := r.sqlc.CreateSale(ctx, sqlc.CreateSaleParams{
		Type:             sqlc.ProductSaleType(sale.Type),
		ItemID:           sale.ItemID,
		DateStarted:      *ValueToPgtype(&pgtype.Timestamptz{}, time.UnixMilli(sale.DateStarted)),
		DateEnded:        *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(sale.DateEnded)),
		CurrentStock:     sale.CurrentStock,
		IsActive:         sale.IsActive,
		DiscountPercent:  pgtype.Int4{Int32: ptr.DerefDefault(sale.DiscountPercent, 0), Valid: sale.DiscountPercent != nil},
		DiscountPrice:    *PtrToPgtype(&pgtype.Int8{}, sale.DiscountPrice),
		MaxDiscountPrice: sale.MaxDiscountPrice,
	})
	if err != nil {
		return model.Sale{}, err
	}

	var discountPercent *int32
	if row.DiscountPercent.Valid {
		discountPercent = &row.DiscountPercent.Int32
	}

	return model.Sale{
		ID:               row.ID,
		Type:             model.SaleType(row.Type),
		DateCreated:      row.DateCreated.Time.UnixMilli(),
		DateStarted:      row.DateStarted.Time.UnixMilli(),
		DateEnded:        ptr.PtrTimeToMilis(PgtypeToPtr[time.Time](row.DateEnded)),
		CurrentStock:     row.CurrentStock,
		Used:             row.Used,
		IsActive:         row.IsActive,
		DiscountPercent:  discountPercent,
		DiscountPrice:    PgtypeToPtr[int64](row.DiscountPrice),
		MaxDiscountPrice: row.MaxDiscountPrice,
	}, nil
}

type UpdateSaleParams struct {
	ID               int64
	Type             *model.SaleType
	ItemID           *int64
	DateStarted      *int64
	DateEnded        *int64
	IsActive         *bool
	DiscountPercent  *int32
	DiscountPrice    *int64
	MaxDiscountPrice *int64
}

func (r *ServiceImpl) UpdateSale(ctx context.Context, params UpdateSaleParams) error {
	return r.sqlc.UpdateSale(ctx, sqlc.UpdateSaleParams{
		ID:               params.ID,
		Type:             *PtrBrandedToPgType(&sqlc.NullProductSaleType{}, params.Type),
		ItemID:           *PtrToPgtype(&pgtype.Int8{}, params.ItemID),
		DateStarted:      *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateStarted)),
		DateEnded:        *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateEnded)),
		IsActive:         *PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		DiscountPercent:  pgtype.Int4{Int32: ptr.DerefDefault(params.DiscountPercent, 0), Valid: params.DiscountPercent != nil},
		DiscountPrice:    *PtrToPgtype(&pgtype.Int8{}, params.DiscountPrice),
		MaxDiscountPrice: *PtrToPgtype(&pgtype.Int8{}, params.MaxDiscountPrice),
	})
}

func (r *ServiceImpl) DeleteSale(ctx context.Context, id int64) error {
	return r.sqlc.DeleteSale(ctx, id)
}
