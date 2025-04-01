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

type ExistsRefundParams struct {
	PaymentID int64
	UserID    int64
}

func (r *RepositoryImpl) ExistsRefund(ctx context.Context, params ExistsRefundParams) (bool, error) {
	return r.sqlc.ExistsRefund(ctx, sqlc.ExistsRefundParams{
		PaymentID: params.PaymentID,
		UserID:    params.UserID,
	})
}

type GetRefundParams struct {
	ID     int64
	UserID *int64
}

func (r *RepositoryImpl) GetRefund(ctx context.Context, params GetRefundParams) (model.Refund, error) {
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
		Address:     row.Address,
		Resources:   row.Resources,
		DateCreated: row.DateCreated.Time.UnixMilli(),
		DateUpdated: row.DateUpdated.Time.UnixMilli(),
	}, nil
}

type ListRefundsParams struct {
	model.PaginationParams
	UserID          *int64
	PaymentID       *int64
	Method          *model.RefundMethod
	Status          *model.Status
	Reason          *string
	Address         *string
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (r *RepositoryImpl) CountRefunds(ctx context.Context, params ListRefundsParams) (int64, error) {
	return r.sqlc.CountRefunds(ctx, sqlc.CountRefundsParams{
		UserID:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
		PaymentID:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.PaymentID),
		Method:          *pgxutil.PtrBrandedToPgType(&sqlc.NullPaymentRefundMethod{}, params.Method),
		Status:          *pgxutil.PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		Reason:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Reason),
		Address:         *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Address),
		DateCreatedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (r *RepositoryImpl) ListRefunds(ctx context.Context, params ListRefundsParams) ([]model.Refund, error) {
	rows, err := r.sqlc.ListRefunds(ctx, sqlc.ListRefundsParams{
		Offset:          params.Offset(),
		Limit:           params.Limit,
		UserID:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
		PaymentID:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.PaymentID),
		Method:          *pgxutil.PtrBrandedToPgType(&sqlc.NullPaymentRefundMethod{}, params.Method),
		Status:          *pgxutil.PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		Reason:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Reason),
		Address:         *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Address),
		DateCreatedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedTo)),
	})
	if err != nil {
		return nil, err
	}

	var refunds []model.Refund
	for _, row := range rows {
		refunds = append(refunds, model.Refund{
			ID:          row.ID,
			PaymentID:   row.PaymentID,
			Method:      model.RefundMethod(row.Method),
			Status:      model.Status(row.Status),
			Reason:      row.Reason,
			Address:     row.Address,
			Resources:   row.Resources,
			DateCreated: row.DateCreated.Time.UnixMilli(),
			DateUpdated: row.DateUpdated.Time.UnixMilli(),
		})
	}

	return refunds, nil
}

func (r *RepositoryImpl) CreateRefund(ctx context.Context, refund model.Refund) (model.Refund, error) {
	row, err := r.sqlc.CreateRefund(ctx, sqlc.CreateRefundParams{
		PaymentID: refund.PaymentID,
		Method:    sqlc.PaymentRefundMethod(refund.Method),
		Status:    sqlc.PaymentStatus(refund.Status),
		Reason:    refund.Reason,
		Address:   refund.Address,
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
	ID      int64
	UserID  *int64
	Method  *model.RefundMethod
	Status  *model.Status
	Reason  *string
	Address *string
}

func (r *RepositoryImpl) UpdateRefund(ctx context.Context, params UpdateRefundParams) error {
	err := r.sqlc.UpdateRefund(ctx, sqlc.UpdateRefundParams{
		ID:      params.ID,
		Method:  *pgxutil.PtrBrandedToPgType(&sqlc.NullPaymentRefundMethod{}, params.Method),
		Status:  *pgxutil.PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		Reason:  *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Reason),
		Address: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Address),
		UserID:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})

	return err
}

func (r *RepositoryImpl) DeleteRefund(ctx context.Context, id int64) error {
	return r.sqlc.DeleteRefund(ctx, id)
}
