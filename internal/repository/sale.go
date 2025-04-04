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

func (r *RepositoryImpl) GetSale(ctx context.Context, id int64) (model.Sale, error) {
	row, err := r.sqlc.GetSale(ctx, id)
	if err != nil {
		return model.Sale{}, err
	}

	var discountPercent *int32
	if row.DiscountPercent.Valid {
		discountPercent = &row.DiscountPercent.Int32
	}

	return model.Sale{
		ID:               row.ID,
		Tag:              pgxutil.PgtypeToPtr[string](row.Tag),
		ProductModelID:   pgxutil.PgtypeToPtr[int64](row.ProductModelID),
		BrandID:          pgxutil.PgtypeToPtr[int64](row.BrandID),
		DateCreated:      row.DateCreated.Time.UnixMilli(),
		DateStarted:      row.DateStarted.Time.UnixMilli(),
		DateEnded:        util.PtrTimeToMilis(pgxutil.PgtypeToPtr[time.Time](row.DateEnded)),
		Quantity:         row.Quantity,
		Used:             row.Used,
		IsActive:         row.IsActive,
		DiscountPercent:  discountPercent,
		DiscountPrice:    pgxutil.PgtypeToPtr[int64](row.DiscountPrice),
		MaxDiscountPrice: row.MaxDiscountPrice,
	}, nil
}

type GetLatestSaleParams struct {
	ProductModelID int64
	BrandID        int64
	Tags           []string
}

func (r *RepositoryImpl) GetAvailableSales(ctx context.Context, params GetLatestSaleParams) ([]model.Sale, error) {
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
			Tag:              pgxutil.PgtypeToPtr[string](row.Tag),
			ProductModelID:   pgxutil.PgtypeToPtr[int64](row.ProductModelID),
			BrandID:          pgxutil.PgtypeToPtr[int64](row.BrandID),
			DateCreated:      row.DateCreated.Time.UnixMilli(),
			DateStarted:      row.DateStarted.Time.UnixMilli(),
			DateEnded:        util.PtrTimeToMilis(pgxutil.PgtypeToPtr[time.Time](row.DateEnded)),
			Quantity:         row.Quantity,
			Used:             row.Used,
			IsActive:         row.IsActive,
			DiscountPercent:  discountPercent,
			DiscountPrice:    pgxutil.PgtypeToPtr[int64](row.DiscountPrice),
			MaxDiscountPrice: row.MaxDiscountPrice,
		})
	}

	return sales, nil
}

type ListSalesParams struct {
	model.PaginationParams
	Tag             *string
	ProductModelID  *int64
	BrandID         *int64
	DateStartedFrom *int64
	DateStartedTo   *int64
	DateEndedFrom   *int64
	DateEndedTo     *int64
	IsActive        *bool
}

func (r *RepositoryImpl) CountSales(ctx context.Context, params ListSalesParams) (int64, error) {
	return r.sqlc.CountSales(ctx, sqlc.CountSalesParams{
		Tag:             *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Tag),
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		BrandID:         *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.BrandID),
		DateStartedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateStartedFrom)),
		DateStartedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateStartedTo)),
		DateEndedFrom:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateEndedFrom)),
		DateEndedTo:     *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateEndedTo)),
		IsActive:        *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.IsActive),
	})
}

