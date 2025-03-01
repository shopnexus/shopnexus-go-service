package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetPaymentParams struct {
	ID     int64
	UserID *int64
}

func (r *Repository) ExistsPayment(ctx context.Context, params GetPaymentParams) (bool, error) {
	return r.sqlc.ExistsPayment(ctx, sqlc.ExistsPaymentParams{
		ID:     params.ID,
		UserID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})
}

func (r *Repository) GetPayment(ctx context.Context, params GetPaymentParams) (model.Payment, error) {
	row, err := r.sqlc.GetPayment(ctx, sqlc.GetPaymentParams{
		ID:     params.ID,
		UserID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})
	if err != nil {
		return model.Payment{}, err
	}

	products, err := r.GetPaymentProducts(ctx, row.ID)
	if err != nil {
		return model.Payment{}, err
	}

	return model.Payment{
		ID:            row.ID,
		UserID:        row.UserID,
		Address:       row.Address,
		PaymentMethod: model.PaymentMethod(row.Method),
		Total:         row.Total,
		Status:        model.Status(row.Status),
		DateCreated:   row.DateCreated.Time.UnixMilli(),
		Products:      products,
	}, nil
}

func (r *Repository) GetPaymentProducts(ctx context.Context, paymentID int64) ([]model.ProductOnPayment, error) {
	rows, err := r.sqlc.GetPaymentProducts(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	var products []model.ProductOnPayment
	for _, row := range rows {
		products = append(products, model.ProductOnPayment{
			ItemQuantityBase: model.ItemQuantityBase[string]{
				ItemID:   row.ProductSerialID,
				Quantity: row.Quantity,
			},
			Price:      row.Price,
			TotalPrice: row.TotalPrice,
		})
	}

	return products, nil
}

func (r *Repository) CreatePayment(ctx context.Context, payment model.Payment) (model.Payment, error) {
	row, err := r.sqlc.CreatePayment(ctx, sqlc.CreatePaymentParams{
		UserID:  payment.UserID,
		Method:  sqlc.PaymentPaymentMethod(payment.PaymentMethod),
		Status:  sqlc.PaymentStatus(payment.Status),
		Address: payment.Address,
		Total:   payment.Total,
	})
	if err != nil {
		return model.Payment{}, err
	}

	var createArgs []sqlc.CreatePaymentProductsParams
	for _, product := range payment.Products {
		createArgs = append(createArgs, sqlc.CreatePaymentProductsParams{
			PaymentID:       row.ID,
			ProductSerialID: product.ItemID,
			Quantity:        product.Quantity,
			Price:           product.Price,
			TotalPrice:      product.TotalPrice,
		})
	}

	_, err = r.sqlc.CreatePaymentProducts(ctx, createArgs)
	if err != nil {
		return model.Payment{}, err
	}

	return model.Payment{
		ID:            row.ID,
		UserID:        row.UserID,
		Address:       row.Address,
		PaymentMethod: model.PaymentMethod(row.Method),
		Total:         row.Total,
		Status:        model.Status(row.Status),
		DateCreated:   row.DateCreated.Time.UnixMilli(),
		Products:      payment.Products,
	}, nil
}

type GetRefundParams struct {
	ID     int64
	UserID *int64
}

func (r *Repository) GetRefund(ctx context.Context, params GetRefundParams) (model.Refund, error) {
	row, err := r.sqlc.GetRefund(ctx, sqlc.GetRefundParams{
		ID:     params.ID,
		UserID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})
	if err != nil {
		return model.Refund{}, err
	}

	return model.Refund{
		ID:          row.ID,
		PaymentID:   row.PaymentID,
		Method:      model.RefundMethod(row.Method),
		Status:      model.Status(row.Status),
		Reason:      row.Reason,
		Address:     pgxutil.PgtypeToPtr[string](row.Address),
		Resources:   row.Resources,
		DateCreated: row.DateCreated.Time.UnixMilli(),
		DateUpdated: row.DateUpdated.Time.UnixMilli(),
	}, nil
}

func (r *Repository) CreateRefund(ctx context.Context, refund model.Refund) (model.Refund, error) {
	row, err := r.sqlc.CreateRefund(ctx, sqlc.CreateRefundParams{
		PaymentID: refund.PaymentID,
		Method:    sqlc.PaymentRefundMethod(refund.Method),
		Status:    sqlc.PaymentStatus(refund.Status),
		Reason:    refund.Reason,
		Address:   *pgxutil.PtrToPgtype(&pgtype.Text{}, refund.Address),
		Resources: refund.Resources,
	})
	if err != nil {
		return model.Refund{}, err
	}

	return model.Refund{
		ID:          row.ID,
		PaymentID:   refund.PaymentID,
		Method:      refund.Method,
		Status:      refund.Status,
		Reason:      refund.Reason,
		Address:     refund.Address,
		Resources:   row.Resources,
		DateCreated: time.Now().UnixMilli(),
		DateUpdated: time.Now().UnixMilli(),
	}, nil
}

type UpdateRefundParams struct {
	ID          int64
	Method      *model.RefundMethod
	Status      *model.Status
	Reason      *string
	Address     *string
	NullAddress bool
}

func (r *Repository) UpdateRefund(ctx context.Context, params UpdateRefundParams) error {
	err := r.sqlc.UpdateRefund(ctx, sqlc.UpdateRefundParams{
		ID:          params.ID,
		Method:      *pgxutil.PtrBrandedToPgType(&sqlc.NullPaymentRefundMethod{}, params.Method),
		Status:      *pgxutil.PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		Reason:      *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Reason),
		Address:     *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Address),
		NullAddress: params.NullAddress,
	})

	return err
}

func (r *Repository) DeleteRefund(ctx context.Context, id int64) error {
	return r.sqlc.DeleteRefund(ctx, id)
}
