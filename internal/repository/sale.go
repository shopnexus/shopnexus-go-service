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

func (r *Repository) CreateSale(ctx context.Context, sale model.Sale) (model.Sale, error) {
	row, err := r.sqlc.CreateSale(ctx, sqlc.CreateSaleParams{
		Tag:             *pgxutil.PtrToPgtype(&pgtype.Text{}, sale.Tag),
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, &sale.ProductModelID),
		BrandID:         *pgxutil.PtrToPgtype(&pgtype.Int8{}, &sale.BrandID),
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
		Tag:             pgxutil.PgtypeToPtr[string](row.Tag),
		ProductModelID:  pgxutil.PgtypeToPtr[int64](row.ProductModelID),
		BrandID:         pgxutil.PgtypeToPtr[int64](row.BrandID),
		DateStarted:     row.DateStarted.Time.UnixMilli(),
		DateEnded:       util.PtrTimeToMilis(pgxutil.PgtypeToPtr[time.Time](row.DateEnded)),
		Quantity:        row.Quantity,
		Used:            row.Used,
		IsActive:        row.IsActive,
		DiscountPercent: pgxutil.PgtypeToPtr[int64](row.DiscountPercent),
		DiscountPrice:   pgxutil.PgtypeToPtr[int64](row.DiscountPrice),
	}, nil
}

type UpdateSaleParams struct {
	ID              int64
	Tag             *string
	ProductModelID  *int64
	BrandID         *int64
	DateStarted     *int64
	DateEnded       *int64
	Quantity        *int64
	Used            *int64
	IsActive        *bool
	DiscountPercent *int64
	DiscountPrice   *int64
}

func (r *Repository) UpdateSale(ctx context.Context, params UpdateSaleParams) error {
	return r.sqlc.UpdateSale(ctx, sqlc.UpdateSaleParams{
		ID:              params.ID,
		Tag:             *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Tag),
		ProductModelID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ProductModelID),
		BrandID:         *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.BrandID),
		DateStarted:     *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateStarted)),
		DateEnded:       *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateEnded)),
		Quantity:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Quantity),
		Used:            *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Used),
		IsActive:        *pgxutil.PtrToPgtype(&pgtype.Bool{}, params.IsActive),
		DiscountPercent: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.DiscountPercent),
		DiscountPrice:   *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.DiscountPrice),
	})
}

func (r *Repository) DeleteSale(ctx context.Context, id int64) error {
	return r.sqlc.DeleteSale(ctx, id)
}