func (r *RepositoryImpl) ListSales(ctx context.Context, params ListSalesParams) ([]model.Sale, error) {
	rows, err := r.sqlc.ListSales(ctx, sqlc.ListSalesParams{
		Limit:           params.Limit,
		Offset:          params.Offset(),
		Tag:             *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Tag),
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		BrandID:         *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.BrandID),
		DateStartedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateStartedFrom)),
		DateStartedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateStartedTo)),
		DateEndedFrom:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateEndedFrom)),
		DateEndedTo:     *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateEndedTo)),
		IsActive:        *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.IsActive),
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
			Tag:              pgxutil.PgtypeToPtr[string](row.Tag),
			ProductModelID:   pgxutil.PgtypeToPtr[int64](row.ProductModelID),
			BrandID:          pgxutil.PgtypeToPtr[int64](row.BrandID),
			DateCreated:      row.DateCreated.Time.UnixMilli(),
			DateStarted:      row.DateStarted.Time.UnixMilli(),
			DateEnded:        util.PtrTimeToMilis(pgxutil.PgtypeToPtr[time.Time](row.DateEnded)),
			Quantity:         row.Quantity,
			Used:             row.Used,
			IsActive:         row.IsActive,
			DiscountPercent:  discountPercent,
			DiscountPrice:    pgxutil.PgtypeToPtr[int64](row.DiscountPrice),
			MaxDiscountPrice: row.MaxDiscountPrice,
		})
	}

	return sales, nil
}

func (r *RepositoryImpl) CreateSale(ctx context.Context, sale model.Sale) (model.Sale, error) {
	row, err := r.sqlc.CreateSale(ctx, sqlc.CreateSaleParams{
		Tag:              *pgxutil.PtrToPgtype(&pgtype.Text{}, sale.Tag),
		ProductModelID:   *pgxutil.PtrToPgtype(&pgtype.Int8{}, &sale.ProductModelID),
		BrandID:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, &sale.BrandID),
		DateStarted:      *pgxutil.ValueToPgtype(&pgtype.Timestamptz{}, time.UnixMilli(sale.DateStarted)),
		DateEnded:        *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(sale.DateEnded)),
		Quantity:         sale.Quantity,
		Used:             0,
		IsActive:         sale.IsActive,
		DiscountPercent:  pgtype.Int4{Int32: util.DerefDefault(sale.DiscountPercent, 0), Valid: sale.DiscountPercent != nil},
		DiscountPrice:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, sale.DiscountPrice),
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
		Tag:              pgxutil.PgtypeToPtr[string](row.Tag),
		ProductModelID:   pgxutil.PgtypeToPtr[int64](row.ProductModelID),
		BrandID:          pgxutil.PgtypeToPtr[int64](row.BrandID),
		DateCreated:      row.DateCreated.Time.UnixMilli(),
		DateStarted:      row.DateStarted.Time.UnixMilli(),
		DateEnded:        util.PtrTimeToMilis(pgxutil.PgtypeToPtr[time.Time](row.DateEnded)),
		Quantity:         row.Quantity,
		Used:             row.Used,
		IsActive:         row.IsActive,
		DiscountPercent:  discountPercent,
		DiscountPrice:    pgxutil.PgtypeToPtr[int64](row.DiscountPrice),
		MaxDiscountPrice: row.MaxDiscountPrice,
	}, nil
}

type UpdateSaleParams struct {
	ID               int64
	Tag              *string
	ProductModelID   *int64
	BrandID          *int64
	DateStarted      *int64
	DateEnded        *int64
	Quantity         *int64
	Used             *int64
	IsActive         *bool
	DiscountPercent  *int32
	DiscountPrice    *int64
	MaxDiscountPrice *int64
}

func (r *RepositoryImpl) UpdateSale(ctx context.Context, params UpdateSaleParams) error {
	return r.sqlc.UpdateSale(ctx, sqlc.UpdateSaleParams{
		ID:               params.ID,
		Tag:              *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Tag),
		ProductModelID:   *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		BrandID:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.BrandID),
		DateStarted:      *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateStarted)),
		DateEnded:        *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateEnded)),
		Quantity:         *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Quantity),
		Used:             *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Used),
		IsActive:         *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		DiscountPercent:  pgtype.Int4{Int32: util.DerefDefault(params.DiscountPercent, 0), Valid: params.DiscountPercent != nil},
		DiscountPrice:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.DiscountPrice),
		MaxDiscountPrice: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.MaxDiscountPrice),
	})
}

func (r *RepositoryImpl) DeleteSale(ctx context.Context, id int64) error {
	return r.sqlc.DeleteSale(ctx, id)
}
