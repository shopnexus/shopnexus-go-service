package storage

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/utils/ptr"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type ExistsRefundParams struct {
	ProductOnPaymentID int64
	UserID             int64
}

func (r *ServiceImpl) ExistsRefund(ctx context.Context, params ExistsRefundParams) (bool, error) {
	return r.sqlc.ExistsRefund(ctx, sqlc.ExistsRefundParams{
		ProductOnPaymentID: params.ProductOnPaymentID,
		UserID:             params.UserID,
	})
}

type GetRefundParams struct {
	ID     int64
	UserID *int64
}

func (r *ServiceImpl) GetRefund(ctx context.Context, params GetRefundParams) (model.Refund, error) {
	row, err := r.sqlc.GetRefund(ctx, sqlc.GetRefundParams{
		ID:     params.ID,
		UserID: *PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})
	if err != nil {
		return model.Refund{}, err
	}

	return model.Refund{
		ID:                 row.ID,
		ProductOnPaymentID: row.ProductOnPaymentID,
		Method:             model.RefundMethod(row.Method),
		Status:             model.Status(row.Status),
		Reason:             row.Reason,
		Address:            row.Address,
		DateCreated:        row.DateCreated.Time.UnixMilli(),
		DateUpdated:        row.DateUpdated.Time.UnixMilli(),
		Resources:          row.Resources,
	}, nil
}

type ListRefundsParams struct {
	model.PaginationParams
	UserID             *int64
	ProductOnPaymentID *int64
	Method             *model.RefundMethod
	Status             *model.Status
	Reason             *string
	Address            *string
	DateCreatedFrom    *int64
	DateCreatedTo      *int64
}

func (r *ServiceImpl) CountRefunds(ctx context.Context, params ListRefundsParams) (int64, error) {
	return r.sqlc.CountRefunds(ctx, sqlc.CountRefundsParams{
		UserID:             *PtrToPgtype(&pgtype.Int8{}, params.UserID),
		ProductOnPaymentID: *PtrToPgtype(&pgtype.Int8{}, params.ProductOnPaymentID),
		Method:             *PtrBrandedToPgType(&sqlc.NullPaymentRefundMethod{}, params.Method),
		Status:             *PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		Reason:             *PtrToPgtype(&pgtype.Text{}, params.Reason),
		Address:            *PtrToPgtype(&pgtype.Text{}, params.Address),
		DateCreatedFrom:    *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:      *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (r *ServiceImpl) ListRefunds(ctx context.Context, params ListRefundsParams) ([]model.Refund, error) {
	rows, err := r.sqlc.ListRefunds(ctx, sqlc.ListRefundsParams{
		Offset:             params.Offset(),
		Limit:              params.Limit,
		UserID:             *PtrToPgtype(&pgtype.Int8{}, params.UserID),
		ProductOnPaymentID: *PtrToPgtype(&pgtype.Int8{}, params.ProductOnPaymentID),
		Method:             *PtrBrandedToPgType(&sqlc.NullPaymentRefundMethod{}, params.Method),
		Status:             *PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		Reason:             *PtrToPgtype(&pgtype.Text{}, params.Reason),
		Address:            *PtrToPgtype(&pgtype.Text{}, params.Address),
		DateCreatedFrom:    *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:      *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedTo)),
	})
	if err != nil {
		return nil, err
	}

	var refunds []model.Refund
	for _, row := range rows {
		refunds = append(refunds, model.Refund{
			ID:                 row.ID,
			ProductOnPaymentID: row.ProductOnPaymentID,
			Method:             model.RefundMethod(row.Method),
			Status:             model.Status(row.Status),
			Reason:             row.Reason,
			Address:            row.Address,
			DateCreated:        row.DateCreated.Time.UnixMilli(),
			DateUpdated:        row.DateUpdated.Time.UnixMilli(),
			Resources:          row.Resources,
		})
	}

	return refunds, nil
}

func (r *ServiceImpl) CreateRefund(ctx context.Context, refund model.Refund) (model.Refund, error) {
	row, err := r.sqlc.CreateRefund(ctx, sqlc.CreateRefundParams{
		ProductOnPaymentID: refund.ProductOnPaymentID,
		Method:             sqlc.PaymentRefundMethod(refund.Method),
		Status:             sqlc.PaymentStatus(refund.Status),
		Reason:             refund.Reason,
		Address:            refund.Address,
	})
	if err != nil {
		return model.Refund{}, err
	}

	if err := r.AddResources(ctx, row.ID, model.ResourceTypeRefund, refund.Resources); err != nil {
		return model.Refund{}, err
	}

	return model.Refund{
		ID:                 row.ID,
		ProductOnPaymentID: refund.ProductOnPaymentID,
		Method:             refund.Method,
		Status:             refund.Status,
		Reason:             refund.Reason,
		Address:            refund.Address,
		DateCreated:        time.Now().UnixMilli(),
		DateUpdated:        time.Now().UnixMilli(),
		Resources:          refund.Resources,
	}, nil
}

type UpdateRefundParams struct {
	ID        int64
	UserID    *int64
	Method    *model.RefundMethod
	Status    *model.Status
	Reason    *string
	Address   *string
	Resources *[]string
}

func (r *ServiceImpl) UpdateRefund(ctx context.Context, params UpdateRefundParams) error {
	err := r.sqlc.UpdateRefund(ctx, sqlc.UpdateRefundParams{
		ID:      params.ID,
		UserID:  *PtrToPgtype(&pgtype.Int8{}, params.UserID),
		Method:  *PtrBrandedToPgType(&sqlc.NullPaymentRefundMethod{}, params.Method),
		Status:  *PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		Reason:  *PtrToPgtype(&pgtype.Text{}, params.Reason),
		Address: *PtrToPgtype(&pgtype.Text{}, params.Address),
	})

	if params.Resources != nil {
		if err := r.UpdateResources(ctx, params.ID, model.ResourceTypeRefund, *params.Resources); err != nil {
			return err
		}
	}

	return err
}

type DeleteRefundParams struct {
	ID     int64
	UserID *int64
}

func (r *ServiceImpl) DeleteRefund(ctx context.Context, params DeleteRefundParams) error {
	return r.sqlc.DeleteRefund(ctx, sqlc.DeleteRefundParams{
		ID:     params.ID,
		UserID: *PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})
}

type CanRefundParams struct {
	ProductOnPaymentID int64
	UserID             *int64
}

func (r *ServiceImpl) CanRefund(ctx context.Context, params CanRefundParams) (bool, error) {
	return r.sqlc.CanRefund(ctx, sqlc.CanRefundParams{
		ID:     params.ProductOnPaymentID,
		UserID: *PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})
}
